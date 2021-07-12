package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	// dispensationkeeper "github.com/Sifchain/sifnode/x/dispensation/keeper"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.0-rc.7", func(ctx sdk.Context, plan types.Plan) {})
}
