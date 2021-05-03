package app

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestExport(t *testing.T) {
	db := db.NewMemDB()
	app := NewInitApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	err := setGenesis(app)
	assert.NoError(t, err)
	_, _, err = app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := db.NewMemDB()
	app := NewInitApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	for acc := range maccPerms {
		require.True(t, app.BankKeeper.BlacklistedAddr(app.SupplyKeeper.GetModuleAddress(acc)))
	}
}

func setGenesis(app *SifchainApp) error {
	genesisState := NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(app.Cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()
	return nil
}
