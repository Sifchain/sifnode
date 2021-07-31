package relayer

import (
	"log"
	"testing"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

const (
	tmProvider        = "Node"
	ethProvider       = "ws://127.0.0.1:7545/"
	contractAddress   = "0x00"
	networkDescriptor = uint32(1)
	privateKeyStr     = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
	validatorMoniker  = "validatorMoniker"
)

func TestNewCosmosSub(t *testing.T) {
	cliContext, _ := client.GetClientTxContext(nil)

	db, err := leveldb.OpenFile("relayerdb", nil)
	require.Equal(t, err, nil)
	privateKey, _ := crypto.HexToECDSA(privateKeyStr)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}

	sugaredLogger := logger.Sugar()
	registryContractAddress := common.HexToAddress(contractAddress)
	sub := NewCosmosSub(oracletypes.NetworkDescriptor(networkDescriptor), privateKey, tmProvider, ethProvider, registryContractAddress,
		db, cliContext, validatorMoniker, false, sugaredLogger)
	require.NotEqual(t, sub, nil)
}

func TestMessageProcessed(t *testing.T) {
	message := txs.CreateTestCosmosMsg(t, types.MsgBurn)
	var claims []types.ProphecyClaimUnique
	claims = append(claims, types.ProphecyClaimUnique{
		ProphecyID: []byte{},
	})

	processed := MessageProcessed(message, claims)
	require.Equal(t, processed, true)
}

func TestMessageNotProcessed(t *testing.T) {
	message := txs.CreateTestCosmosMsg(t, types.MsgBurn)
	var claims []types.ProphecyClaimUnique
	claims = append(claims, types.ProphecyClaimUnique{
		ProphecyID: []byte{},
	})

	processed := MessageProcessed(message, claims)
	require.Equal(t, processed, false)
}

func TestMyDecode(t *testing.T) {
	wrongData := []byte("wrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongDatawrongData")
	_, err := MyDecode(wrongData)
	require.Error(t, err)
}
