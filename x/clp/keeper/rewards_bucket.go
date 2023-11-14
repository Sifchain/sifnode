package keeper

import (
	"fmt"

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

// AddToRewardsBucket adds an amount to a specific RewardsBucket in the store,
// or creates a new RewardsBucket if one does not already exist for the denom.
func (k Keeper) AddToRewardsBucket(ctx sdk.Context, denom string, amount sdk.Int) error {
	if denom == "" {
		return types.ErrDenomCantBeEmpty
	}
	if amount.IsNegative() {
		return types.ErrAmountCantBeNegative
	}

	rewardsBucket, found := k.GetRewardsBucket(ctx, denom)
	if !found {
		// Initialize a new RewardsBucket if it does not exist
		rewardsBucket = types.RewardsBucket{
			Denom:  denom,
			Amount: sdk.NewInt(0),
		}
	}

	// Add the amount to the current or new rewards
	newAmount := rewardsBucket.Amount.Add(amount)
	rewardsBucket.Amount = newAmount

	k.SetRewardsBucket(ctx, rewardsBucket)

	return nil
}

// SubtractFromRewardsBucket subtracts an amount from a specific RewardsBucket in the store
func (k Keeper) SubtractFromRewardsBucket(ctx sdk.Context, denom string, amount sdk.Int) error {
	if denom == "" {
		return types.ErrDenomCantBeEmpty
	}
	if amount.IsNegative() {
		return types.ErrAmountCantBeNegative
	}

	rewardsBucket, found := k.GetRewardsBucket(ctx, denom)
	if !found {
		return fmt.Errorf(types.ErrRewardsBucketNotFound.Error(), denom)
	}

	// Check if the rewards bucket has enough to subtract
	if rewardsBucket.Amount.LT(amount) {
		return fmt.Errorf(types.ErrNotEnoughBalanceInRewardsBucket.Error(), denom)
	}

	// Subtract the amount from the current rewards
	newAmount := rewardsBucket.Amount.Sub(amount)
	rewardsBucket.Amount = newAmount

	k.SetRewardsBucket(ctx, rewardsBucket)
	return nil
}

// AddMultipleCoinsToRewardsBuckets adds multiple coin amounts to their respective RewardsBuckets
func (k Keeper) AddMultipleCoinsToRewardsBuckets(ctx sdk.Context, coins sdk.Coins) (sdk.Coins, error) {
	for _, coin := range coins {
		err := k.AddToRewardsBucket(ctx, coin.Denom, coin.Amount)
		if err != nil {
			return nil, err
		}
	}

	// return a list of all the coins added to rewards buckets
	return coins, nil
}

func (k Keeper) ShouldDistributeRewards(ctx sdk.Context, epochIdentifier string) bool {
	params := k.GetRewardsParams(ctx)
	return epochIdentifier == params.RewardsEpochIdentifier
}

// DistributeLiquidityProviderRewards distributes rewards to a liquidity provider
func (k Keeper) DistributeLiquidityProviderRewards(ctx sdk.Context, lp types.LiquidityProvider) error {
	return nil
}
