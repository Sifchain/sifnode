package relayer

import (
	"log"
	"math/big"
	"testing"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

const (
	tmProvider        = "Node"
	ethProvider       = "ws://127.0.0.1:7545/"
	contractAddress   = "0x00"
	networkDescriptor = uint32(1)
	privateKeyStr     = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
	validatorMoniker  = "validatorMoniker"
	sifnodeGrpc       = "0.0.0.0:9090"
)

func TestNewCosmosSub(t *testing.T) {

	privateKey, _ := crypto.HexToECDSA(privateKeyStr)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}

	sugaredLogger := logger.Sugar()
	registryContractAddress := common.HexToAddress(contractAddress)

	ethereumChainId := big.NewInt(1)
	maxFeePerGas := big.NewInt(3000)
	maxPriorityFeePerGas := big.NewInt(3000)

	sub := NewCosmosSub(oracletypes.NetworkDescriptor(networkDescriptor),
		privateKey,
		tmProvider,
		ethProvider,
		registryContractAddress,
		client.Context{},
		validatorMoniker,
		sugaredLogger,
		maxFeePerGas,
		maxPriorityFeePerGas,
		ethereumChainId,
		sifnodeGrpc)

	require.NotEqual(t, sub, nil)
}
