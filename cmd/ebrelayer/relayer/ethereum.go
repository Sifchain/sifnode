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
	tmLog "github.com/tendermint/tendermint/libs/log"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

const (
	transactionInterval = 10 * time.Second
	trailingBlocks      = 50
	levelDbFile         = "relayerdb"
)

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	Cdc                     *codec.Codec
	EthProvider             string
	RegistryContractAddress common.Address
	ValidatorName           string
	ValidatorAddress        sdk.ValAddress
	CliCtx                  sdkContext.CLIContext
	TxBldr                  authtypes.TxBuilder
	PrivateKey              *ecdsa.PrivateKey
	TempPassword            string
	EventsBuffer            types.EthEventBuffer
	Logger                  tmLog.Logger
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
func NewEthereumSub(inBuf io.Reader, rpcURL string, cdc *codec.Codec, validatorMoniker, chainID, ethProvider string,
	registryContractAddress common.Address, privateKey *ecdsa.PrivateKey, mnemonic string, logger tmLog.Logger) (EthereumSub, error) {

	tempPassword, _ := password.Generate(32, 5, 0, false, false)
	keybase, info, err := NewKeybase(validatorMoniker, mnemonic, tempPassword)
	if err != nil {
		return EthereumSub{}, err
	}

	validatorAddress := sdk.ValAddress(info.GetAddress())

	// Load CLI context and Tx builder
	cliCtx, err := LoadTendermintCLIContext(cdc, validatorAddress, validatorMoniker, rpcURL, chainID)
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
		RegistryContractAddress: registryContractAddress,
		ValidatorName:           validatorMoniker,
		ValidatorAddress:        validatorAddress,
		CliCtx:                  cliCtx,
		TxBldr:                  txBldr,
		PrivateKey:              privateKey,
		TempPassword:            tempPassword,
		EventsBuffer:            types.NewEthEventBuffer(),
		Logger:                  logger,
	}, nil
}

// LoadTendermintCLIContext : loads CLI context for tendermint txs
func LoadTendermintCLIContext(appCodec *amino.Codec, validatorAddress sdk.ValAddress, validatorName string,
	rpcURL string, chainID string) (sdkContext.CLIContext, error) {
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
		log.Println(err)
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
		sub.Logger.Error("SetupWebsocketEthClient failed: ", err.Error())
		completionEvent.Add(1)
		go sub.Start(completionEvent)
		return
	}
	defer client.Close()
	sub.Logger.Info("Started Ethereum websocket with provider:", sub.EthProvider)

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		sub.Logger.Error(err.Error())
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
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(client, sub.RegistryContractAddress, txs.BridgeBank)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
	}

	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)

	// Listen the new header
	heads := make(chan *ctypes.Header)
	defer close(heads)
	subHead, err := client.SubscribeNewHead(context.Background(), heads)
	if err != nil {
		log.Println(err)
		return
	}
	defer subHead.Unsubscribe()

	db, err := leveldb.OpenFile(levelDbFile, nil)
	if err != nil {
		log.Println("Error opening leveldb: ", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("db.Close filed: ", err.Error())
		}
	}()

	ethLevelDBKey := "ethereumLastProcessedBlock"
	var lastProcessedBlock *big.Int

	data, err := db.Get([]byte(ethLevelDBKey), nil)
	if err != nil {
		log.Println("Error getting the last ethereum block from level db", err)
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
			sub.Logger.Error("subHead failed: ", err.Error())
			completionEvent.Add(1)
			go sub.Start(completionEvent)
			return
		case newHead := <-heads:
			sub.Logger.Info(fmt.Sprintf("New header %d with hash %v", newHead.Number, newHead.Hash()))

			startingBigInt := newHead.Number
			endingBlock := startingBigInt.Sub(startingBigInt, big.NewInt(trailingBlocks))

			// if the current block number - trailing blocks is negative, don't bother
			// going deeper into the function.
			if endingBlock.Cmp(big.NewInt(0)) == -1 {
				log.Println("Ending block index negative. Cancelling run")
				continue
			}

			// If the last processed block is the default (0), then go and set it to the difference of ending block minus 1
			// The user who starts this must provide a valid last processed block
			if lastProcessedBlock.Cmp(big.NewInt(0)) == 0 {
				lastProcessedBlock.Sub(endingBlock, big.NewInt(1))
			}

			sub.Logger.Info(fmt.Sprintf("Processing events from block %d to %d", lastProcessedBlock, endingBlock))

			// query event data from this specific block range
			ethLogs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
				FromBlock: lastProcessedBlock,
				ToBlock:   endingBlock,
				Addresses: []common.Address{bridgeBankAddress},
			})
			if err != nil {
				log.Printf("Error getting events on block %d from bridgebank: %v", newHead.Number, err)
				// if you have an error getting the logs from the block, continue and keep
				// the current last processed block so we keep retrying
				continue
			}

			// Assumption here is that we will repeat a failing block because we return if there is an error retrieving logs
			log.Printf("Successfully received bridgebank events from block %d to %d ", lastProcessedBlock, endingBlock)

			var events []types.EthereumEvent

			// loop over ethlogs, and build an array of burn/lock events
			for _, ethLog := range ethLogs {
				log.Printf("Processed events from block %v", ethLog.BlockNumber)
				event, isBurnLock, err := sub.logToEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI, ethLog)
				if err != nil {
					log.Println("Continuing processing events: ", err)
					continue
				}
				if !isBurnLock {
					log.Println("not burn or lock event, continue events: ", err)
					continue
				}
				events = append(events, event)
			}

			if len(events) > 0 {
				if err := sub.handleEthereumEvent(events); err != nil {
					log.Println("handleEthereumEvent failed: ", err.Error())
				}
				time.Sleep(transactionInterval)
			}
			// save the current ending block to the lastprocessed block to ensure we keep reading blocks sequentially
			lastProcessedBlock = endingBlock
			err = db.Put([]byte(ethLevelDBKey), lastProcessedBlock.Bytes(), nil)
			if err != nil {
				// if you can't write to leveldb, then error out as something is seriously amiss
				log.Fatalf("Error saving lastProcessedBlock to leveldb: %v", err)
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
		sub.Logger.Error(err.Error())
		return event, false, err
	}
	event.BridgeContractAddress = contractAddress
	event.EthereumChainID = clientChainID
	if eventName == types.LogBurn.String() {
		event.ClaimType = ethbridge.BurnText
	} else {
		event.ClaimType = ethbridge.LockText
	}
	sub.Logger.Info(event.String())

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
			sub.Logger.Info(err.Error())
		} else {
			prophecyClaims = append(prophecyClaims, prophecyClaim)
		}
	}
	fmt.Println("prophecyClaims length: ", len(prophecyClaims))

	if len(events) == 0 {
		return nil
	}

	return txs.RelayToCosmos(sub.Cdc, sub.ValidatorName, sub.TempPassword, prophecyClaims, sub.CliCtx, sub.TxBldr)
}
