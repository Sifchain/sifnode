package txs

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"time"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit            = uint64(500000)
)

var GasPriceMinimum *big.Int = big.NewInt(60000000000)

func sleepThread(seconds time.Duration) {
	time.Sleep(time.Second * seconds)
}

// RelayProphecyClaimToEthereum relays the provided ProphecyClaim to CosmosBridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(
	claim types.CosmosMsg,
	sugaredLogger *zap.SugaredLogger,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	cosmosBridgeInstance *cosmosbridge.CosmosBridge,
	) error {

	// Send transaction
	sugaredLogger.Infow(
		"Sending new ProphecyClaim to CosmosBridge.",
		"CosmosSender", claim.CosmosSender,
		"CosmosSenderSequence", claim.CosmosSenderSequence,
	)

	amount := claim.Amount.BigInt()

	tx, err := cosmosBridgeInstance.NewProphecyClaim(
		auth,
		uint8(claim.ClaimType),
		claim.CosmosSender,
		claim.CosmosSenderSequence,
		claim.EthereumReceiver,
		claim.Symbol,
		amount,
	)

	// sleep 2 seconds to wait for tx to go through before querying.
	sleepThread(2)

	if err != nil {
		return err
	}

	sugaredLogger.Infow("get NewProphecyClaim tx hash:", "ProphecyClaimHash", tx.Hash().Hex())

	// var receipt *eth.types.Receipt
	var receipt *ctypes.Receipt
	maxRetries := 60
	i := 0
	// if there is an error getting the tx, or if the tx fails, retry 60 times
	for i < maxRetries {
		// Get the transaction receipt
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
	
		if err != nil {
			sleepThread(1)
		} else {
			break
		}
		i++
	}

	if i == maxRetries {
		return errors.New("hit max tx receipt query retries")
	}

	sugaredLogger.Infow(
		"Successfully received transaction receipt after retry",
		"txReceipt", receipt,
	)

	return nil
}

// InitRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func InitRelayConfig(
	provider string,
	registry common.Address,
	event types.Event,
	key *ecdsa.PrivateKey,
	sugaredLogger *zap.SugaredLogger,
	) (
		*ethclient.Client,
		*bind.TransactOpts,
		common.Address,
		error,
	) {
	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		sugaredLogger.Errorw("failed to connect ethereum node.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	// Load the validator's address
	sender, err := LoadSender()
	if err != nil {
		sugaredLogger.Errorw("failed to load validator address.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), sender)
	sugaredLogger.Infow("Current eth operator at pending nonce.", "pendingNonce", nonce)

	if err != nil {
		sugaredLogger.Errorw("failed to broadcast transaction.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		sugaredLogger.Errorw("failed to get gas price.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(key)

	sugaredLogger.Infow("ethereum tx current nonce from client api.",
		"nonce", nonce,
		"suggestedGasPrice", gasPrice)

	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))

	quarterGasPrice := big.NewInt(0)
	quarterGasPrice = quarterGasPrice.Div(gasPrice, big.NewInt(4))

	gasPrice.Sub(gasPrice, quarterGasPrice)
	sugaredLogger.Infow("final gas price after adjustment.",
		"finalGasPrice", gasPrice)

	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	sugaredLogger.Infow("nonce before send transaction.",
		"transactOptsAuth.Nonce", transactOptsAuth.Nonce)

	var targetContract ContractRegistry
	switch event {
	// ProphecyClaims are sent to the CosmosBridge contract
	case types.MsgBurn, types.MsgLock:
		targetContract = CosmosBridge
	// OracleClaims are sent to the Oracle contract
	case types.LogNewProphecyClaim:
		targetContract = Oracle
	default:
		panic("invalid target contract address")
	}

	// Get the specific contract's address
	target, err := GetAddressFromBridgeRegistry(client, registry, targetContract, sugaredLogger)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmos bridger contract address from registry.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err

	}
	return client, transactOptsAuth, target, nil
}
