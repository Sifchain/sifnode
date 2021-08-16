package app

import (
	tokenregistrymigrations "github.com/Sifchain/sifnode/x/tokenregistry/migrations"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const upgradeNameV093 = "0.9.3-ibc"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.1", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.2", func(ctx sdk.Context, plan types.Plan) {})
	SetupUpgradeV093(app)
}

func SetupUpgradeV093(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler(upgradeNameV093, func(ctx sdk.Context, plan types.Plan) {
		app.Logger().Info("Running upgrade handler for " + upgradeNameV093 + " with new store " + tokenregistrytypes.StoreKey)
		tokenregistrymigrations.Init(ctx, app.TokenRegistryKeeper)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeNameV093 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{tokenregistrytypes.StoreKey},
		}

		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
