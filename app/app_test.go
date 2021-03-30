package app

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"
)

// func TestExport(t *testing.T) {
// 	db := db.NewMemDB()
// 	app := NewSifApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
// 	err := setGenesis(app)
// 	assert.NoError(t, err)
// 	_, _, err = app.ExportAppStateAndValidators(false, []string{})
// 	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
// }

// // ensure that black listed addresses are properly set in bank keeper
// func TestBlackListedAddrs(t *testing.T) {
// 	db := db.NewMemDB()
// 	app := NewSifApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

// 	for acc := range maccPerms {
// 		require.True(t, app.BankKeeper.BlacklistedAddr(app.SupplyKeeper.GetModuleAddress(acc)))
// 	}
// }

func setGenesis(app *SifchainApp) error {
	encCfg := MakeTestEncodingConfig()

	genesisState := NewDefaultGenesisState(encCfg.Marshaler)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
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
