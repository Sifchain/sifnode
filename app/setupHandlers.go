package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("changePoolFormula", GetPoolChangeFunc(app))
	app.UpgradeKeeper.SetUpgradeHandler("release-20210401000000", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("release-20210407042000", GetAddDispensation(app))
	app.UpgradeKeeper.SetUpgradeHandler("release-20210414000000", func(ctx sdk.Context, plan types.Plan) {})
}
