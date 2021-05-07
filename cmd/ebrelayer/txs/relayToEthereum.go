package txs

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit            = uint64(500000)
	transactionInterval = 60 * time.Second
)

// RelayProphecyClaimToEthereum relays the provided ProphecyClaim to CosmosBridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(provider string, contractAddress common.Address, event types.Event,
	claim ethbridge.ProphecyClaim, key *ecdsa.PrivateKey, sugaredLogger *zap.SugaredLogger) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target, err := initRelayConfig(provider, contractAddress, event, key, sugaredLogger)
	if err != nil {
		sugaredLogger.Errorw("failed in init relay config.",
			errorMessageKey, err.Error())
		return err
	}

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return err
	}

	// Send transaction
	sugaredLogger.Infow("Sending new ProphecyClaim to CosmosBridge.",
		"CosmosSender", claim.CosmosSender,
		"CosmosSenderSequence", claim.CosmosSenderSequence)
	claimType, err := strconv.Atoi(claim.ClaimType)
	if err != nil {
		sugaredLogger.Errorw("failed convert ClaimType", errorMessageKey, err.Error())
		return err
	}
	cosmosSenderSequence, err := strconv.ParseInt(string(claim.CosmosSenderSequence), 0, 64)
	if err != nil {
		sugaredLogger.Errorw("failed convert CosmosSenderSequence", errorMessageKey, err.Error())
		return err
	}
	amount, err := strconv.ParseInt(claim.Amount, 0, 64)
	if err != nil {
		sugaredLogger.Errorw("failed convert Amount", errorMessageKey, err.Error())
		return err
	}
	tx, err := cosmosBridgeInstance.NewProphecyClaim(auth, uint8(claimType),
		claim.CosmosSender, big.NewInt(cosmosSenderSequence), common.BytesToAddress(claim.EthereumReceiver), claim.Symbol, big.NewInt(amount))

	if err != nil {
		sugaredLogger.Errorw("failed to send ProphecyClaim to CosmosBridge.",
			errorMessageKey, err.Error())
		return err
	}

	sugaredLogger.Infow("get NewProphecyClaim tx hash:", "ProphecyClaimHash", tx.Hash().Hex())

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		sugaredLogger.Errorw("failed to get transaction receipt.",
			errorMessageKey, err.Error())
		return err
	}

	switch receipt.Status {
	case 0:
		sugaredLogger.Infow("trasaction failed:")
	case 1:
		sugaredLogger.Infow("trasaction is successful:")
	}
	return nil
}

// initRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func initRelayConfig(provider string, registry common.Address, event types.Event, key *ecdsa.PrivateKey,
	sugaredLogger *zap.SugaredLogger) (*ethclient.Client, *bind.TransactOpts, common.Address, error) {
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

	// rate limit the bridge so that nonce is handled correctly
	time.Sleep(transactionInterval)

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
