package cli

import (
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
	clpcmd := GetCmdCreatePool(cdc)
	SetupViper()
	viper.Set(FlagExternalAssetAmount, "100")
	viper.Set(FlagNativeAssetAmount, "100")
	viper.Set(FlagAssetSourceChain, "ethereum")
	viper.Set(FlagAssetSymbol, "ETH")
	viper.Set(FlagAssetTicker, "ceth")
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
	clpcmd := GetCmdAddLiquidity(cdc)
	SetupViper()
	viper.Set(FlagExternalAssetAmount, "100")
	viper.Set(FlagNativeAssetAmount, "100")
	viper.Set(FlagAssetSourceChain, "ethereum")
	viper.Set(FlagAssetSymbol, "ETH")
	viper.Set(FlagAssetTicker, "ceth")
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
	clpcmd := GetCmdRemoveLiquidity(cdc)
	SetupViper()
	viper.Set(FlagWBasisPoints, "100")
	viper.Set(FlagAsymmetry, "1000")
	viper.Set(FlagAssetSourceChain, "ethereum")
	viper.Set(FlagAssetSymbol, "ETH")
	viper.Set(FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--sourceChain", "ethereum",
		"--symbol", "ETH",
		"--ticker", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err := clpcmd.Execute()
	assert.NoError(t, err)

	viper.Set(FlagWBasisPoints, "%%")
	viper.Set(FlagAsymmetry, "1000")
	viper.Set(FlagAssetSourceChain, "ethereum")
	viper.Set(FlagAssetSymbol, "ETH")
	viper.Set(FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--sourceChain", "ethereum",
		"--symbol", "ETH",
		"--ticker", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err = clpcmd.Execute()
	assert.Error(t, err)

	viper.Set(FlagWBasisPoints, "100")
	viper.Set(FlagAsymmetry, "asdef")
	viper.Set(FlagAssetSourceChain, "ethereum")
	viper.Set(FlagAssetSymbol, "ETH")
	viper.Set(FlagAssetTicker, "ceth")
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
	clpcmd := GetCmdSwap(cdc)
	SetupViper()

	viper.Set(FlagSentAssetSourceChain, "ethereum")
	viper.Set(FlagSentAssetSymbol, "ETH")
	viper.Set(FlagSentAssetTicker, "ceth")
	viper.Set(FlagReceivedAssetSourceChain, "dash")
	viper.Set(FlagReceivedAssetSymbol, "DASH")
	viper.Set(FlagReceivedAssetTicker, "cdash")
	viper.Set(FlagAmount, "100")
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
	clpcmd := GetCmdDecommissionPool(cdc)
	SetupViper()
	viper.Set(FlagAssetTicker, "ceth")
	clpcmd.SetArgs([]string{
		"--ticker", "ceth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}
