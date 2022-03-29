package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) GetMinCreatePoolThreshold(ctx sdk.Context) (res uint64) {
	keeper.paramstore.Get(ctx, types.KeyMinCreatePoolThreshold, &res)
	return res
}

func (keeper Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	keeper.paramstore.GetParamSet(ctx, &params)
	return params
}

// set the params
func (keeper Keeper) SetParams(ctx sdk.Context, params types.Params) {
	keeper.paramstore.SetParamSet(ctx, &params)
}
