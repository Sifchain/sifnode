package relayer

import (
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	tmProvider      = "Node"
	ethProvider     = "ws://127.0.0.1:7545/"
	contractAddress = "0x00"
	privateKeyStr   = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
	// privateKey              = txs.LoadPrivateKey()
	// logger                  = tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))
)

func TestNewCosmosSub(t *testing.T) {

	privateKey, _ := crypto.HexToECDSA(privateKeyStr)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	registryContractAddress := common.HexToAddress(contractAddress)
	a := relayer.NewCosmosSub(tmProvider, ethProvider, registryContractAddress,
		privateKey, logger)
	fmt.Println(a)
}
