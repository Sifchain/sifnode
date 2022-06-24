package relayer

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sifchain/sifnode/x/instrumentation"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/artifacts/contracts/CosmosBridge.sol"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"google.golang.org/grpc"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"

	"go.uber.org/zap"
)

const (
	errorMessageKey      = "errorMessage"
	cosmosSleepDuration  = 1
	maxCosmosQueryBlocks = 5000
	// ProphecyLifeTime signature info life time on chain
	blockTimeInSecond = 5
	secondsPerDay     = 60 * 60 * 24
	daysPerMonth      = 30
	ProphecyLifeTime  = (secondsPerDay * daysPerMonth) / blockTimeInSecond
)

// CosmosSub defines a Cosmos listener that relays events to Ethereum and Cosmos
type CosmosSub struct {
	TmProvider              string
	EthProvider             string
	PrivateKey              *ecdsa.PrivateKey
	SugaredLogger           *zap.SugaredLogger
	NetworkDescriptor       oracletypes.NetworkDescriptor
	RegistryContractAddress common.Address
	CliContext              client.Context
	ValidatorName           string
	MaxFeePerGas            *big.Int
	MaxPriorityFeePerGas    *big.Int
	EthereumChainId         *big.Int
	SifnodeGrpc             string
}

// NewCosmosSub initializes a new CosmosSub
func NewCosmosSub(networkDescriptor oracletypes.NetworkDescriptor,
	privateKey *ecdsa.PrivateKey,
	tmProvider,
	ethProvider string,
	registryContractAddress common.Address,
	cliContext client.Context,
	validatorName string,
	sugaredLogger *zap.SugaredLogger,
	maxFeePerGas,
	maxPriorityFeePerGas,
	ethereumChainId *big.Int,
	sifnodeGrpc string) CosmosSub {

	return CosmosSub{
		NetworkDescriptor:       networkDescriptor,
		TmProvider:              tmProvider,
		PrivateKey:              privateKey,
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		CliContext:              cliContext,
		ValidatorName:           validatorName,
		SugaredLogger:           sugaredLogger,
		MaxFeePerGas:            maxFeePerGas,
		MaxPriorityFeePerGas:    maxPriorityFeePerGas,
		EthereumChainId:         ethereumChainId,
		SifnodeGrpc:             sifnodeGrpc,
	}
}

// Start a Cosmos chain subscription
func (sub CosmosSub) Start(txFactory tx.Factory, completionEvent *sync.WaitGroup, symbolTranslator *symbol_translator.SymbolTranslator) {
	defer completionEvent.Done()
	time.Sleep(time.Second)
	client, err := tmclient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		return
	}

	if err := client.Start(); err != nil {
		sub.SugaredLogger.Errorw("failed to start a sifchain client.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent, symbolTranslator)
		return
	}

	defer client.Stop() //nolint:errcheck

	sub.SugaredLogger.Info("Initialized cosmos subscription")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	for {
		select {
		// Handle any errors
		case <-quit:
			log.Println("we receive the quit signal and exit")
			return
		default:
			sub.CheckSequenceAndProcess(txFactory, client)
			time.Sleep(time.Second * cosmosSleepDuration)
		}
	}
}

// CheckSequenceAndProcess check the lock burn Sequence and process the event
func (sub CosmosSub) CheckSequenceAndProcess(txFactory tx.Factory,
	client *tmclient.HTTP) {

	valAddr, err := GetValAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the validator address from validator moniker",
			errorMessageKey, err.Error())
		return
	}

	// get lock burn Sequence and start block number from cosmos
	globalSequence, blockNumber, err := sub.GetGlobalSequenceBlockNumberFromCosmos(sub.NetworkDescriptor, valAddr.String())
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the lock burn Sequence from cosmos rpc",
			errorMessageKey, err.Error())
		return
	}

	// If block number for the next global sequence is zero, means no new lock/burn transaction happened in Sifnode side.
	// just return here, to avoid the unnecessary events querying from block 0 to current block.
	if blockNumber == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cosmosSleepDuration)*time.Second)
	defer cancel()
	block, err := client.Block(ctx, nil)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the block via cosmos client",
			errorMessageKey, err.Error())
		return
	}
	currentBlockHeight := uint64(block.Block.Header.Height)
	if currentBlockHeight-blockNumber > maxCosmosQueryBlocks {
		currentBlockHeight = blockNumber + maxCosmosQueryBlocks

	}
	sub.ProcessLockBurnWithScope(txFactory, client, globalSequence, blockNumber, currentBlockHeight)
}

// ProcessLockBurnWithScope scan blocks in scope and handle all burn lock events
func (sub CosmosSub) ProcessLockBurnWithScope(txFactory tx.Factory, client *tmclient.HTTP, globalSequence, fromBlockNumber, toBlockNumber uint64) {
	sub.SugaredLogger.Infow("ProcessLockBurnWithScope",
		"globalSequence", globalSequence,
		"fromBlockNumber", fromBlockNumber,
		"toBlockNumber", toBlockNumber)

	for blockNumber := fromBlockNumber; blockNumber <= toBlockNumber; blockNumber++ {
		tmpBlockNumber := int64(blockNumber)

		ctx := context.Background()
		block, err := client.BlockResults(ctx, &tmpBlockNumber)

		if err != nil {
			sub.SugaredLogger.Infow("sifchain client failed to get a block.",
				errorMessageKey, err.Error())
			continue
		}

		for _, txLog := range block.TxsResults {
			for _, event := range txLog.Events {

				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.MsgBurn, types.MsgLock:
					sub.SugaredLogger.Infow("claimtype cosmos.go: ", "claimType: ", claimType)

					// the relayer for signature aggregator not handle burn and lock
					cosmosMsg, err := txs.BurnLockEventToCosmosMsg(event.GetAttributes(), sub.SugaredLogger)
					if err != nil {
						sub.SugaredLogger.Errorw("sifchain client receive a malformed burn/lock message.",
							errorMessageKey, err.Error())
						/**
						When fail to parse the event data, only thing relayer can do just skip it.
						If the wrong event's global sequence is not we expected, then no impact for relayer
						to process next event. But if glocal sequence is next one,
						this continue would cause the bridge to be stuck, since sifnode requires each relayer
						strictly submit the prophecy according to global sequence
						This is same as skipping a msg, the next cosmosMsg.GlobalSequence
						will be +2 of current global sequence. then only solution is upgrade the relayer to
						the version can handle this specific event.
						*/
						panicString := fmt.Sprintf("Could not create BurnLockEventToCosmos Message, IF YOU SEE THIS FILE A BUG REPORT: %s", err.Error())
						panic(panicString)
					}

					sub.SugaredLogger.Infow(
						"Received message from sifchain: ",
						"msg", cosmosMsg,
						"NetworkDescriptor", cosmosMsg.NetworkDescriptor,
						"GlobalSequence", globalSequence,
					)

					if cosmosMsg.NetworkDescriptor == sub.NetworkDescriptor {
						// if global Sequence is expected, sign prophecy and send back to cosmos
						// if global Sequence is less than expected, just ignore the event. it is normal to see processed Sequence coexist with expected one
						// if global Sequence is larger than expected, it is wrong and we must miss something.
						if cosmosMsg.GlobalSequence == globalSequence+1 {
							instrumentation.PeggyCheckpointZap(sub.SugaredLogger,
								instrumentation.ReceiveCosmosBurnMessage,
								zap.Reflect("cosmosMsg", cosmosMsg),
								zap.Reflect("sub", sub),
								"globalSequence", globalSequence)

							sub.witnessSignAndBroadcastProphecy(txFactory, cosmosMsg)
							// update expected global Sequence
							globalSequence++

						} else if cosmosMsg.GlobalSequence > globalSequence+1 {
							sub.SugaredLogger.Errorw(
								"The global Sequence is invalid",
								"expected global Sequence is:", globalSequence+1,
								"global Sequence from message is:", cosmosMsg.GlobalSequence,
							)
							return
						}
					}
				}
			}
		}

	}
}

// getOracleClaimType sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) types.Event {
	var claimType types.Event
	switch eventType {
	case types.MsgBurn.String():
		claimType = types.MsgBurn
	case types.MsgLock.String():
		claimType = types.MsgLock
	case types.ProphecyCompleted.String():
		claimType = types.ProphecyCompleted
	default:
		claimType = types.Unsupported
	}
	return claimType
}

func tryInitRelayConfig(sub CosmosSub) (*ethclient.Client, *bind.TransactOpts, common.Address, error) {
	client, auth, target, err := txs.InitRelayConfig(
		sub.EthProvider,
		sub.RegistryContractAddress,
		sub.PrivateKey,
		sub.MaxFeePerGas,
		sub.MaxPriorityFeePerGas,
		sub.EthereumChainId,
		sub.SugaredLogger,
	)

	if err != nil {
		sub.SugaredLogger.Errorw("tryInitRelayConfig", errorMessageKey, err.Error())
		return nil, nil, common.Address{}, errors.New("tryInitRelayConfig")
	}

	return client, auth, target, err
}

// witness node sign against prophecyID of lock and burn message and send the singnature in message back to Sifnode.
func (sub CosmosSub) witnessSignAndBroadcastProphecy(
	txFactory tx.Factory,
	cosmosMsg types.CosmosMsg,
) {
	/**
	err in this func are ignored, and they would cause an issue because we
	assume any cosmos msgs would be processed in here, since globalSequence++ is done in any case
	*/
	sub.SugaredLogger.Infow("handle burn lock message.",
		"cosmosMessage", cosmosMsg.String())

	sub.SugaredLogger.Infow(
		"get the prophecy claim.",
		"cosmosMsg", cosmosMsg,
	)

	valAddr, err := GetValAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		// err is being thrown silently.
		sub.SugaredLogger.Infow(
			"get the prophecy claim.",
			"cosmosMsg", err,
		)
	}

	signData := txs.PrefixMsg(cosmosMsg.ProphecyID)
	address := crypto.PubkeyToAddress(sub.PrivateKey.PublicKey)
	signature, err := txs.SignClaim(signData, sub.PrivateKey)
	if err != nil {
		sub.SugaredLogger.Infow(
			"failed to sign the prophecy id",
			errorMessageKey, err.Error(),
		)
	}

	signProphecy := ethbridgetypes.NewMsgSignProphecy(valAddr.String(), cosmosMsg.NetworkDescriptor,
		cosmosMsg.ProphecyID, address.String(), string(signature))

	instrumentation.PeggyCheckpointZap(sub.SugaredLogger, instrumentation.WitnessSignProphecy, zap.Reflect("prophecy", signProphecy))

	txs.SignProphecyToCosmos(txFactory, signProphecy, sub.CliContext, sub.SugaredLogger)

}

// GetGlobalSequenceBlockNumberFromCosmos get global Sequence block number via rpc
func (sub CosmosSub) GetGlobalSequenceBlockNumberFromCosmos(
	networkDescriptor oracletypes.NetworkDescriptor,
	relayerValAddress string) (uint64, uint64, error) {

	gRpcClientConn, err := grpc.Dial(sub.SifnodeGrpc, grpc.WithInsecure())
	if err != nil {
		return 0, 0, err
	}
	defer gRpcClientConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := ethbridgetypes.NewQueryClient(gRpcClientConn)

	// Get lockburn sequence per networkdescriptor+relayer address
	witnessLockBurnSequenceRequest := ethbridgetypes.QueryWitnessLockBurnSequenceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
	response, err := client.WitnessLockBurnSequence(ctx, &witnessLockBurnSequenceRequest)
	if err != nil {
		return 0, 0, err
	}
	globalSequence := response.WitnessLockBurnSequence

	sub.SugaredLogger.Debugw("Retrieved witness lockburn sequence", "globalSequence", globalSequence)

	// Get the block number of the global sequence to be processed next
	sequenceToBlockNumberRequest := ethbridgetypes.QueryGlobalSequenceBlockNumberRequest{
		NetworkDescriptor: networkDescriptor,
		GlobalSequence:    globalSequence + 1,
	}

	response2, err := client.GlobalSequenceBlockNumber(ctx, &sequenceToBlockNumberRequest)
	if err != nil {
		return 0, 0, err
	}

	sub.SugaredLogger.Debugw("Retrieved block number", "response", &response2)

	return globalSequence, response2.BlockNumber, nil
}

// GetLastNonceSubmitted get last nonce submitted in cosmos bridge contract
func GetLastNonceSubmitted(client *ethclient.Client, cosmosBridgeAddress common.Address, sugaredLogger *zap.SugaredLogger) (*big.Int, error) {

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(cosmosBridgeAddress, client)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return nil, err
	}
	return cosmosBridgeInstance.LastNonceSubmitted(nil)
}

// GetAccAddressFromKeyring get the address from key ring and keyname
func GetAccAddressFromKeyring(k keyring.Keyring, keyname string) (sdk.AccAddress, error) {
	keyInfo, err := k.Key(keyname)
	if err != nil {
		return nil, err
	}
	return keyInfo.GetAddress(), nil
}
