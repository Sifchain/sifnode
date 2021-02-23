package txs

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit            = uint64(200000)
	transactionInterval = 60 * time.Second
)

// RelayProphecyClaimToEthereum relays the provided ProphecyClaim to CosmosBridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(provider string, contractAddress common.Address, event types.Event,
	claim ProphecyClaim, key *ecdsa.PrivateKey) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target, err := initRelayConfig(provider, contractAddress, event, key)
	if err != nil {
		return err
	}

	// Initialize CosmosBridge instance
	fmt.Println("\nFetching CosmosBridge contract...")
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		log.Println(err)
		return err
	}

	// Send transaction
	fmt.Println("Sending new ProphecyClaim to CosmosBridge...")
	tx, err := cosmosBridgeInstance.NewProphecyClaim(auth, uint8(claim.ClaimType),
		claim.CosmosSender, claim.CosmosSenderSequence, claim.EthereumReceiver, claim.Symbol, claim.Amount.BigInt())

	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("NewProphecyClaim tx hash:", tx.Hash().Hex())

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Println(err)
		return err
	}

	switch receipt.Status {
	case 0:
		fmt.Println("Tx Status: 0 - Failed")
	case 1:
		fmt.Println("Tx Status: 1 - Successful")
	}
	return nil
}

// initRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func initRelayConfig(provider string, registry common.Address, event types.Event, key *ecdsa.PrivateKey,
) (*ethclient.Client, *bind.TransactOpts, common.Address, error) {
	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Println(err)
		return nil, nil, common.Address{}, err
	}

	// Load the validator's address
	sender, err := LoadSender()
	if err != nil {
		log.Println(err)
		return nil, nil, common.Address{}, err
	}

	// rate limit the bridge so that nonce is handled correctly
	time.Sleep(transactionInterval)

	nonce, err := client.PendingNonceAt(context.Background(), sender)
	log.Println("Current eth operator at pending nonce: ", nonce)
	if err != nil {
		log.Println("Error broadcasting tx: ", err)
		return nil, nil, common.Address{}, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return nil, nil, common.Address{}, err
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(key)

	log.Printf("ethereum tx current nonce from client api is %d\n", nonce)

	fmt.Println("suggested gas price: ", gasPrice)

	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))
	fmt.Println("suggested gas price after multiplying by 2: ", gasPrice)

	quarterGasPrice := big.NewInt(0)
	quarterGasPrice = quarterGasPrice.Div(gasPrice, big.NewInt(4))
	fmt.Println("quarterGasPrice: ", quarterGasPrice)
	fmt.Println("gasPrice after: ", gasPrice)

	gasPrice.Sub(gasPrice, quarterGasPrice)
	fmt.Println("suggested gas price after subtracting 1/4: ", gasPrice)

	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	log.Println("transactOptsAuth.Nonce: ", transactOptsAuth.Nonce)

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
	target, err := GetAddressFromBridgeRegistry(client, registry, targetContract)
	if err != nil {
		log.Println(err)
		return nil, nil, common.Address{}, err

	}
	return client, transactOptsAuth, target, nil
}
