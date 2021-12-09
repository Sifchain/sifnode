package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sifchain/sifnode/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/stretchr/testify/require"
)

func TestGenerateEntry(t *testing.T) {
	cmd, _ := NewRootCmd()
	cmd.SetArgs([]string{
		"query",
		"tokenregistry",
		"generate",
		"--token_base_denom", "uatom",
		"--token_ibc_channel_id", "channel-0",
		"--token_decimals", "6",
	})
	app.SetConfig(false)
	homeDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(homeDir)
	err = svrcmd.Execute(cmd, homeDir)
	require.NoError(t, err)
}
