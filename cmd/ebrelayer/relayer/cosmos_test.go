package relayer

import (
	"bufio"
	"os"
	"testing"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	tmProvider       = "Node"
	ethProvider      = "ws://127.0.0.1:7545/"
	contractAddress  = "0x00"
	privateKeyStr    = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
	rpcURL           = "tcp://localhost:26657"
	chainID          = "sifchain"
	validatorMoniker = "shadowfiend"
	mnemonic         = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
)

func TestNewCosmosSub(t *testing.T) {

	var cmd cobra.Command
	inBuf := bufio.NewReader(cmd.InOrStdin())
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	var cdc codec.Codec
	cosmosContext, err := types.NewCosmosContext(inBuf, &cdc, rpcURL, validatorMoniker, mnemonic, chainID, logger)
	// TODO need simulate the response from sifnoded
	require.NotEqual(t, err, nil)

	privateKey, _ := crypto.HexToECDSA(privateKeyStr)
	registryContractAddress := common.HexToAddress(contractAddress)
	sub := NewCosmosSub(tmProvider, ethProvider, registryContractAddress,
		privateKey, cosmosContext, logger)
	require.NotEqual(t, sub, nil)
}
