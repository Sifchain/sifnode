package relayer

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

const (
	tmProvider      = "Node"
	ethProvider     = "ws://127.0.0.1:7545/"
	contractAddress = "0x00"
)

func TestNewCosmosSub(t *testing.T) {
	db, err := leveldb.OpenFile("relayerdb", nil)
	require.Equal(t, err, nil)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}

	sugaredLogger := logger.Sugar()
	registryContractAddress := common.HexToAddress(contractAddress)
	var key *ecdsa.PrivateKey // this isn't actually used
	sub := NewCosmosSub(tmProvider, ethProvider, registryContractAddress,
		key, db, sugaredLogger)
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
