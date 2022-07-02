package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetLiquidityProtectionParams(ctx sdk.Context, params *types.LiquidityProtectionParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LiquidityProtectionParamsPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetLiquidityProtectionParams(ctx sdk.Context) *types.LiquidityProtectionParams {
	params := types.LiquidityProtectionParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LiquidityProtectionParamsPrefix)
	k.cdc.MustUnmarshal(bz, &params)
	return &params
}
