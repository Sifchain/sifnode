package cmd_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Sifchain/sifnode/app"
	sifnodedcmd "github.com/Sifchain/sifnode/cmd/sifnoded/cmd"
	"github.com/Sifchain/sifnode/x/oracle"
)

func TestAddGenesisValidatorCmd(t *testing.T) {
	homeDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(homeDir)

	initCmd, _ := sifnodedcmd.NewRootCmd()
	initBuf := new(bytes.Buffer)
	initCmd.SetOut(initBuf)
	initCmd.SetErr(initBuf)
	initCmd.SetArgs([]string{"init", "test", "--home=" + homeDir})

	app.SetConfig(false)
	expectedValidatorBech32 := "sifvaloper1rwqp4q88ue83ag3kgnmxxypq0td59df4782tjn"
	expectedValidator, err := sdk.ValAddressFromBech32(expectedValidatorBech32)
	require.NoError(t, err)

	addValCmd, _ := sifnodedcmd.NewRootCmd()
	addValBuf := new(bytes.Buffer)
	addValCmd.SetOut(addValBuf)
	addValCmd.SetErr(addValBuf)
	addValCmd.SetArgs([]string{"add-genesis-validators", expectedValidatorBech32, "--home=" + homeDir})

	// Run init
	err = svrcmd.Execute(initCmd, homeDir)
	require.NoError(t, err)
	// Run add-genesis-validators
	err = svrcmd.Execute(addValCmd, homeDir)
	require.NoError(t, err)
	// Load genesis state from temp home dir and parse JSON
	serverCtx := server.GetServerContextFromCmd(addValCmd)
	genFile := serverCtx.Config.GenesisFile()
	appState, _, err := genutiltypes.GenesisStateFromGenFile(genFile)
	require.NoError(t, err)
	// Setup app to get oracle keeper and ctx.
	sifapp := app.Setup(false)
	ctx := sifapp.BaseApp.NewContext(false, tmproto.Header{})
	// Run loaded genesis through InitGenesis on oracle module
	mm := module.NewManager(
		oracle.NewAppModule(sifapp.OracleKeeper),
	)
	_ = mm.InitGenesis(ctx, sifapp.AppCodec(), appState)
	// Assert validator
	validators := sifapp.OracleKeeper.GetOracleWhiteList(ctx)
	require.Equal(t, []sdk.ValAddress{expectedValidator}, validators)
}
