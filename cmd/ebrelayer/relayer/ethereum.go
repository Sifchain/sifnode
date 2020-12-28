package relayer

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmLog "github.com/tendermint/tendermint/libs/log"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	EthProvider             string
	RegistryContractAddress common.Address
	PrivateKey              *ecdsa.PrivateKey
	EventsBuffer            types.EthEventBuffer
	CosmosContext           *types.CosmosContext
	Logger                  tmLog.Logger
}

// NewEthereumSub initializes a new EthereumSub
func NewEthereumSub(ethProvider string,
	registryContractAddress common.Address, privateKey *ecdsa.PrivateKey, cosmosContext *types.CosmosContext, logger tmLog.Logger) (EthereumSub, error) {

	return EthereumSub{
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		PrivateKey:              privateKey,
		EventsBuffer:            types.NewEthEventBuffer(),
		CosmosContext:           cosmosContext,
		Logger:                  logger,
	}, nil
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
	prophecyClaim, err := txs.EthereumEventToEthBridgeClaim(sub.CosmosContext.ValidatorAddress, &event)
	if err != nil {
		return err
	}
	return txs.RelayToCosmos(sub.CosmosContext, &prophecyClaim)
}
