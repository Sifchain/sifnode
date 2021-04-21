package app

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func GetAddDispensation(app *SifchainApp) func(ctx sdk.Context, plan upgradetypes.Plan) {
	return func(ctx sdk.Context, plan upgradetypes.Plan) {
		dispensation.InitGenesis(ctx, app.DispensationKeeper, dispensation.DefaultGenesisState())
	}
}
