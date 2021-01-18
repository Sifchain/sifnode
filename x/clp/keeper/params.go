package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MinCreatePoolThreshold(ctx sdk.Context) uint {
	var minThreshold uint
	k.paramSpace.Get(ctx, types.KeyMinCreatePoolThreshold, &minThreshold)
	return minThreshold
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.MinCreatePoolThreshold(ctx))
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.Set(ctx, types.ParamStoreKey, &params)
}
