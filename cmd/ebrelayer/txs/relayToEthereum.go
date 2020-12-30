package txs

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit = uint64(3000000)
)

var (
	callMsg *ethereum.CallMsg
)

// RelayProphecyClaimToEthereum relays the provided ProphecyClaim to CosmosBridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(provider string, contractAddress common.Address, event types.Event,
	claim ProphecyClaim, key *ecdsa.PrivateKey, cethAmount *big.Int) (uint64, error) {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target, err := initRelayConfig(provider, contractAddress, event, key)
	if err != nil {
		return 0, err
	}

	// Estimate the cost for the transaction
	if callMsg != nil {
		// Update gas price
		callMsg.GasPrice = auth.GasPrice
	}
	estimateGas, err := client.EstimateGas(context.Background(), *callMsg)

	if err != nil {
		return 0, err
	}

	// If ceth amount is lower than estimated gas
	if cethAmount.Cmp(big.NewInt(int64(estimateGas))) > 0 {
		return 0, errors.New("ceth paid not enough")
	}

	// Initialize CosmosBridge instance
	fmt.Println("\nFetching CosmosBridge contract...")
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// Send transaction
	fmt.Println("Sending new ProphecyClaim to CosmosBridge...")
	tx, err := cosmosBridgeInstance.NewProphecyClaim(auth, uint8(claim.ClaimType),
		claim.CosmosSender, claim.CosmosSenderSequence, claim.EthereumReceiver, claim.Symbol, claim.Amount.BigInt())
	if err != nil {
		log.Println(err)
		return 0, err
	}
	fmt.Println("NewProphecyClaim tx hash:", tx.Hash().Hex())

	// Init call message from transaction data
	if callMsg == nil {
		callMsg = &ethereum.CallMsg{
			From:     auth.From,
			To:       &contractAddress,
			Gas:      0,
			GasPrice: auth.GasPrice,
			Value:    big.NewInt(0),
			Data:     tx.Data(),
		}
	}

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Println(err)
		return 0, err
	}

	switch receipt.Status {
	case 0:
		fmt.Println("Tx Status: 0 - Failed")
		return 0, errors.New("NewProphecyClaim transaction failed ")
	case 1:
		fmt.Println("Tx Status: 1 - Successful")
		return receipt.GasUsed, nil
	}

	return 0, nil
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

	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Println(err)
		return nil, nil, common.Address{}, err

	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return nil, nil, common.Address{}, err

	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(key)
	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

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
