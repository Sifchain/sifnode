package app

import (
	"fmt"
	dispensation "github.com/Sifchain/sifnode/x/dispensation/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const releaseVersion = "0.12.0"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)
		//fromVM := make(map[string]uint64)
		//for moduleName := range app.mm.Modules {
		//	fromVM[moduleName] = 1
		//}
		//fromVM[dispensation.ModuleName] = 0
		fromVersion, exists := vm[dispensation.ModuleName]
		fmt.Println("fromversion ", fromVersion, "| exists :", exists)
		return app.mm.RunMigrations(ctx, app.configurator, vm)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == releaseVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
