package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"
)

const upgradeName = "0.10.0-rc"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.10.0-rc.2", func(ctx sdk.Context, plan types.Plan, fromVM m.VersionMap) (m.VersionMap, error) {
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})
	app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + upgradeName)
		app.IBCKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())
		fromVM := make(map[string]uint64)
		// Old Modules can execute Migrations if needed .
		// State migration logic should be added to a migrator function
		// Migrating modules should increment the ConsensusVersion
		// FromVersion NotEqual to ConsensusVersion is required to trigger a migration.
		for moduleName := range app.mm.Modules {
			fromVM[moduleName] = 1
		}
		// New Modules must execute Init Genesis
		fromVM[authz.ModuleName] = 0
		fromVM[feegrant.ModuleName] = 0
		fromVM["vesting"] = 0
		fromVM["crisis"] = 0
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName, "vesting", "crisis"},
		}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
