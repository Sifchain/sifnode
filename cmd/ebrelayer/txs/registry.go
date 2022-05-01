package txs

// DONTCOVER

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	bridgeregistry "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/artifacts/contracts/BridgeRegistry.sol"
)

// ContractRegistry is an enum for the bridge contract types
type ContractRegistry byte

const (
	// BridgeBank bridgeBank contract
	BridgeBank ContractRegistry = iota + 1
	// CosmosBridge cosmosBridge contract
	CosmosBridge
)

// String returns the event type as a string
func (d ContractRegistry) String() string {
	return [...]string{"bridgebank", "cosmosbridge"}[d-1]
}

// GetAddressFromBridgeRegistry queries the requested contract address from the BridgeRegistry contract
func GetAddressFromBridgeRegistry(client *ethclient.Client, registry common.Address, target ContractRegistry,
	sugaredLogger *zap.SugaredLogger) (common.Address, error) {
	sender, err := LoadSender()
	if err != nil {
		sugaredLogger.Errorw("failed to get sender", errorMessageKey, err.Error())
		return common.Address{}, err
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		sugaredLogger.Errorw("failed to get header", errorMessageKey, err.Error())

		return common.Address{}, err
	}

	// Set up CallOpts auth
	auth := bind.CallOpts{
		Pending:     true,
		From:        sender,
		BlockNumber: header.Number,
		Context:     context.Background(),
	}

	// Initialize BridgeRegistry instance
	registryInstance, err := bridgeregistry.NewBridgeRegistry(registry, client)
	if err != nil {
		sugaredLogger.Errorw("failed to get registry contract address", errorMessageKey, err.Error())
		return common.Address{}, err
	}

	var address common.Address
	switch target {
	case BridgeBank:
		address, err = registryInstance.BridgeBank(&auth)
	case CosmosBridge:
		address, err = registryInstance.CosmosBridge(&auth)
	default:
		panic("invalid target contract address")
	}

	if err != nil {
		sugaredLogger.Errorw("failed to get contract address from registry", errorMessageKey, err.Error())
		return common.Address{}, err
	}

	return address, nil
}
