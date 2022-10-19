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

var bridgeBankAddress = common.HexToAddress(nullAddress)
var cosmosBridgeAddress = common.HexToAddress(nullAddress)

// String returns the event type as a string
func (d ContractRegistry) String() string {
	return [...]string{"bridgebank", "cosmosbridge"}[d-1]
}

// GetAddressFromBridgeRegistry queries the requested contract address from the BridgeRegistry contract
func GetAddressFromBridgeRegistry(client *ethclient.Client, registry common.Address, target ContractRegistry,
	sugaredLogger *zap.SugaredLogger) (common.Address, error) {
	// Return address if already got and stored
	switch target {
	case BridgeBank:
		if bridgeBankAddress.Hex() != nullAddress {
			return bridgeBankAddress, nil
		}
	case CosmosBridge:
		if cosmosBridgeAddress.Hex() != nullAddress {
			return cosmosBridgeAddress, nil
		}
	default:
		panic("invalid target contract address")
	}

	// load sender for query
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

	address, err := registryInstance.BridgeBank(&auth)
	if err != nil {
		sugaredLogger.Errorw("failed to get bridge bank address from registry", errorMessageKey, err.Error())
		return common.Address{}, err
	}
	bridgeBankAddress = address

	address, err = registryInstance.CosmosBridge(&auth)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmos bridge address from registry", errorMessageKey, err.Error())
		return common.Address{}, err
	}
	cosmosBridgeAddress = address

	switch target {
	case BridgeBank:
		return bridgeBankAddress, nil
	case CosmosBridge:
		return cosmosBridgeAddress, nil
	default:
		panic("invalid target contract address")
	}
}
