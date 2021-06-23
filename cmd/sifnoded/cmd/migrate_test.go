package cmd

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/app"
)

func TestMigrateGenesisDataCmd(t *testing.T) {
	cmd, _ := NewRootCmd()
	migrateOutputBuf := new(bytes.Buffer)
	cmd.SetOut(migrateOutputBuf)
	// This test file has been run through sifnoded migrate, and IBC state added.
	cmd.SetArgs([]string{"migrate-data", "v0.9", "testdata/v039_exported_migrated_state.json"})

	app.SetConfig(false)

	homeDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(homeDir)

	err = svrcmd.Execute(cmd, homeDir)
	require.NoError(t, err)

	cmd, _ = NewRootCmd()
	cmd.SetArgs([]string{"init", "test", "--home=" + homeDir})
	err = svrcmd.Execute(cmd, homeDir)
	require.NoError(t, err)

	cmd, _ = NewRootCmd()
	cmd.SetArgs([]string{"unsafe-reset-all"})
	err = svrcmd.Execute(cmd, homeDir)
	require.NoError(t, err)

	err = ioutil.WriteFile(homeDir+"/config/genesis.json", migrateOutputBuf.Bytes(), fs.ModeExclusive)
	require.NoError(t, err)

	cmd, _ = NewRootCmd()
	cmd.SetArgs([]string{"validate-genesis"})
	err = svrcmd.Execute(cmd, homeDir)
	require.NoError(t, err)
}
