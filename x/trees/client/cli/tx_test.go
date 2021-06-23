package cli_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/faucet/client/cli"
	"github.com/Sifchain/sifnode/x/faucet/test"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func SetupViper() {
	viper.Set(flags.FlagKeyringBackend, flags.DefaultKeyringBackend)
	viper.Set(flags.FlagGenerateOnly, true)
	viper.Set(flags.FlagChainID, "sifchainTest")

}

func TestGetCmdRequestCoins(t *testing.T) {
	cdc := test.MakeTestCodec()
	faucetcmd := cli.GetCmdRequestCoins(cdc)
	SetupViper()

	faucetcmd.SetArgs([]string{
		"1000rowan"})
	err := faucetcmd.Execute()
	assert.NoError(t, err)
}

func TestGetCmdAddCoins(t *testing.T) {
	cdc := test.MakeTestCodec()
	faucetcmd := cli.GetCmdAddCoins(cdc)
	SetupViper()

	faucetcmd.SetArgs([]string{
		"1000rowan"})
	err := faucetcmd.Execute()
	assert.NoError(t, err)
}
