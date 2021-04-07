package app

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

func GetAddDispensation(app *SifchainApp) func(ctx sdk.Context, plan upgrade.Plan) {
	return func(ctx sdk.Context, plan upgrade.Plan) {
		dispensation.InitGenesis(ctx, app.DispensationKeeper, dispensation.DefaultGenesisState())
	}
}
