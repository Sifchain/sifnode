package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCurrentRewardPeriod(ctx sdk.Context, params *types.RewardParams) *types.RewardPeriod {
	height := uint64(ctx.BlockHeight())
	for _, period := range params.RewardPeriods {
		if height >= period.RewardPeriodStartBlock && height <= period.RewardPeriodEndBlock {
			// mod 0 is undefined - in which case we'll run every block
			if period.RewardPeriodMod == 0 {
				period.RewardPeriodMod = 1
			}
			return period
		}
	}
	return nil
}

func CalcBlockDistribution(period *types.RewardPeriod) sdk.Uint {
	periodLength := period.RewardPeriodEndBlock - period.RewardPeriodStartBlock + 1
	return period.RewardPeriodAllocation.QuoUint64(periodLength)
}

func (k Keeper) DistributeDepthRewards(ctx sdk.Context, blockDistribution sdk.Uint, period *types.RewardPeriod, pools []*types.Pool) error {
	height := uint64(ctx.BlockHeight())
	if height == period.RewardPeriodStartBlock {
		rewardsParams := k.GetRewardsParams(ctx)
		rewardsParams.RewardPeriodStartTime = ctx.BlockTime().String()
		k.SetRewardParams(ctx, rewardsParams)
	}

	if blockDistribution.IsZero() {
		return nil
	}

	totalDepth, err := k.calcTotalDepth(ctx, pools, period, height)
	if err != nil {
		return err
	}

	if totalDepth.LTE(sdk.ZeroDec()) {
		return nil
	}

	tuples, coinsToMint := CollectPoolRewardTuples(pools, blockDistribution, totalDepth, period)
	moduleBalancePreMinting := k.GetModuleRowan(ctx)
	rewardCoins := sdk.NewCoins(sdk.NewCoin(types.NativeSymbol, sdk.NewIntFromBigInt(coinsToMint.BigInt())))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, rewardCoins)
	if err != nil {
		return err
	}

	shouldDistribute := period.RewardPeriodDistribute
	poolRowanMap := make(PoolRowanMap)
	lpRowanMap := make(LpRowanMap)
	lpPoolMap := make(LpPoolMap)

	var partitions map[types.Asset][]*types.LiquidityProvider
	if shouldDistribute {
		partitions, err = k.GetAllLiquidityProvidersPartitions(ctx)
		if err != nil {
			fireLPPGetLPsErrorEvent(ctx, err)
		}
	}

	for _, e := range tuples {
		if shouldDistribute {
			lps, exists := partitions[*e.Pool.ExternalAsset]
			if !exists { // TODO: fire event
				k.Logger(ctx).Error(fmt.Sprintf("No liquidity providers for asset %s ", e.Pool.ExternalAsset.Symbol))

				// if this fails, we add rewards to pool instead
				k.addRewardsToPool(ctx, e.Pool, e.Reward)
				continue
			}

			lpsFiltered := FilterValidLiquidityProviders(ctx, lps)
			CollectRewards(ctx, e.Pool, e.Reward, e.Pool.PoolUnits, lpsFiltered, poolRowanMap, lpRowanMap, lpPoolMap)
		} else {
			k.addRewardsToPool(ctx, e.Pool, e.Reward)
		}
	}

	if shouldDistribute {
		poolRowanMapSum := sdk.ZeroUint()
		for _, rowan := range poolRowanMap {
			poolRowanMapSum = poolRowanMapSum.Add(rowan)
		}

		if !coinsToMint.Equal(poolRowanMapSum) {
			k.Logger(ctx).Info(fmt.Sprintln("coinsToMint", coinsToMint.String(), " != poolRowanMapSum", poolRowanMapSum.String()))
		}

		// this updates poolRowanMap in case coin tranfers fails
		k.TransferRewards(ctx, poolRowanMap, lpRowanMap, lpPoolMap)

		for pool, rowan := range poolRowanMap {
			if rowan.Equal(sdk.ZeroUint()) {
				continue
			}

			pool.RewardPeriodNativeDistributed = pool.RewardPeriodNativeDistributed.Add(rowan)
			k.SetPool(ctx, pool) // nolint:errcheck
		}

		// As we have already minted all coins we wanted to distribute, check if we could distribute them all.
		// If not, burn what we could not distribute
		moduleBalancePostTransfer := k.GetModuleRowan(ctx)
		diff := moduleBalancePostTransfer.Sub(moduleBalancePreMinting).Amount // post is always >= pre

		if !diff.IsZero() {
			k.BurnRowan(ctx, diff) // nolint:errcheck
		}
		coinsMinted := sdk.NewIntFromBigInt(coinsToMint.BigInt()).Sub(diff)
		fireRewardsEvent(ctx, "rewards/distribution", coinsMinted, PoolRowanMapToLPPools(poolRowanMap))
	} else {
		fireRewardsEvent(ctx, "rewards/accumulation", sdk.NewIntFromBigInt(coinsToMint.BigInt()), poolRewardsToLPPools(tuples))
	}

	return nil
}

func fireRewardsEvent(ctx sdk.Context, typeStr string, coinsMinted sdk.Int, lpPools []LPPool) {
	data := PrintPools(lpPools)
	successEvent := sdk.NewEvent(
		typeStr,
		sdk.NewAttribute("total_amount", coinsMinted.String()),
		sdk.NewAttribute("amounts", data),
	)

	ctx.EventManager().EmitEvents(sdk.Events{successEvent})
}

func poolRewardsToLPPools(poolRewards []PoolReward) []LPPool {
	arr := make([]LPPool, 0, len(poolRewards))
	for _, poolReward := range poolRewards {
		arr = append(arr, LPPool{Pool: poolReward.Pool, Amount: poolReward.Reward})
	}

	return arr
}

func (k Keeper) GetModuleRowan(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, types.GetCLPModuleAddress(), types.NativeSymbol)
}

func (k Keeper) BurnRowan(ctx sdk.Context, amount sdk.Int) error {
	coin := sdk.NewCoin(types.NativeSymbol, amount)
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
}

func (k Keeper) TransferRewards(ctx sdk.Context, poolRowanMap PoolRowanMap, lpRowanMap LpRowanMap, lpPoolMap LpPoolMap) {
	k.TransferProviderDistributionGeneric(ctx, poolRowanMap, lpRowanMap, lpPoolMap, "rewards/distribution_error", "rewards/distribution")
}

type PoolReward struct {
	Pool   *types.Pool
	Reward sdk.Uint
}

func CollectPoolRewardTuples(pools []*types.Pool, blockDistribution sdk.Uint, totalDepth sdk.Dec, period *types.RewardPeriod) ([]PoolReward, sdk.Uint) {
	coinsToMint := sdk.ZeroUint()
	var tuples []PoolReward //nolint

	remaining := blockDistribution
	for _, pool := range pools {
		if remaining.IsZero() { // we out of money... kthxbye
			return tuples, coinsToMint
		}

		m := GetPoolMultiplier(pool.ExternalAsset.Symbol, period)
		poolDistribution := calcPoolDistribution(m, pool.NativeAssetBalance, totalDepth, blockDistribution)

		if poolDistribution.IsZero() {
			continue
		}

		// This is it: last pool out of pools to get rewards
		if poolDistribution.GT(remaining) {
			poolDistribution = remaining
		}

		coinsToMint = coinsToMint.Add(poolDistribution)
		remaining = remaining.Sub(poolDistribution)
		tuples = append(tuples, PoolReward{Pool: pool, Reward: poolDistribution})
	}

	return tuples, coinsToMint
}

func (k Keeper) addRewardsToPool(ctx sdk.Context, pool *types.Pool, poolDistribution sdk.Uint) {
	pool.NativeAssetBalance = pool.NativeAssetBalance.Add(poolDistribution)
	pool.RewardPeriodNativeDistributed = pool.RewardPeriodNativeDistributed.Add(poolDistribution)

	err := k.SetPool(ctx, pool)
	// this is impossible but... defensive programming
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("The impossible happened: Unable to set pool for asset %s error %s", pool.ExternalAsset.Symbol, err.Error()))
	}
}

func calcPoolDistribution(multiplier sdk.Dec, poolNativeBalance sdk.Uint, totalDepth sdk.Dec, blockDistributionAmount sdk.Uint) sdk.Uint {
	weight := sdk.NewDecFromBigInt(poolNativeBalance.BigInt()).Mul(multiplier).Quo(totalDepth)
	poolDistribution := weight.Mul(sdk.NewDecFromBigInt(blockDistributionAmount.BigInt()))

	return sdk.NewUintFromBigInt(poolDistribution.TruncateInt().BigInt())
}

func CollectRewards(ctx sdk.Context, pool *types.Pool, poolDistribution sdk.Uint, poolUnits sdk.Uint, lps []ValidLiquidityProvider, poolRowanMap PoolRowanMap, lpRowanMap LpRowanMap, lpPoolMap LpPoolMap) {
	rowanToDistribute := CollectProviderDistribution(ctx, pool, sdk.NewDecFromBigInt(poolDistribution.BigInt()), sdk.NewDec(1), poolUnits, lps, lpRowanMap, lpPoolMap)
	poolRowanMap[pool] = rowanToDistribute
}

func (k Keeper) calcTotalDepth(ctx sdk.Context, pools []*types.Pool, period *types.RewardPeriod, height uint64) (sdk.Dec, error) {
	totalDepth := sdk.ZeroDec()
	for _, pool := range pools {
		m := GetPoolMultiplier(pool.ExternalAsset.Symbol, period)
		totalDepth = totalDepth.Add(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()).Mul(m))
		if height == period.RewardPeriodStartBlock {
			pool.RewardPeriodNativeDistributed = sdk.ZeroUint()
			err := k.SetPool(ctx, pool)
			if err != nil {
				return sdk.Dec{}, err
			}
		}
	}

	return totalDepth, nil
}

func (k Keeper) UseUnlockedLiquidity(ctx sdk.Context, lp types.LiquidityProvider, units sdk.Uint, any bool) error {
	// Ensure there is enough liquidity requested for unlock, and also passed lock period.
	// Reduce liquidity in one or more unlock records.
	// Remove unlock records with zero units remaining.
	params := k.GetRewardsParams(ctx)
	currentHeight := ctx.BlockHeight()
	lockPeriod := params.LiquidityRemovalLockPeriod

	unitsLeftToUse := units
	for _, record := range lp.Unlocks {
		if any || record.RequestHeight+int64(lockPeriod) <= currentHeight {
			if unitsLeftToUse.GT(record.Units) {
				// use all this record's unit's and continue with remaining
				unitsLeftToUse = unitsLeftToUse.Sub(record.Units)
				record.Units = sdk.ZeroUint()
			} else {
				// use a portion of this record's units and break
				record.Units = record.Units.Sub(unitsLeftToUse)
				unitsLeftToUse = sdk.ZeroUint()
				break
			}
		}
	}

	if lockPeriod != 0 && !unitsLeftToUse.IsZero() {
		return types.ErrBalanceNotAvailable
	}

	// prune records.
	//var records []*types.LiquidityUnlock
	records := make([]*types.LiquidityUnlock, 0)
	for _, record := range lp.Unlocks {
		/* move to begin blocker
		if currentHeight >= record.RequestHeight + int64(lockPeriod) + cancelPeriod {
			// prune auto cancelled record
			continue
		}*/
		if record.Units.IsZero() {
			// prune used / zero record
			continue
		}
		records = append(records, record)
	}

	lp.Unlocks = records
	k.SetLiquidityProvider(ctx, &lp)

	return nil
}

func (k Keeper) PruneUnlockRecords(ctx sdk.Context, lp *types.LiquidityProvider, lockPeriod, cancelPeriod uint64) {
	currentHeight := ctx.BlockHeight()

	var write bool
	//var records []*types.LiquidityUnlock
	records := make([]*types.LiquidityUnlock, 0)
	for _, record := range lp.Unlocks {
		if currentHeight >= record.RequestHeight+int64(lockPeriod)+int64(cancelPeriod) {
			// prune auto cancelled record
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeCancelUnlock,
					sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
					sdk.NewAttribute(types.AttributeKeyPool, lp.Asset.Symbol),
					sdk.NewAttribute(types.AttributeKeyUnits, record.Units.String()),
				),
			})
			write = true
			continue
		}
		if record.Units.IsZero() {
			// prune used / zero record
			write = true
			continue
		}
		records = append(records, record)
	}
	if write {
		lp.Unlocks = records
		k.SetLiquidityProvider(ctx, lp)
	}
}

func GetPoolMultiplier(asset string, period *types.RewardPeriod) sdk.Dec {
	for _, m := range period.RewardPeriodPoolMultipliers {
		if types.StringCompare(asset, m.PoolMultiplierAsset) {
			if m.Multiplier != nil && !m.Multiplier.IsNil() {
				return *m.Multiplier
			}
		}
	}

	return *period.RewardPeriodDefaultMultiplier
}

func (k Keeper) SetBlockDistributionAccu(ctx sdk.Context, blockDistribution sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	bytes, _ := blockDistribution.Marshal()
	store.Set(types.RewardsBlockDistributionPrefix, bytes)
}

func (k Keeper) GetBlockDistributionAccu(ctx sdk.Context) sdk.Uint {
	blockDistribution := sdk.ZeroUint()
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.RewardsBlockDistributionPrefix)

	if bytes == nil {
		return blockDistribution
	}

	_ = blockDistribution.Unmarshal(bytes)

	return blockDistribution
}
