package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetMinCreatePoolThreshold(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMinCreatePoolThreshold, &res)

	return res
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)

	return params
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
