package cli_test

import (
	"github.com/Sifchain/sifnode/x/clp/client/cli"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func SetupViper() {
	viper.Set(flags.FlagKeyringBackend, flags.DefaultKeyringBackend)
	viper.Set(flags.FlagGenerateOnly, true)
	viper.Set(flags.FlagChainID, "sifchainTest")

}

func TestGetCmdCreatePool(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdCreatePool(cdc)
	SetupViper()
	viper.Set(cli.FlagExternalAssetAmount, "100")
	viper.Set(cli.FlagNativeAssetAmount, "100")
	viper.Set(cli.FlagAssetSourceChain, "ethereum")
	viper.Set(cli.FlagAssetSymbol, "ETH")
	viper.Set(cli.FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--externalAmount", "100",
		"--nativeAmount", "100",
		"--sourceChain", "ethereum",
		"--symbol", "ETH",
		"--ticker", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdAddLiquidity(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdAddLiquidity(cdc)
	SetupViper()
	viper.Set(cli.FlagExternalAssetAmount, "100")
	viper.Set(cli.FlagNativeAssetAmount, "100")
	viper.Set(cli.FlagAssetSourceChain, "ethereum")
	viper.Set(cli.FlagAssetSymbol, "ETH")
	viper.Set(cli.FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--externalAmount", "100",
		"--nativeAmount", "100",
		"--sourceChain", "ethereum",
		"--symbol", "ETH",
		"--ticker", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdRemoveLiquidity(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdRemoveLiquidity(cdc)
	SetupViper()
	viper.Set(cli.FlagWBasisPoints, "100")
	viper.Set(cli.FlagAsymmetry, "1000")
	viper.Set(cli.FlagAssetSourceChain, "ethereum")
	viper.Set(cli.FlagAssetSymbol, "ETH")
	viper.Set(cli.FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--sourceChain", "ethereum",
		"--symbol", "ETH",
		"--ticker", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err := clpcmd.Execute()
	assert.NoError(t, err)

	viper.Set(cli.FlagWBasisPoints, "%%")
	viper.Set(cli.FlagAsymmetry, "1000")
	viper.Set(cli.FlagAssetSourceChain, "ethereum")
	viper.Set(cli.FlagAssetSymbol, "ETH")
	viper.Set(cli.FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--sourceChain", "ethereum",
		"--symbol", "ETH",
		"--ticker", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err = clpcmd.Execute()
	assert.Error(t, err)

	viper.Set(cli.FlagWBasisPoints, "100")
	viper.Set(cli.FlagAsymmetry, "asdef")
	viper.Set(cli.FlagAssetSourceChain, "ethereum")
	viper.Set(cli.FlagAssetSymbol, "ETH")
	viper.Set(cli.FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--sourceChain", "ethereum",
		"--symbol", "ETH",
		"--ticker", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err = clpcmd.Execute()
	assert.Error(t, err)
}

func TestGetCmdSwap(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdSwap(cdc)
	SetupViper()

	viper.Set(cli.FlagSentAssetSourceChain, "ethereum")
	viper.Set(cli.FlagSentAssetSymbol, "ETH")
	viper.Set(cli.FlagSentAssetTicker, "ceth")
	viper.Set(cli.FlagReceivedAssetSourceChain, "dash")
	viper.Set(cli.FlagReceivedAssetSymbol, "DASH")
	viper.Set(cli.FlagReceivedAssetTicker, "cdash")
	viper.Set(cli.FlagAmount, "100")
	clpcmd.SetArgs([]string{
		"--sentAmount", "100",
		"--receivedSourceChain", "dash",
		"--receivedSymbol", "DASH",
		"--receivedTicker", "cdash",
		"--sentSourceChain", "ethereum",
		"--sentSymbol", "ETH",
		"--sentTicker", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdDecommissionPool(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdDecommissionPool(cdc)
	SetupViper()
	viper.Set(cli.FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--ticker", "ceth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}
