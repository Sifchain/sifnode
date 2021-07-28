package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tokenregistrymigrations "github.com/Sifchain/sifnode/x/tokenregistry/migrations"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0-rc.7", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0-rc.8", func(ctx sdk.Context, plan types.Plan) {
		tokenregistrymigrations.Init(ctx, app.TokenRegistryKeeper)
	})
}
