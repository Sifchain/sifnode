package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSwapFeeRate(ctx sdk.Context, params *types.SwapFeeRate) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.SwapFeeRatePrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetSwapFeeRate(ctx sdk.Context) *types.SwapFeeRate {
	params := types.SwapFeeRate{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SwapFeeRatePrefix)
	k.cdc.MustUnmarshal(bz, &params)
	return &params
}
