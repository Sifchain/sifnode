package relayer

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/syndtr/goleveldb/leveldb"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	bridgebank "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/bridgebank"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	transactionInterval = 10 * time.Second
	ethereumWakeupTimer = 60
)

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	EthProvider             string
	TmProvider              string
	RegistryContractAddress common.Address
	ValidatorName           string
	ValidatorAddress        sdk.ValAddress
	CliCtx                  client.Context
	PrivateKey              *ecdsa.PrivateKey
	DB                      *leveldb.DB
	SugaredLogger           *zap.SugaredLogger
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
	validatorMoniker,
	ethProvider string,
	registryContractAddress common.Address,
	validatorAddress sdk.ValAddress,
	db *leveldb.DB,
	sugaredLogger *zap.SugaredLogger,
) EthereumSub {

	return EthereumSub{
		EthProvider:             ethProvider,
		TmProvider:              nodeURL,
		RegistryContractAddress: registryContractAddress,
		ValidatorName:           validatorMoniker,
		ValidatorAddress:        validatorAddress,
		CliCtx:                  cliCtx,
		DB:                      db,
		SugaredLogger:           sugaredLogger,
	}
}

// Start an Ethereum event handler
func (sub EthereumSub) Start(txFactory tx.Factory, completionEvent *sync.WaitGroup) {
	defer completionEvent.Done()
	time.Sleep(time.Second)
	ethClient, err := SetupWebsocketEthClient(sub.EthProvider)
	if err != nil {
		sub.SugaredLogger.Errorw("SetupWebsocketEthClient failed.",
			errorMessageKey, err.Error())

		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent)
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
		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent)
		return
	}

	// We will check logs for new events
	logs := make(chan ctypes.Log)
	defer close(logs)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	// get the bridgebank address from the registry contract
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
		return
	}

	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)

	t := time.NewTicker(time.Second * ethereumWakeupTimer)

	for {
		select {
		// Handle any errors
		case <-quit:
			return
		case <-t.C:
			sub.CheckNonceAndProcess(txFactory, networkID, ethClient, tmClient, bridgeBankAddress, bridgeBankContractABI)
		}
	}
}

// CheckNonceAndProcess check the lock burn nonce and process the event
func (sub EthereumSub) CheckNonceAndProcess(txFactory tx.Factory,
	networkID *big.Int,
	ethClient *ethclient.Client,
	tmClient *tmclient.HTTP,
	bridgeBankAddress common.Address,
	bridgeBankContractABI abi.ABI) {
	// get current block height
	currentBlock, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return
	}
	currentBlockHeight := currentBlock.Number.Uint64()

	lockBurnNonce, err := sub.GetLockBurnNonceFromCosmos(oracletypes.NetworkDescriptor(networkID.Uint64()), string(sub.ValidatorAddress))
	if err != nil {
		return
	}

	fromBlock := sub.GetBlockHeightWithLockBurnNonce(ethClient,
		bridgeBankAddress,
		currentBlockHeight,
		lockBurnNonce)

	sub.HandleEthereumEventWithScope(ethClient,
		bridgeBankAddress,
		bridgeBankContractABI,
		fromBlock,
		currentBlockHeight,
		lockBurnNonce,
		networkID,
		txFactory)
}

// GetBlockHeightWithLockBurnNonce return the block height with specific lock burn nonce
func (sub EthereumSub) GetBlockHeightWithLockBurnNonce(client *ethclient.Client,
	bridgeBankAddress common.Address,
	currentHeight uint64,
	lockBurnNonce uint64) uint64 {

	bridgeBankInstance, err := bridgebank.NewBridgeBank(bridgeBankAddress, client)
	if err != nil {
		return currentHeight
	}

	for currentHeight > 0 {
		callOps := bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: big.NewInt(int64(currentHeight)),
			Context:     nil,
		}

		lockBurnNonceWithHeight, err := bridgeBankInstance.LockBurnNonce(&callOps)
		if err != nil {
			return currentHeight
		}

		// the fist time when get the same nonce, the block height to update nonce should be next one
		if lockBurnNonceWithHeight.Uint64() == lockBurnNonce {
			return currentHeight + 1
		}

		currentHeight--
	}

	return 0
}

// GetLockBurnNonceFromCosmos via rpc
func (sub EthereumSub) GetLockBurnNonceFromCosmos(
	networkDescriptor oracletypes.NetworkDescriptor,
	relayerValAddress string) (uint64, error) {
	conn, err := grpc.Dial(sub.TmProvider)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := ethbridgetypes.NewQueryClient(conn)
	request := ethbridgetypes.QueryLockBurnNonceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
	response, err := client.LockBurnNonce(ctx, &request)
	if err != nil {
		return 0, err
	}
	return response.LockBurnNonce, nil
}

// HandleEthereumEventWithScope process event one by one
func (sub EthereumSub) HandleEthereumEventWithScope(client *ethclient.Client,
	bridgeBankAddress common.Address,
	bridgeBankContractABI abi.ABI,
	fromBlock uint64, toBlock uint64, startLockBurnNonce uint64,
	networkID *big.Int,
	txFactory tx.Factory) {

	events := []types.EthereumEvent{}
	// query event data from this specific block range
	ethLogs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{bridgeBankAddress},
	})

	if err != nil {
		return
	}
	// loop over ethlogs, and build an array of burn/lock events
	for _, ethLog := range ethLogs {
		log.Printf("Processed events from block %v", ethLog.BlockNumber)
		event, isBurnLock, err := sub.logToEvent(oracletypes.NetworkDescriptor(networkID.Uint64()),
			bridgeBankAddress,
			bridgeBankContractABI,
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
		if err := sub.handleEthereumEvent(txFactory, events); err != nil {
			sub.SugaredLogger.Errorw("failed to handle ethereum event.",
				errorMessageKey, err.Error())
		}
		time.Sleep(transactionInterval)
	}
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
func (sub EthereumSub) Replay(txFactory tx.Factory) {
	tmClient, err := tmclient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		return
	}

	ethClient, err := SetupRPCEthClient(sub.EthProvider)
	if err != nil {
		log.Printf("failed to connect ethereum node, error is %s\n", err.Error())
		return
	}
	defer ethClient.Close()

	networkID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		log.Printf("failed to get chain ID, error is %s\n", err.Error())
		return
	}

	// get the bridgebank address from the registry contract
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
		return
	}

	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)

	sub.CheckNonceAndProcess(txFactory, networkID, ethClient, tmClient, bridgeBankAddress, bridgeBankContractABI)
}

// logToEvent unpacks an Ethereum event
func (sub EthereumSub) logToEvent(networkDescriptor oracletypes.NetworkDescriptor, contractAddress common.Address,
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
	event.NetworkDescriptor = networkDescriptor
	if eventName == types.LogBurn.String() {
		event.ClaimType = ethbridgetypes.ClaimType_CLAIM_TYPE_BURN
	} else {
		event.ClaimType = ethbridgetypes.ClaimType_CLAIM_TYPE_LOCK
	}
	sub.SugaredLogger.Infow("receive an event.",
		"event", event.String())

	// Add the event to the record
	types.NewEventWrite(cLog.TxHash.Hex(), event)
	return event, true, nil
}

func GetValAddressFromKeyring(k keyring.Keyring, keyname string) (sdk.ValAddress, error) {
	keyInfo, err := k.Key(keyname)
	if err != nil {
		return nil, err
	}
	return sdk.ValAddress(keyInfo.GetAddress()), nil
}

// handleEthereumEvent unpacks an Ethereum event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func (sub EthereumSub) handleEthereumEvent(txFactory tx.Factory, events []types.EthereumEvent) error {
	var prophecyClaims []*ethbridgetypes.EthBridgeClaim
	valAddr, err := GetValAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		return err
	}
	for _, event := range events {
		prophecyClaim, err := txs.EthereumEventToEthBridgeClaim(valAddr, event)
		if err != nil {
			sub.SugaredLogger.Errorw(".",
				errorMessageKey, err.Error())
		} else {
			prophecyClaims = append(prophecyClaims, &prophecyClaim)
		}
	}
	sub.SugaredLogger.Infow("relay prophecy claims to cosmos.",
		"prophecy claims length", len(prophecyClaims))

	if len(events) == 0 {
		return nil
	}

	return txs.RelayToCosmos(txFactory, prophecyClaims, sub.CliCtx, sub.SugaredLogger)
}
