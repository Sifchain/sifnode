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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sethvargo/go-password/password"
	"github.com/tendermint/go-amino"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

const (
	transactionInterval = 10 * time.Second
)

// TODO: Move relay functionality out of EthereumSub into a new Relayer parent struct

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	Cdc                     *codec.Codec
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
		TmProvider:              rpcURL,
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
		sub.Logger.Error(err.Error())
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

	// Start BridgeBank subscription, prepare contract ABI and LockLog event signature
	bridgeBankAddress, subBridgeBank := sub.startContractEventSub(logs, client, txs.BridgeBank)
	defer subBridgeBank.Unsubscribe()
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

	for {
		select {
		// Handle any errors
		case <-quit:
			return
		case err := <-subBridgeBank.Err():
			sub.Logger.Error(err.Error())
			completionEvent.Add(1)
			go sub.Start(completionEvent)
			return
		case err := <-subHead.Err():
			sub.Logger.Error(err.Error())
			completionEvent.Add(1)
			go sub.Start(completionEvent)
			return
		case newHead := <-heads:
			sub.Logger.Info(fmt.Sprintf("New header %d with hash %v", newHead.Number, newHead.Hash()))

			// Add new header info to buffer
			sub.EventsBuffer.AddHeader(newHead.Number, newHead.Hash(), newHead.ParentHash)
			for {
				fifty := big.NewInt(50)
				fifty.Add(fifty, sub.EventsBuffer.MinHeight)
				if fifty.Cmp(newHead.Number) <= 0 {
					events := sub.EventsBuffer.GetHeaderEvents()
					for _, event := range events {
						err := sub.handleEthereumEvent(event)
						time.Sleep(transactionInterval)
						if err != nil {
							sub.Logger.Error(err.Error())
							completionEvent.Add(1)
						}
					}
					sub.EventsBuffer.RemoveHeight()
				} else {
					break
				}
			}

		// vLog is raw event data
		case vLog := <-logs:
			sub.Logger.Info(fmt.Sprintf("Witnessed tx %s on block %d\n", vLog.TxHash.Hex(), vLog.BlockNumber))
			event, isBurnLock, err := sub.logToEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI, vLog)
			if err != nil {
				sub.Logger.Error("Failed to get event from ethereum log")
			} else if isBurnLock {
				sub.Logger.Info("Add event into buffer")
				sub.EventsBuffer.AddEvent(big.NewInt(int64(vLog.BlockNumber)), vLog.BlockHash, event)
			}
		}
	}
}

func (sub EthereumSub) getAllClaims(fromBlock int64, toBlock int64) []types.EthereumBridgeClaim {
	sub.Logger.Info(fmt.Sprintf("Replay get all ethereum bridge claim from block %d to block %d\n", fromBlock, toBlock))

	var claimArray []types.EthereumBridgeClaim
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.Logger.Error("failed to initialize a client", "err", err)
		return claimArray
	}
	client.SetLogger(sub.Logger)

	if err := client.Start(); err != nil {
		sub.Logger.Error("failed to start a client", "err", err)
		return claimArray
	}

	defer client.Stop() //nolint:errcheck

	for blockNumber := fromBlock; blockNumber < toBlock; {
		tmpBlockNumber := blockNumber
		block, err := client.BlockResults(&tmpBlockNumber)
		blockNumber++
		sub.Logger.Info(fmt.Sprintf("Replay start to process block %d", blockNumber))

		if err != nil {
			sub.Logger.Error(fmt.Sprintf("failed to start a client %s", err))
			continue
		}

		for _, log := range block.TxsResults {
			for _, event := range log.Events {
				sub.Logger.Info(fmt.Sprintf("Replay get an event %s", event.GetType()))
				if event.GetType() == "create_claim" {
					claim, err := txs.AttributesToEthereumBridgeClaim(event.GetAttributes())
					if err != nil {
						continue
					}

					// Check if sender is me
					if claim.CosmosSender.Equals(sub.ValidatorAddress) {
						sub.Logger.Info(fmt.Sprintf("We got an eth bridge claim message %s", claim.EthereumSender.String()))
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
	sub.Logger.Info(fmt.Sprintf("ethereum replay for %d block to %d block\n", fromBlock, toBlock))

	bridgeClaims := sub.getAllClaims(cosmosFromBlock, cosmosToBlock)
	sub.Logger.Info(fmt.Sprintf("found out %d bridgeClaims\n", len(bridgeClaims)))

	client, err := SetupRPCEthClient(sub.EthProvider)
	if err != nil {
		sub.Logger.Error(err.Error())
		return
	}
	defer client.Close()

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		sub.Logger.Error(err.Error())
		return
	}

	// Get the contract address for this subscription
	subContractAddress, err := txs.GetAddressFromBridgeRegistry(client, sub.RegistryContractAddress, txs.BridgeBank)
	if err != nil {
		sub.Logger.Error(err.Error())
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
		sub.Logger.Error(err.Error())
		return
	}

	for _, log := range logs {
		// fmt.Printf("log is %v", log)
		// Before deal with it, we need check in cosmos if it is already handled by myself bofore.
		event, isBurnLock, err := sub.logToEvent(clientChainID, subContractAddress, bridgeBankContractABI, log)
		if err != nil {
			sub.Logger.Error("Failed to get event from ethereum log")
		} else if isBurnLock {
			sub.Logger.Info(fmt.Sprintf("found out a burn lock event\n"))
			if !EventProcessed(bridgeClaims, event) {
				err := sub.handleEthereumEvent(event)
				time.Sleep(transactionInterval)
				if err != nil {
					sub.Logger.Error(err.Error())
				}
			} else {
				sub.Logger.Info(fmt.Sprintf("event already processed by me\n"))
			}
		}
	}

}

// startContractEventSub : starts an event subscription on the specified Peggy contract
func (sub EthereumSub) startContractEventSub(logs chan ctypes.Log, client *ethclient.Client,
	contractName txs.ContractRegistry) (common.Address, ethereum.Subscription) {
	// Get the contract address for this subscription
	subContractAddress, err := txs.GetAddressFromBridgeRegistry(client, sub.RegistryContractAddress, contractName)
	if err != nil {
		sub.Logger.Error(err.Error())
	}

	// We need the address in []bytes for the query
	subQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subContractAddress},
	}

	// Start the contract subscription
	contractSub, err := client.SubscribeFilterLogs(context.Background(), subQuery, logs)
	if err != nil {
		sub.Logger.Error(err.Error())
	}
	sub.Logger.Info(fmt.Sprintf("Subscribed to %v contract at address: %s", contractName, subContractAddress.Hex()))
	return subContractAddress, contractSub
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
func (sub EthereumSub) handleEthereumEvent(event types.EthereumEvent) error {
	prophecyClaim, err := txs.EthereumEventToEthBridgeClaim(sub.ValidatorAddress, event)
	if err != nil {
		sub.Logger.Info(err.Error())
		return err
	}

	return txs.RelayToCosmos(sub.Cdc, sub.ValidatorName, sub.TempPassword, &prophecyClaim, sub.CliCtx, sub.TxBldr)
}
