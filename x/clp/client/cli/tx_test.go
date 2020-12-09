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
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--externalAmount", "100",
		"--nativeAmount", "100",
		"--symbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdAddLiquidity(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdAddLiquidity(cdc)
	SetupViper()
	viper.Set(cli.FlagExternalAssetAmount, "100")
	viper.Set(cli.FlagNativeAssetAmount, "100")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--externalAmount", "100",
		"--nativeAmount", "100",
		"--symbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdRemoveLiquidity(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdRemoveLiquidity(cdc)
	SetupViper()
	viper.Set(cli.FlagWBasisPoints, "100")
	viper.Set(cli.FlagAsymmetry, "1000")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--symbol", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err := clpcmd.Execute()
	assert.NoError(t, err)

	viper.Set(cli.FlagWBasisPoints, "%%")
	viper.Set(cli.FlagAsymmetry, "1000")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--symbol", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err = clpcmd.Execute()
	assert.Error(t, err)

	viper.Set(cli.FlagWBasisPoints, "100")
	viper.Set(cli.FlagAsymmetry, "asdef")
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--wBasis", "100",
		"--asymmetry", "1000",
		"--symbol", "eth"})
	clpcmd.SetOut(ioutil.Discard)
	err = clpcmd.Execute()
	assert.Error(t, err)
}

func TestGetCmdSwap(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdSwap(cdc)
	SetupViper()

	viper.Set(cli.FlagSentAssetSymbol, "eth")
	viper.Set(cli.FlagReceivedAssetSymbol, "dash")
	viper.Set(cli.FlagAmount, "100")
	clpcmd.SetArgs([]string{
		"--sentAmount", "100",
		"--receivedSymbol", "dash",
		"--sentSymbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdDecommissionPool(t *testing.T) {
	cdc := test.MakeTestCodec()
	clpcmd := cli.GetCmdDecommissionPool(cdc)
	SetupViper()
	viper.Set(cli.FlagAssetSymbol, "eth")
	clpcmd.SetArgs([]string{
		"--symbol", "eth"})
	err := clpcmd.Execute()
	assert.NoError(t, err)
}
