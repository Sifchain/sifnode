package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
// TODO: Define if your module needs Parameters, if not this can be deleted

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/utx0/sifnode/x/clp/types"
)

// GetParams returns the total set of clp parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the clp parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}
*/
func (k Keeper) MinCreatePoolThreshold(ctx sdk.Context) (res uint) {
	k.paramstore.Get(ctx, types.KeyMinCreatePoolThreshold, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.MinCreatePoolThreshold(ctx))
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
