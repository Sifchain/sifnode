package relayer

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sethvargo/go-password/password"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/go-amino"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

const (
	transactionInterval = 10 * time.Second
	trailingBlocks      = 50
	ethLevelDBKey       = "ethereumLastProcessedBlock"
)

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	Cdc                     *codec.LegacyAmino
	EthProvider             string
	TmProvider              string
	RegistryContractAddress common.Address
	ValidatorName           string
	ValidatorAddress        sdk.ValAddress
	CliCtx                  sdkContext.CLIContext
	TxBldr                  authtypes.TxBuilder
	PrivateKey              *ecdsa.PrivateKey
	TempPassword            string
	EventsBuffer            types.EthEventBuffer
	DB                      *leveldb.DB
	SugaredLogger           *zap.SugaredLogger
}

// NewKeybase create a new keybase instance
func NewKeybase(validatorMoniker, mnemonic, password string) (keys.Keybase, keys.Info, error) {
	keybase := keys.NewInMemory()
	hdpath := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	info, err := keybase.CreateAccount(validatorMoniker, mnemonic, "", password, hdpath.String(), keys.Secp256k1)
	if err != nil {
		return nil, nil, err
	}

	return keybase, info, nil
}

// NewEthereumSub initializes a new EthereumSub
func NewEthereumSub(inBuf io.Reader, rpcURL string, cdc *codec.LegacyAmino, validatorMoniker, chainID, ethProvider string,
	registryContractAddress common.Address, privateKey *ecdsa.PrivateKey, mnemonic string, db *leveldb.DB, sugaredLogger *zap.SugaredLogger) (EthereumSub, error) {

	tempPassword, _ := password.Generate(32, 5, 0, false, false)
	keybase, info, err := NewKeybase(validatorMoniker, mnemonic, tempPassword)
	if err != nil {
		sugaredLogger.Errorw("failed to initialize a new key base.",
			errorMessageKey, err.Error())
		return EthereumSub{}, err
	}

	validatorAddress := sdk.ValAddress(info.GetAddress())

	// Load CLI context and Tx builder
	cliCtx, err := LoadTendermintCLIContext(cdc, validatorAddress, validatorMoniker, rpcURL, chainID, sugaredLogger)
	if err != nil {
		return EthereumSub{}, err
	}

	txBldr := authtypes.NewTxBuilderFromCLI(inBuf).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID).
		WithKeybase(keybase)

	return EthereumSub{
		Cdc:                     cdc,
		EthProvider:             ethProvider,
		TmProvider:              rpcURL,
		RegistryContractAddress: registryContractAddress,
		ValidatorName:           validatorMoniker,
		ValidatorAddress:        validatorAddress,
		CliCtx:                  cliCtx,
		TxBldr:                  txBldr,
		PrivateKey:              privateKey,
		TempPassword:            tempPassword,
		EventsBuffer:            types.NewEthEventBuffer(),
		DB:                      db,
		SugaredLogger:           sugaredLogger,
	}, nil
}

// LoadTendermintCLIContext : loads CLI context for tendermint txs
func LoadTendermintCLIContext(appCodec *amino.Codec, validatorAddress sdk.ValAddress, validatorName string,
	rpcURL string, chainID string, sugaredLogger *zap.SugaredLogger) (sdkContext.CLIContext, error) {
	// Create the new CLI context
	cliCtx := sdkContext.NewCLIContext().
		WithCodec(appCodec).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(validatorName)

	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	// Confirm that the validator's address exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists(sdk.AccAddress(validatorAddress))
	if err != nil {
		sugaredLogger.Errorw("validator address not exists.",
			errorMessageKey, err.Error())
		return sdkContext.CLIContext{}, err
	}
	return cliCtx, nil
}

// Start an Ethereum chain subscription
func (sub EthereumSub) Start(completionEvent *sync.WaitGroup) {
	defer completionEvent.Done()
	time.Sleep(time.Second)
	client, err := SetupWebsocketEthClient(sub.EthProvider)
	if err != nil {
		sub.SugaredLogger.Errorw("SetupWebsocketEthClient failed.",
			errorMessageKey, err.Error())

		completionEvent.Add(1)
		go sub.Start(completionEvent)
		return
	}
	defer client.Close()
	sub.SugaredLogger.Infow("Started Ethereum websocket with provider:",
		"Ethereum provider", sub.EthProvider)

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get network ID.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(completionEvent)
		return
	}

	// We will check logs for new events
	logs := make(chan ctypes.Log)
	defer close(logs)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	// get the bridgebank address from the registry contract
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(client, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
	}

	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)

	// Listen the new header
	heads := make(chan *ctypes.Header)
	defer close(heads)
	subHead, err := client.SubscribeNewHead(context.Background(), heads)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to subscribe new head.",
			errorMessageKey, err.Error())
		return
	}
	defer subHead.Unsubscribe()

	var lastProcessedBlock *big.Int

	data, err := sub.DB.Get([]byte(ethLevelDBKey), nil)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the last ethereum block from level db.",
			errorMessageKey, err.Error())
		lastProcessedBlock = big.NewInt(0)
	} else {
		lastProcessedBlock = new(big.Int).SetBytes(data)
	}

	for {
		select {
		// Handle any errors
		case <-quit:
			return
		case err := <-subHead.Err():
			sub.SugaredLogger.Errorw("failed to subscribe ethereum header.",
				errorMessageKey, err.Error())
			completionEvent.Add(1)
			go sub.Start(completionEvent)
			return
		case newHead := <-heads:
			sub.SugaredLogger.Infow("receive new ethereum header.",
				"ethereum block number", newHead.Number,
				"ethereum block hash", newHead.Hash())

			startingBigInt := newHead.Number
			endingBlock := startingBigInt.Sub(startingBigInt, big.NewInt(trailingBlocks))

			// if the current block number - trailing blocks is negative, don't bother
			// going deeper into the function.
			if endingBlock.Cmp(big.NewInt(0)) == -1 {
				sub.SugaredLogger.Infow("Ending block index negative. Cancelling run.")
				continue
			}

			// If the last processed block is the default (0), then go and set it to the difference of ending block minus 1
			// The user who starts this must provide a valid last processed block
			if lastProcessedBlock.Cmp(big.NewInt(0)) == 0 {
				lastProcessedBlock = endingBlock
			}

			sub.SugaredLogger.Infow("Processing events from blocks.",
				"lastProcessedBlock", lastProcessedBlock,
				"endingBlock", endingBlock)

			// query event data from this specific block range
			ethLogs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: lastProcessedBlock,
				ToBlock:   endingBlock,
				Addresses: []common.Address{bridgeBankAddress},
			})

			if err != nil {
				sub.SugaredLogger.Errorw("failed to get events from bridgebank.",
					errorMessageKey, err.Error(),
					"block number", newHead.Number)
				// if you have an error getting the logs from the block, continue and keep
				// the current last processed block so we keep retrying
				continue
			}
			// Assumption here is that we will repeat a failing block because we return if there is an error retrieving logs
			sub.SugaredLogger.Infow("received events from bridgebank.",
				"lastProcessedBlock", lastProcessedBlock,
				"endingBlock", endingBlock)

			var events []types.EthereumEvent

			// loop over ethlogs, and build an array of burn/lock events
			for _, ethLog := range ethLogs {
				log.Printf("Processed events from block %v", ethLog.BlockNumber)
				event, isBurnLock, err := sub.logToEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI, ethLog)
				if err != nil {
					sub.SugaredLogger.Errorw("failed to transform from log to event.",
						errorMessageKey, err.Error())
					continue
				}
				if !isBurnLock {
					sub.SugaredLogger.Infow("not burn or lock event, continue events.")
					continue
				}
				events = append(events, event)
			}

			if len(events) > 0 {
				if err := sub.handleEthereumEvent(events); err != nil {
					sub.SugaredLogger.Errorw("failed to handle ethereum event.",
						errorMessageKey, err.Error())
				}
				time.Sleep(transactionInterval)
			}
			// add 1 to the current block so we don't reprocess it
			endingBlock = endingBlock.Add(endingBlock, big.NewInt(1))
			// save the current ending block + 1 to the lastprocessed block to ensure we keep reading blocks sequentially and don't repeat blocks
			err = sub.DB.Put([]byte(ethLevelDBKey), endingBlock.Bytes(), nil)
			if err != nil {
				// if you can't write to leveldb, then error out as something is seriously amiss
				log.Fatalf("Error saving lastProcessedBlock to leveldb: %v", err)
			}
			lastProcessedBlock = endingBlock
		}
	}
}

func (sub EthereumSub) getAllClaims(fromBlock int64, toBlock int64) []types.EthereumBridgeClaim {
	log.Printf("Replay get all ethereum bridge claim from block %d to block %d\n", fromBlock, toBlock)

	var claimArray []types.EthereumBridgeClaim
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client, error is %s\n", err.Error())
		return claimArray
	}

	if err := client.Start(); err != nil {
		log.Printf("failed to start a client, error is %s\n", err.Error())
		return claimArray
	}

	defer client.Stop() //nolint:errcheck

	for blockNumber := fromBlock; blockNumber < toBlock; {
		tmpBlockNumber := blockNumber
		ctx := context.Background()
		block, err := client.BlockResults(ctx, &tmpBlockNumber)
		blockNumber++
		log.Printf("Replay start to process block %d\n", blockNumber)

		if err != nil {
			log.Printf("failed to start a client %s\n", err.Error())
			continue
		}

		for _, result := range block.TxsResults {
			for _, event := range result.Events {
				log.Printf("Replay get an event %s\n", event.GetType())
				if event.GetType() == "create_claim" {
					claim, err := txs.AttributesToEthereumBridgeClaim(event.GetAttributes())
					if err != nil {
						continue
					}

					// Check if sender is me
					if claim.CosmosSender.Equals(sub.ValidatorAddress) {
						log.Printf("We got a eth bridge claim message %s\n", claim.EthereumSender.String())
						claimArray = append(claimArray, claim)
					}
				}
			}
		}
	}

	return claimArray
}

// EventProcessed check if the event processed by relayer
func EventProcessed(bridgeClaims []types.EthereumBridgeClaim, event types.EthereumEvent) bool {
	for _, claim := range bridgeClaims {
		if event.From == claim.EthereumSender && event.Nonce.Cmp(claim.Nonce.BigInt()) == 0 {
			return true
		}
	}
	return false
}

// Replay the missed events
func (sub EthereumSub) Replay(fromBlock int64, toBlock int64, cosmosFromBlock int64, cosmosToBlock int64) {
	log.Printf("ethereum replay for %d block to %d block\n", fromBlock, toBlock)

	bridgeClaims := sub.getAllClaims(cosmosFromBlock, cosmosToBlock)
	log.Printf("found out %d bridgeClaims\n", len(bridgeClaims))

	client, err := SetupRPCEthClient(sub.EthProvider)
	if err != nil {
		log.Printf("failed to connect ethereum node, error is %s\n", err.Error())
		return
	}
	defer client.Close()

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Printf("failed to get chain ID, error is %s\n", err.Error())
		return
	}

	// Get the contract address for this subscription
	subContractAddress, err := txs.GetAddressFromBridgeRegistry(client, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Printf("failed to get contract address, error is %s\n", err.Error())
		return
	}
	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)
	// We need the address in []bytes for the query
	subQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subContractAddress},
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
	}

	logs, err := client.FilterLogs(context.Background(), subQuery)
	if err != nil {
		log.Printf("failed to get filter log, error is %s\n", err.Error())
		return
	}

	for _, ethLog := range logs {
		// Before deal with it, we need check in cosmos if it is already handled by myself bofore.
		event, isBurnLock, err := sub.logToEvent(clientChainID, subContractAddress, bridgeBankContractABI, ethLog)
		if err != nil {
			log.Println("Failed to get event from ethereum log")
		} else if isBurnLock {
			log.Println(fmt.Sprintf("found out a burn lock event"))
			if !EventProcessed(bridgeClaims, event) {
				err := sub.handleEthereumEvent([]types.EthereumEvent{event})
				time.Sleep(transactionInterval)
				if err != nil {
					log.Printf("failed to handle ethereum event, error is %s\n", err.Error())
				}
			} else {
				log.Println("event already processed by me.")
			}
		}
	}
}

// logToEvent unpacks an Ethereum event
func (sub EthereumSub) logToEvent(clientChainID *big.Int, contractAddress common.Address,
	contractABI abi.ABI, cLog ctypes.Log) (types.EthereumEvent, bool, error) {
	// Parse the event's attributes via contract ABI
	event := types.EthereumEvent{}
	eventLogLockSignature := contractABI.Events[types.LogLock.String()].ID().Hex()
	eventLogBurnSignature := contractABI.Events[types.LogBurn.String()].ID().Hex()

	var eventName string
	switch cLog.Topics[0].Hex() {
	case eventLogBurnSignature:
		eventName = types.LogBurn.String()
	case eventLogLockSignature:
		eventName = types.LogLock.String()
	}

	// If event is not expected
	if eventName == "" {
		return event, false, nil
	}

	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.SugaredLogger.Errorw(".",
			errorMessageKey, err.Error())
		return event, false, err
	}
	event.BridgeContractAddress = contractAddress
	event.EthereumChainID = clientChainID
	if eventName == types.LogBurn.String() {
		event.ClaimType = ethbridge.ClaimType_CLAIM_TYPE_BURN
	} else {
		event.ClaimType = ethbridge.ClaimType_CLAIM_TYPE_LOCK
	}
	sub.SugaredLogger.Infow("receive an event.",
		"event", event.String())

	// Add the event to the record
	types.NewEventWrite(cLog.TxHash.Hex(), event)
	return event, true, nil
}

// handleEthereumEvent unpacks an Ethereum event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func (sub EthereumSub) handleEthereumEvent(events []types.EthereumEvent) error {
	var prophecyClaims []ethbridge.EthBridgeClaim

	for _, event := range events {
		prophecyClaim, err := txs.EthereumEventToEthBridgeClaim(sub.ValidatorAddress, event)
		if err != nil {
			sub.SugaredLogger.Errorw(".",
				errorMessageKey, err.Error())
		} else {
			prophecyClaims = append(prophecyClaims, prophecyClaim)
		}
	}
	sub.SugaredLogger.Infow("relay prophecy claims to cosmos.",
		"prophecy claims length", len(prophecyClaims))

	if len(events) == 0 {
		return nil
	}

	return txs.RelayToCosmos(sub.Cdc, sub.ValidatorName, sub.TempPassword, prophecyClaims, sub.CliCtx, sub.TxBldr, sub.SugaredLogger)
}
