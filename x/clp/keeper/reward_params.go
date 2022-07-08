package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetRewardParams(ctx sdk.Context, params *types.RewardParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RewardParamPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetRewardsParams(ctx sdk.Context) *types.RewardParams {
	params := types.RewardParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardParamPrefix)
	k.cdc.MustUnmarshal(bz, &params)
	return &params
}
