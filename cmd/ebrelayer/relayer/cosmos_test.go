package relayer

import (
	"log"
	"math/big"
	"testing"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	tmProvider      = "Node"
	ethProvider     = "ws://127.0.0.1:7545/"
	contractAddress = "0x00"
	privateKeyStr   = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
)

func TestNewCosmosSub(t *testing.T) {

	privateKey, _ := crypto.HexToECDSA(privateKeyStr)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}

	sugaredLogger := logger.Sugar()
	registryContractAddress := common.HexToAddress(contractAddress)
	sub := NewCosmosSub(tmProvider, ethProvider, registryContractAddress,
		privateKey, sugaredLogger)
	require.NotEqual(t, sub, nil)
}

func TestMessageProcessed(t *testing.T) {
	message := txs.CreateTestCosmosMsg(t, types.MsgBurn)
	var claims []types.ProphecyClaimUnique
	claims = append(claims, types.ProphecyClaimUnique{
		CosmosSender:         []byte(txs.TestCosmosAddress1),
		CosmosSenderSequence: big.NewInt(txs.TestCosmosAddressSequence),
	})

	processed := MessageProcessed(message, claims)
	require.Equal(t, processed, true)
}

func TestMessageNotProcessed(t *testing.T) {
	message := txs.CreateTestCosmosMsg(t, types.MsgBurn)
	var claims []types.ProphecyClaimUnique
	claims = append(claims, types.ProphecyClaimUnique{
		CosmosSender:         []byte(txs.TestCosmosAddress1),
		CosmosSenderSequence: big.NewInt(txs.TestCosmosAddressSequence + 1),
	})

	processed := MessageProcessed(message, claims)
	require.Equal(t, processed, false)
}

func TestMyDecode(t *testing.T) {
	wrongData := []byte("wrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongData")
	_, err := MyDecode(wrongData)
	require.Error(t, err)
}
