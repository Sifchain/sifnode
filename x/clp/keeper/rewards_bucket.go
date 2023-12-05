package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k Keeper) ShouldDistributeRewardsToLPWallet(ctx sdk.Context) bool {
	params := k.GetRewardsParams(ctx)
	return params.RewardsDistribute
}

// DistributeLiquidityProviderRewards distributes rewards to a liquidity provider
func (k Keeper) DistributeLiquidityProviderRewards(ctx sdk.Context, lp *types.LiquidityProvider, asset string, rewardAmount sdk.Int) error {
	// get the liquidity provider address
	lpAddress, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	if err != nil {
		return err
	}

	// distribute rewards to the liquidity provider
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lpAddress, sdk.NewCoins(sdk.NewCoin(asset, rewardAmount)))
	if err != nil {
		return err
	}

	// subtract the reward amount from the rewards bucket
	err = k.SubtractFromRewardsBucket(ctx, asset, rewardAmount)
	if err != nil {
		return err
	}

	return nil
}

// calculate the reward share for each liquidity provider
func (k Keeper) CalculateRewardShareForLiquidityProviders(
	ctx sdk.Context,
	lps []*types.LiquidityProvider,
) []sdk.Dec {
	// sum up the liquidity provider total units
	totalUnits := sdk.ZeroInt()
	for _, lp := range lps {
		totalUnits = totalUnits.Add(sdk.NewIntFromBigInt(lp.LiquidityProviderUnits.BigInt()))
	}

	// create a list of the reward share based on lp units and totalUnits
	rewardShares := make([]sdk.Dec, len(lps))
	for i, lp := range lps {
		rewardShares[i] = sdk.NewDecFromInt(sdk.NewIntFromBigInt(lp.LiquidityProviderUnits.BigInt())).Quo(sdk.NewDecFromInt(totalUnits))
	}

	return rewardShares
}

// CalculateRewardAmountForLiquidityProviders calculates the reward amount for each liquidity provider
func (k Keeper) CalculateRewardAmountForLiquidityProviders(
	ctx sdk.Context,
	rewardShares []sdk.Dec,
	rewardsBucketAmount sdk.Int,
) []sdk.Int {
	rewardAmounts := make([]sdk.Int, len(rewardShares))
	for i, rewardShare := range rewardShares {
		rewardAmounts[i] = rewardShare.MulInt(rewardsBucketAmount).TruncateInt()
	}
	return rewardAmounts
}

// AddRewardAmountToLiquidityPool adds a new reward amount to a liquidity pool
func (k Keeper) AddRewardAmountToLiquidityPool(ctx sdk.Context, liquidityProvider *types.LiquidityProvider, asset types.Asset, rewardAmount sdk.Int) error {
	if liquidityProvider.Asset.Equals(asset) == false {
		return types.ErrInValidAsset
	}

	pool, err := k.GetPool(ctx, asset.Symbol)
	if err != nil {
		return types.ErrPoolDoesNotExist
	}

	nativeAssetDepth, externalAssetDepth := pool.ExtractDebt(pool.NativeAssetBalance, pool.ExternalAssetBalance, false)

	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	sellNativeSwapFeeRate := k.GetSwapFeeRate(ctx, types.GetSettlementAsset(), false)
	buyNativeSwapFeeRate := k.GetSwapFeeRate(ctx, asset, false)

	newPoolUnits, lpUnits, _, _, err := CalculatePoolUnits(
		pool.PoolUnits,
		nativeAssetDepth,
		externalAssetDepth,
		sdk.ZeroUint(),
		sdk.Uint(rewardAmount),
		sellNativeSwapFeeRate,
		buyNativeSwapFeeRate,
		pmtpCurrentRunningRate)
	if err != nil {
		return err
	}

	// Update pool total share units
	pool.PoolUnits = newPoolUnits

	// Add to external asset balance
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(sdk.Uint(rewardAmount))

	// Subtract from rewards bucket
	err = k.SubtractFromRewardsBucket(ctx, asset.Symbol, rewardAmount)
	if err != nil {
		return err
	}

	// Update LP units
	liquidityProvider.LiquidityProviderUnits = liquidityProvider.LiquidityProviderUnits.Add(lpUnits)

	// Save new pool balances
	err = k.SetPool(ctx, &pool)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}

	// Save LP
	k.SetLiquidityProvider(ctx, liquidityProvider)

	return nil
}
