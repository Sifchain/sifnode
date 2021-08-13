package app

import (
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tokenregistrymigrations "github.com/Sifchain/sifnode/x/tokenregistry/migrations"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.1", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.2", func(ctx sdk.Context, plan types.Plan) {
		tokenregistrymigrations.Init(ctx, app.TokenRegistryKeeper)
	})
	app.BaseApp.SetStoreLoader(func(ms sdk.CommitMultiStore) error {
		return ms.LoadLatestVersionAndUpgrade(&upgradetypes.StoreUpgrades{
			Added: []string{tokenregistrytypes.ModuleName},
		})
	})
}
