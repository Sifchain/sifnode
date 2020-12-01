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

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

// TODO: Move relay functionality out of EthereumSub into a new Relayer parent struct

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
	Logger                  tmLog.Logger
}

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
	eventLogLockSignature := bridgeBankContractABI.Events[types.LogLock.String()].ID().Hex()
	eventLogBurnSignature := bridgeBankContractABI.Events[types.LogBurn.String()].ID().Hex()

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
		// vLog is raw event data
		case vLog := <-logs:
			sub.Logger.Info(fmt.Sprintf("Witnessed tx %s on block %d\n", vLog.TxHash.Hex(), vLog.BlockNumber))
			log.Println("Found event from the ethereum bridgebank contract: ", types.LogLock.String())
			var err error
			switch vLog.Topics[0].Hex() {
			case eventLogBurnSignature:
				err = sub.handleEthereumEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI,
					types.LogBurn.String(), vLog)
			case eventLogLockSignature:
				log.Println("Found loglock event from the ethereum bridgebank contract: ", types.LogLock.String())
				err = sub.handleEthereumEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI,
					types.LogLock.String(), vLog)
			}
			// TODO: Check local events store for status, if retryable, attempt relay again
			if err != nil {
				sub.Logger.Error(err.Error())
				completionEvent.Add(1)
				go sub.Start(completionEvent)
				return
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

// handleEthereumEvent unpacks an Ethereum event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func (sub EthereumSub) handleEthereumEvent(clientChainID *big.Int, contractAddress common.Address,
	contractABI abi.ABI, eventName string, cLog ctypes.Log) error {
	// Parse the event's attributes via contract ABI
	event := types.EthereumEvent{}
	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.Logger.Error("error unpacking: %v", err)
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

	prophecyClaim, err := txs.EthereumEventToEthBridgeClaim(sub.ValidatorAddress, &event)
	if err != nil {
		return err
	}
	return txs.RelayToCosmos(sub.Cdc, sub.ValidatorName, sub.TempPassword, &prophecyClaim, sub.CliCtx, sub.TxBldr)
}
