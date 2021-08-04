package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.1", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.2", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.3", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.4", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.5", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.6", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.7", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.8", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("0.9.9", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("1.0.0", func(ctx sdk.Context, plan types.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("1.0.1", func(ctx sdk.Context, plan types.Plan) {})
}
