package relayer

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	bridgeBankContract "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/artifacts/contracts/BridgeBank/BridgeBank.sol"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"github.com/Sifchain/sifnode/x/instrumentation"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	trailingBlocks        = 50
	ethereumSleepDuration = 1
	maxQueryBlocks        = 5000
)

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	EthProvider             string
	TmProvider              string
	RegistryContractAddress common.Address
	ValidatorName           string
	ValidatorAddress        sdk.ValAddress
	NetworkDescriptor       oracletypes.NetworkDescriptor
	CliCtx                  client.Context
	PrivateKey              *ecdsa.PrivateKey
	SugaredLogger           *zap.SugaredLogger
	SifnodeGrpc             string
}

// NewKeybase create a new keybase instance
func NewKeybase(validatorMoniker, mnemonic, password string) (keyring.Keyring, keyring.Info, error) {
	kr := keyring.NewInMemory()
	hdpath := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	info, err := kr.NewAccount(validatorMoniker, mnemonic, password, hdpath.String(), hd.Secp256k1)
	if err != nil {
		return nil, nil, err
	}

	return kr, info, nil
}

// NewEthereumSub initializes a new EthereumSub
func NewEthereumSub(
	cliCtx client.Context,
	nodeURL string,
	validatorMoniker string,
	networkDescriptor oracletypes.NetworkDescriptor,
	ethProvider string,
	registryContractAddress common.Address,
	sugaredLogger *zap.SugaredLogger,
	sifnodeGrpc string,
) EthereumSub {

	return EthereumSub{
		EthProvider:             ethProvider,
		TmProvider:              nodeURL,
		NetworkDescriptor:       networkDescriptor,
		RegistryContractAddress: registryContractAddress,
		ValidatorName:           validatorMoniker,
		ValidatorAddress:        nil,
		CliCtx:                  cliCtx,
		SugaredLogger:           sugaredLogger,
		SifnodeGrpc:             sifnodeGrpc,
	}
}

// Start an Ethereum chain subscription
func (sub EthereumSub) Start(txFactory tx.Factory,
	completionEvent *sync.WaitGroup,
	symbolTranslator *symbol_translator.SymbolTranslator) {

	defer completionEvent.Done()
	ethClient, err := SetupWebsocketEthClient(sub.EthProvider)

	if err != nil {
		sub.SugaredLogger.Errorw("SetupWebsocketEthClient failed.",
			errorMessageKey, err.Error())

		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent, symbolTranslator)
		return
	}
	defer ethClient.Close()
	sub.SugaredLogger.Infow("Started Ethereum websocket with provider:",
		"Ethereum provider", sub.EthProvider)

	tmClient, err := tmclient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		return
	}

	networkID := big.NewInt(int64(sub.NetworkDescriptor))

	validatorAddress, err := GetValAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		log.Fatal("Error getting validator address: ", err.Error())
	}

	sub.ValidatorAddress = validatorAddress

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	// get the bridgebank address from the registry contract
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
	}

	for {
		select {
		// Handle any errors
		case <-quit:
			return
		default:
			if !sub.CheckNonceAndProcess(txFactory,
				networkID,
				ethClient,
				tmClient,
				bridgeBankAddress,
				symbolTranslator) {
				// CheckNonceAndProcess did no work, so we pause for a bit
				time.Sleep(time.Second * ethereumSleepDuration)
			}
		}
	}
}

// CheckNonceAndProcess check the lock burn nonce and process the event
func (sub EthereumSub) CheckNonceAndProcess(txFactory tx.Factory,
	networkID *big.Int,
	ethClient *ethclient.Client,
	tmClient *tmclient.HTTP,
	bridgeBankAddress common.Address,
	symbolTranslator *symbol_translator.SymbolTranslator) (processedBlocks bool) {
	// get current block height
	currentBlock, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the current block from ethereum client",
			errorMessageKey, err.Error())
		return false
	}
	currentBlockHeight := currentBlock.Number

	// If current block is less than 50, just return.
	if currentBlockHeight.Cmp(big.NewInt(trailingBlocks)) <= 0 {
		return
	}

	endBlockHeight := big.NewInt(0)
	endBlockHeight = endBlockHeight.Sub(currentBlockHeight, big.NewInt(trailingBlocks))

	// get lock burn nonce from cosmos
	lockBurnSequence, err := sub.GetLockBurnSequenceFromCosmos(oracletypes.NetworkDescriptor(networkID.Uint64()), sub.ValidatorAddress.String())
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the lock burn nonce from cosmos rpc",
			errorMessageKey, err.Error())
		return
	}

	topics := [][]common.Hash{}
	// add the log type as first topic, the first search filter will be lockTopic or burnTopic
	abi, err := bridgeBankContract.BridgeBankMetaData.GetAbi()
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get events from abi")
		return
	}
	lockTopic := abi.Events[types.LogLock.String()].ID
	burnTopic := abi.Events[types.LogBurn.String()].ID
	topics = append(topics, []common.Hash{lockTopic, burnTopic})

	// add the lock burn nonce as second topic, combined search filter will be (lockTopic or burnTopic)
	// and lockBurnNonceTopic
	var lockBurnNonceTopic [32]byte
	bigLockBurnSequence := (&big.Int{}).SetUint64(lockBurnSequence + 1)
	paddedLockBurnSequence := math.PaddedBigBytes(bigLockBurnSequence, 32)
	copy(lockBurnNonceTopic[:], paddedLockBurnSequence[:32])
	topics = append(topics, []common.Hash{lockBurnNonceTopic})

	// query the exact block number with the lock burn nonce
	filterQuery := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   endBlockHeight,
		Addresses: []common.Address{bridgeBankAddress},
		Topics:    topics,
	}
	ethLogs, err := ethClient.FilterLogs(context.Background(), filterQuery)

	if err != nil {
		sub.SugaredLogger.Errorw("failed to filter the logs from ethereum client",
			errorMessageKey, err.Error())
		return
	}

	fromBlockNumber := uint64(0)
	lenEthLogs := len(ethLogs)
	if lenEthLogs != 1 {
		sub.SugaredLogger.Debugw("no results from filter", "lenEthLogs", lenEthLogs)
		return
	}

	bridgeBankInstance, err := bridgeBankContract.NewBridgeBank(bridgeBankAddress, ethClient)
	if err != nil {
		sub.SugaredLogger.Errorw("NewBridgeBank",
			errorMessageKey, err.Error())
		return
	}

	event, isBurnLock, err := sub.logToEvent(oracletypes.NetworkDescriptor(networkID.Uint64()),
		bridgeBankAddress,
		bridgeBankInstance,
		ethLogs[0])

	if err != nil {
		sub.SugaredLogger.Errorw("failed to transform from log to event.",
			errorMessageKey, err.Error())
		return
	}
	if !isBurnLock {
		sub.SugaredLogger.Infow("not burn or lock event, continue events.")
		return
	}

	if event.Nonce.Uint64() != lockBurnSequence+1 {
		sub.SugaredLogger.Errorw("the lock burn nonce is not expected.")
		return
	}

	// get the block height for the specific lock burn nonce
	fromBlockNumber = ethLogs[0].BlockNumber

	events := []types.EthereumEvent{}
	// get a new topics, exclude the lock burn nonce since we already get block number
	topics = [][]common.Hash{}
	topics = append(topics, []common.Hash{lockTopic, burnTopic})

	for endBlock := endBlockHeight.Uint64(); fromBlockNumber <= endBlock; endBlock = endBlockHeight.Uint64() {

		// query block scope limited to maxQueryBlocks
		if endBlock > fromBlockNumber+maxQueryBlocks {
			endBlock = fromBlockNumber + maxQueryBlocks
		}

		// query the events with block scope
		ethLogs, err = ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlockNumber)),
			ToBlock:   big.NewInt(int64(endBlock)),
			Addresses: []common.Address{bridgeBankAddress},
			Topics:    topics,
		})

		if err != nil {
			sub.SugaredLogger.Errorw("failed to filter the logs from ethereum client",
				errorMessageKey, err.Error())
			return
		}

		// loop over ethlogs, and build an array of burn/lock events
		for _, ethLog := range ethLogs {
			event, isBurnLock, err := sub.logToEvent(oracletypes.NetworkDescriptor(networkID.Uint64()),
				bridgeBankAddress,
				bridgeBankInstance,
				ethLog)

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
			if lockBurnSequence, err = sub.handleEthereumEvent(txFactory, events, symbolTranslator, lockBurnSequence); err != nil {
				sub.SugaredLogger.Errorw("failed to handle ethereum event.",
					errorMessageKey, err.Error())
				return
			}
			processedBlocks = true
		}

		// update fromBlockNumber
		fromBlockNumber += maxQueryBlocks
	}
	return processedBlocks
}

// Replay the missed events
func (sub EthereumSub) Replay(txFactory tx.Factory, symbolTranslator *symbol_translator.SymbolTranslator) {

	ethClient, err := SetupWebsocketEthClient(sub.EthProvider)
	if err != nil {
		sub.SugaredLogger.Errorw("SetupWebsocketEthClient failed.",
			errorMessageKey, err.Error())

		return
	}
	defer ethClient.Close()
	sub.SugaredLogger.Infow("Started Ethereum websocket with provider:",
		"Ethereum provider", sub.EthProvider)

	tmClient, err := tmclient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		return
	}

	networkID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get network ID.",
			errorMessageKey, err.Error())

		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	// get the bridgebank address from the registry contract
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
	}

	sub.CheckNonceAndProcess(txFactory,
		networkID,
		ethClient,
		tmClient,
		bridgeBankAddress,
		symbolTranslator)
}

// logToEvent unpacks an Ethereum event
func (sub EthereumSub) logToEvent(networkDescriptor oracletypes.NetworkDescriptor,
	contractAddress common.Address,
	bridgeBank *bridgeBankContract.BridgeBank,
	cLog ctypes.Log) (types.EthereumEvent, bool, error) {
	// Parse the event's attributes via contract ABI
	event := types.EthereumEvent{}
	event.BridgeContractAddress = contractAddress

	if decodedEvent, err := bridgeBank.BridgeBankFilterer.ParseLogLock(cLog); err == nil {
		event.ClaimType = ethbridgetypes.ClaimType_CLAIM_TYPE_LOCK
		event.To = append(event.To, decodedEvent.To...)
		event.Symbol = decodedEvent.Symbol
		event.Name = decodedEvent.Name
		event.Decimals = decodedEvent.Decimals
		event.NetworkDescriptor = int32(networkDescriptor)
		event.Value = decodedEvent.Value
		event.Nonce = (&big.Int{}).Set(decodedEvent.Nonce)
		event.From = decodedEvent.From
		event.Token = decodedEvent.Token
	}
	if decodedEvent, err := bridgeBank.BridgeBankFilterer.ParseLogBurn(cLog); err == nil {
		event.ClaimType = ethbridgetypes.ClaimType_CLAIM_TYPE_BURN
		event.From = decodedEvent.From
		event.To = append(event.To, decodedEvent.To...)
		event.CosmosDenom = decodedEvent.Denom
		event.Token = decodedEvent.Token
		event.Value = decodedEvent.Value
		event.Nonce = (&big.Int{}).Set(decodedEvent.Nonce)
		event.Decimals = decodedEvent.Decimals
		event.NetworkDescriptor = int32(networkDescriptor)
	}

	instrumentation.PeggyCheckpointZap(
		sub.SugaredLogger,
		instrumentation.EthereumEvent,
		zap.Reflect("event", event),
		"txhash", cLog.TxHash.Hex(),
	)

	// Add the event to the record
	types.NewEventWrite(cLog.TxHash.Hex(), event)
	return event, true, nil
}

// handleEthereumEvent unpacks an Ethereum event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func (sub EthereumSub) handleEthereumEvent(txFactory tx.Factory,
	events []types.EthereumEvent,
	symbolTranslator *symbol_translator.SymbolTranslator,
	lockBurnNonce uint64) (uint64, error) {

	var ethBridgeClaims []*ethbridgetypes.EthBridgeClaim

	valAddr, err := GetValAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		return lockBurnNonce, err
	}
	for _, event := range events {
		ethBridgeClaim, err := txs.EthereumEventToEthBridgeClaim(valAddr, event, symbolTranslator, sub.SugaredLogger)
		if err != nil {
			sub.SugaredLogger.Errorw(".",
				errorMessageKey, err.Error())
		} else {
			// lockBurnNonce is zero, means the relayer is new one, never process event before
			// then it start from current event and sifnode will accept it
			if lockBurnNonce == 0 || ethBridgeClaim.EthereumLockBurnSequence == lockBurnNonce+1 {
				ethBridgeClaims = append(ethBridgeClaims, &ethBridgeClaim)
				instrumentation.PeggyCheckpointZap(sub.SugaredLogger, instrumentation.EthereumProphecyClaim, zap.Reflect("event", event), "prophecyClaim", ethBridgeClaim)
				lockBurnNonce = ethBridgeClaim.EthereumLockBurnSequence
			} else {
				sub.SugaredLogger.Infow("lock burn nonce is not expected",
					"nextLockBurnNonce", lockBurnNonce,
					"prophecyClaim.EthereumLockBurnNonce", ethBridgeClaim.EthereumLockBurnSequence)
				return lockBurnNonce, errors.New("lock burn nonce is not expected")
			}

		}
	}
	sub.SugaredLogger.Infow("relay prophecy claims to cosmos.",
		"prophecy claims length", len(ethBridgeClaims))

	if len(events) == 0 {
		return lockBurnNonce, nil
	}

	return lockBurnNonce, txs.RelayToCosmos(txFactory, ethBridgeClaims, sub.CliCtx, sub.SugaredLogger)
}

// GetLockBurnNonceFromCosmos via rpc
func (sub EthereumSub) GetLockBurnSequenceFromCosmos(
	networkDescriptor oracletypes.NetworkDescriptor,
	relayerValAddress string) (uint64, error) {

	conn, err := grpc.Dial(sub.SifnodeGrpc, grpc.WithInsecure())
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cosmosSleepDuration)*time.Second*10)
	defer cancel()
	queryClient := ethbridgetypes.NewQueryClient(conn)
	request := ethbridgetypes.QueryEthereumLockBurnSequenceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
	response, err := queryClient.EthereumLockBurnSequence(ctx, &request)
	if err != nil {
		return 0, err
	}
	return response.EthereumLockBurnSequence, nil
}

// GetValAddressFromKeyring get validator address from keyring
func GetValAddressFromKeyring(k keyring.Keyring, keyname string) (sdk.ValAddress, error) {
	keyInfo, err := k.Key(keyname)
	if err != nil {
		return nil, err
	}
	return sdk.ValAddress(keyInfo.GetAddress()), nil
}
