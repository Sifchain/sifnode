package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tokenregistrymigrations "github.com/Sifchain/sifnode/x/tokenregistry/migrations"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0", func(ctx sdk.Context, plan types.Plan) {})
}
