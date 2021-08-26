package app

import (
	tokenregistrymigrations "github.com/Sifchain/sifnode/x/tokenregistry/migrations"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const upgradeNameV095 = "0.9.5"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.1", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.2", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.2-ibc.7", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.3", func(ctx sdk.Context, plan types.Plan) {})
	SetupHandlersForV095(app)
}

func SetupHandlersForV095(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.5", func(ctx sdk.Context, plan types.Plan) {
		app.Logger().Info("Running upgrade handler for " + upgradeNameV095 + " with new store " + tokenregistrytypes.StoreKey)
		// Install initial token registry entries for non-ibc tokens.
		tokenregistrymigrations.Init(ctx, app.TokenRegistryKeeper)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeNameV095 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
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
