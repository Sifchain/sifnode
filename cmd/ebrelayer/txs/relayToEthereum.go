package txs

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit            = uint64(500000)
	transactionInterval = 60 * time.Second
)

// RelayProphecyClaimToEthereum relays the provided ProphecyClaim to CosmosBridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(provider string, contractAddress common.Address, event types.Event,
	claim ProphecyClaim, key *ecdsa.PrivateKey, sugaredLogger *zap.SugaredLogger) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target, err := initRelayConfig(provider, contractAddress, event, key, sugaredLogger)
	if err != nil {
		sugaredLogger.Errorw("failed in init relay config.",
			"error message", err.Error())
		return err
	}

	// Initialize CosmosBridge instance
	// log.Println("\nFetching CosmosBridge contract...")
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmosBridge instance.",
			"error message", err.Error())
		return err
	}

	// Send transaction
	// log.Println("Sending new ProphecyClaim to CosmosBridge...")
	sugaredLogger.Infow("Sending new ProphecyClaim to CosmosBridge.",
		"CosmosSender", claim.CosmosSender,
		"CosmosSenderSequence", claim.CosmosSenderSequence)

	tx, err := cosmosBridgeInstance.NewProphecyClaim(auth, uint8(claim.ClaimType),
		claim.CosmosSender, claim.CosmosSenderSequence, claim.EthereumReceiver, claim.Symbol, claim.Amount.BigInt())

	if err != nil {
		sugaredLogger.Errorw("failed to send ProphecyClaim to CosmosBridge.",
			"error message", err.Error())
		return err
	}

	sugaredLogger.Infow("get NewProphecyClaim tx hash:", "ProphecyClaimHash", tx.Hash().Hex())

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		sugaredLogger.Errorw("failed to get transaction receipt.",
			"error message", err.Error())
		return err
	}

	switch receipt.Status {
	case 0:
		sugaredLogger.Infow("trasaction failed:")
	case 1:
		// log.Println("Tx Status: 1 - Successful")
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
		// log.Println(err)
		sugaredLogger.Errorw("failed to connect ethereum node.",
			"error message", err.Error())
		return nil, nil, common.Address{}, err
	}

	// Load the validator's address
	sender, err := LoadSender()
	if err != nil {
		sugaredLogger.Errorw("failed to load validator address.",
			"error message", err.Error())
		// log.Println(err)
		return nil, nil, common.Address{}, err
	}

	// rate limit the bridge so that nonce is handled correctly
	time.Sleep(transactionInterval)

	nonce, err := client.PendingNonceAt(context.Background(), sender)
	// log.Println("Current eth operator at pending nonce: ", nonce)
	sugaredLogger.Infow("Current eth operator at pending nonce.", "pending nonce", nonce)

	if err != nil {
		// log.Println("Error broadcasting tx: ", err)
		sugaredLogger.Errorw("failed to broadcast transaction.",
			"error message", err.Error())
		return nil, nil, common.Address{}, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		// log.Println(err)
		sugaredLogger.Errorw("failed to get gas price.",
			"error message", err.Error())
		return nil, nil, common.Address{}, err
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(key)

	// log.Printf("ethereum tx current nonce from client api is %d\n", nonce)

	// log.Println("suggested gas price: ", gasPrice)

	sugaredLogger.Infow("ethereum tx current nonce from client api.",
		"nonce", nonce,
		"suggested gas price", gasPrice)

	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))
	// log.Println("suggested gas price after multiplying by 2: ", gasPrice)

	quarterGasPrice := big.NewInt(0)
	quarterGasPrice = quarterGasPrice.Div(gasPrice, big.NewInt(4))
	// log.Println("quarterGasPrice: ", quarterGasPrice)
	// log.Println("gasPrice after: ", gasPrice)

	gasPrice.Sub(gasPrice, quarterGasPrice)
	// log.Println("suggested gas price after subtracting 1/4: ", gasPrice)
	sugaredLogger.Infow("final gas price after adjustment.",
		"final gas price", gasPrice)

	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	// log.Println("transactOptsAuth.Nonce: ", transactOptsAuth.Nonce)
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
		// log.Println(err)
		sugaredLogger.Errorw("failed to get cosmos bridger contract address from registry.",
			"error message", err.Error())
		return nil, nil, common.Address{}, err

	}
	return client, transactOptsAuth, target, nil
}
