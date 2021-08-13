package app

import (
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tokenregistrymigrations "github.com/Sifchain/sifnode/x/tokenregistry/migrations"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.1", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.2", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.3", func(ctx sdk.Context, plan types.Plan) {
		tokenregistrymigrations.Init(ctx, app.TokenRegistryKeeper)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == "0.9.3" && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{tokenregistrytypes.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
