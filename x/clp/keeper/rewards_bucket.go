package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRewardsBucket set a specific rewardsBucket in the store from its index
func (k Keeper) SetRewardsBucket(ctx sdk.Context, rewardsBucket types.RewardsBucket) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardsBucketKeyPrefix))
	b := k.cdc.MustMarshal(&rewardsBucket)
	store.Set(types.RewardsBucketKey(
		rewardsBucket.Denom,
	), b)
}

// GetRewardsBucket returns a rewardsBucket from its index
func (k Keeper) GetRewardsBucket(
	ctx sdk.Context,
	denom string,

) (val types.RewardsBucket, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardsBucketKeyPrefix))

	b := store.Get(types.RewardsBucketKey(
		denom,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveRewardsBucket removes a rewardsBucket from the store
func (k Keeper) RemoveRewardsBucket(
	ctx sdk.Context,
	denom string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardsBucketKeyPrefix))
	store.Delete(types.RewardsBucketKey(
		denom,
	))
}

// GetAllRewardsBucket returns all rewardsBucket
func (k Keeper) GetAllRewardsBucket(ctx sdk.Context) (list []types.RewardsBucket) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RewardsBucketKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RewardsBucket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
