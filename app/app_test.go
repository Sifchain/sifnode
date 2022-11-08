package app

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func TestAppUpgrade_CannotDeleteLatestVersion(t *testing.T) {
	// Skip this test that shows how SDK allows stores to get into irrecoverable state and panics
	// TODO: Open PR on SDK with this test.
	t.Skip()
	encCfg := MakeTestEncodingConfig()
	db := dbm.NewMemDB()
	app := NewSifApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		false,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		encCfg,
		EmptyAppOptions{},
	)
	err := app.LoadLatestVersion()
	require.NoError(t, err)

	genesisState := NewDefaultGenesisState(encCfg.Marshaler)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	// Commit enough to trigger beginning of pruning soon
	for i := 1; i < 105; i++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: types.Header{Height: int64(i)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(i)})
		app.Commit()
	}

	app = NewSifApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		false,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		encCfg,
		EmptyAppOptions{},
		func(app *baseapp.BaseApp) {
			cms := rootmulti.NewStore(db, app.Logger())
			cms.SetPruning(storetypes.PruneDefault)
			app.SetCMS(cms)
		},
	)
	// Mount and load newStore which will be loaded with version 0,
	// because it did not load with an "add" upgrade to set it's initialVersion,
	// while the other stores will be at version = blockHeight (i.e 1 here).
	app.MountKVStores(map[string]*sdk.KVStoreKey{
		"newStore": sdk.NewKVStoreKey("newStore"),
	})
	err = app.LoadLatestVersion()
	require.NoError(t, err)

	// Commit until just before default pruning is triggered
	for i := 105; i <= 109; i++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: types.Header{Height: int64(i)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(i)})
		app.Commit()
	}

	// Next commit will panic when trying to prune a height on the new store,
	// that is >= the current version of the new store.
	defer func() {
		err := recover()
		require.EqualError(t, err.(error), errors.Errorf("cannot delete latest saved version (%d)", 6).Error())
	}()

	app.BeginBlock(abci.RequestBeginBlock{Header: types.Header{Height: int64(110)}})
	app.EndBlock(abci.RequestEndBlock{Height: int64(110)})
	app.Commit()
}

func TestAppUpgrade_CannotLoadCorruptStoreUsingLatestHeight(t *testing.T) {
	// Skip this test that shows how SDK allows stores to get installed with 0 height,
	// then, after pruning halts chain, it becomes impossible to load them at correct height,
	// loosing any data that was added since last upgrade.
	// Trying to rename (i.e copy) also fails as old store still errors.
	t.Skip()
	encCfg := MakeTestEncodingConfig()
	db := dbm.NewMemDB()
	app := NewSifApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		false,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		encCfg,
		EmptyAppOptions{},
	)
	err := app.LoadLatestVersion()
	require.NoError(t, err)

	genesisState := NewDefaultGenesisState(encCfg.Marshaler)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	// Commit enough to trigger beginning of pruning soon
	for i := 1; i < 105; i++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: types.Header{Height: int64(i)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(i)})
		app.Commit()
	}

	app = NewSifApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		false,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		encCfg,
		EmptyAppOptions{},
		func(app *baseapp.BaseApp) {
			cms := rootmulti.NewStore(db, app.Logger())
			cms.SetPruning(storetypes.PruneDefault)
			app.SetCMS(cms)
		},
	)
	// Mount and load newStore which will be loaded with version 0,
	// because it did not load with an "add" upgrade to set it's initialVersion,
	// while the other stores will be at version = blockHeight (i.e 1 here).
	app.MountKVStores(map[string]*sdk.KVStoreKey{
		"newStore": sdk.NewKVStoreKey("newStore"),
	})
	err = app.LoadLatestVersion()
	require.NoError(t, err)

	// Commit until just before default pruning is triggered
	for i := 105; i <= 109; i++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: types.Header{Height: int64(i)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(i)})
		app.Commit()
	}

	// Try to recover the store by setting an "add" store upgrade,
	// after it has already been committed to above, before setting initial height.
	app = NewSifApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		db,
		nil,
		false,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		encCfg,
		EmptyAppOptions{},
		func(app *baseapp.BaseApp) {
			cms := rootmulti.NewStore(db, app.Logger())
			cms.SetPruning(storetypes.PruneDefault)
			app.SetCMS(cms)
		},
	)
	// Mount and load newStore which will be loaded with version 0,
	// because it did not load with an "add" upgrade to set it's initialVersion,
	// while the other stores will be at version = blockHeight (i.e 1 here).
	app.MountKVStores(map[string]*sdk.KVStoreKey{
		"newStore": sdk.NewKVStoreKey("newStore"),
	})
	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(110, &storetypes.StoreUpgrades{
		Added: []string{"newStore"},
	}))

	// Trigger the upgrade store loader set above, instead of the default store loader.
	// Load will error when trying to load new store at current block height,
	// because it is > the latest store version (i.e there is an earlier version of the store).
	err = app.LoadLatestVersion()
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.Errorf("initial version set to %v, but found earlier version %v", 110, 1).Error())
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

func TestSimAppExportAndBlockedAddrs(t *testing.T) {
	encCfg := MakeTestEncodingConfig()
	db := dbm.NewMemDB()
	app := NewSifApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encCfg, EmptyAppOptions{})
	SetConfig(false)
	for acc := range maccPerms {
		require.True(
			t,
			app.BankKeeper.BlockedAddr(app.AccountKeeper.GetModuleAddress(acc)),
			"ensure that blocked addresses are properly set in bank keeper",
		)
	}

	genesisState := NewDefaultGenesisState(encCfg.Marshaler)
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewSifApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, encCfg, EmptyAppOptions{})
	_, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestAddressFormatValidation(t *testing.T) {
	addr0 := make([]byte, 0)
	err := sdk.VerifyAddressFormat(addr0)
	assert.Error(t, err, "addresses cannot be empty: unknown address")
	addr5 := make([]byte, 5)
	err = sdk.VerifyAddressFormat(addr5)
	assert.NoError(t, err)
	addr20 := make([]byte, 20)
	err = sdk.VerifyAddressFormat(addr20)
	assert.NoError(t, err)
	addr32 := make([]byte, 32)
	err = sdk.VerifyAddressFormat(addr32)
	assert.NoError(t, err)
	addr256 := make([]byte, 256)
	err = sdk.VerifyAddressFormat(addr256)
	assert.Error(t, err, "address max length is 255, got 256: unknown address")
}
