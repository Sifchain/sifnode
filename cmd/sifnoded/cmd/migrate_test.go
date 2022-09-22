package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/app"
)

func TestMigrateGenesisDataCmd(t *testing.T) {
	homeDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(homeDir)
	cmd, _ := NewRootCmd()
	migrateOutputBuf := new(bytes.Buffer)
	cmd.SetOut(migrateOutputBuf)
	cmd.SetArgs([]string{"migrate-data", "v0.9", "testdata/v039_exported_migrated_state.json"})
	app.SetConfig(false)
	err = svrcmd.Execute(cmd, homeDir)
	require.NoError(t, err)
	cmd, _ = NewRootCmd()
	cmd.SetArgs([]string{"init", "test", "--home=" + homeDir})
	err = svrcmd.Execute(cmd, homeDir)
	require.NoError(t, err)
	err = ioutil.WriteFile(homeDir+"/config/genesis.json", migrateOutputBuf.Bytes(), 0600)
	require.NoError(t, err)
}
