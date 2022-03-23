package keeper

import (
	"strings"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlock prunes / auto-cancels the providers expired unlock requests.
func (k Keeper) BeginBlock(ctx sdk.Context) {
	params := k.GetParams(ctx)
	for it := k.GetLiquidityProviderIterator(ctx); it.Valid(); it.Next() {
		value := it.Value()
		var lp types.LiquidityProvider
		if len(value) <= 0 {
			continue
		}
		err := k.cdc.Unmarshal(value, &lp)
		if err != nil {
			continue
		}
		k.PruneUnlockRecords(ctx, lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)
	}
}

// EndBlock distributes rewards to pools based on the current reward period configurations.
func EndBlock(ctx sdk.Context, _ abci.RequestEndBlock, keeper Keeper) []abci.ValidatorUpdate {
	params := keeper.GetParams(ctx)
	pools := keeper.GetPools(ctx)
	currentPeriod := keeper.GetCurrentRewardPeriod(ctx, params)
	if currentPeriod != nil {
		err := keeper.DistributeDepthRewards(ctx, currentPeriod, pools)
		if err != nil {
			panic(err)
		}
	}

	//keeper.PruneRewardPeriods(ctx, params)

	return []abci.ValidatorUpdate{}
}

func (keeper Keeper) GetCurrentRewardPeriod(ctx sdk.Context, params types.Params) *types.RewardPeriod {
	height := uint64(ctx.BlockHeight())
	for _, period := range params.RewardPeriods {
		if height >= period.StartBlock && height <= period.EndBlock {
			return period
		}
	}
	return nil
}

/*
func (k Keeper) PruneRewardPeriods(ctx sdk.Context, params types.Params) {
	height := uint64(ctx.BlockHeight())
	var write bool
	var periods []*types.RewardPeriod
	for _, period := range params.RewardPeriods {
		if period.EndBlock > height {
			write = true
			continue
		}

		periods = append(periods, period)
	}

	if write {
		params.RewardPeriods = periods
		k.SetParams(ctx, params)
	}
}
*/
func (k Keeper) DistributeDepthRewards(ctx sdk.Context, period *types.RewardPeriod, pools []*types.Pool) error {
	//rewardExecution := k.GetRewardExecution(ctx)
	//if rewardExecution.Id != period.Id {
	//	rewardExecution.Id = period.Id
	//	rewardExecution.Distributed = sdk.ZeroUint()
	//}
	// todo remove unecesssary remaining counter here and in store
	//remaining := period.Allocation.Sub(rewardExecution.Distributed)
	periodLength := period.EndBlock - period.StartBlock
	blockDistribution := period.Allocation.QuoUint64(periodLength)

	if /*remaining.IsZero() ||*/ blockDistribution.IsZero() {
		return nil
	}

	totalDepth := sdk.ZeroDec()
	for _, pool := range pools {
		m := k.GetPoolMultiplier(pool.ExternalAsset.Symbol, period)
		totalDepth = totalDepth.Add(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()).Mul(m))
	}

	for _, pool := range pools {
		m := k.GetPoolMultiplier(pool.ExternalAsset.Symbol, period)
		weight := sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()).Mul(m).Quo(totalDepth)
		blockDistributionDec := sdk.NewDecFromBigInt(blockDistribution.BigInt())
		poolDistributionDec := weight.Mul(blockDistributionDec)
		poolDistribution := sdk.NewUintFromBigInt(poolDistributionDec.TruncateInt().BigInt())
		//if poolDistribution.GT(remaining) {
		//		poolDistribution = remaining
		//}
		if poolDistribution.IsZero() {
			continue
		}
		rewardCoins := sdk.NewCoins(sdk.NewCoin(types.GetSettlementAsset().Symbol, sdk.NewIntFromBigInt(poolDistribution.BigInt())))
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, rewardCoins)
		if err != nil {
			return err
		}
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(poolDistribution)
		//remaining = remaining.Sub(poolDistribution)
		//rewardExecution.Distributed = rewardExecution.Distributed.Add(poolDistribution)
		err = k.SetPool(ctx, pool)
		if err != nil {
			return err
		}
	}

	//k.SetRewardExecution(ctx, rewardExecution)

	return nil
}

func (k Keeper) GetRewardExecution(ctx sdk.Context) types.RewardExecution {
	var rewardExecution types.RewardExecution
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardExecutionPrefix)
	if bz == nil {
		return types.RewardExecution{}
	}
	k.cdc.MustUnmarshal(bz, &rewardExecution)
	return rewardExecution
}

func (k Keeper) SetRewardExecution(ctx sdk.Context, execution types.RewardExecution) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&execution)
	store.Set(types.RewardExecutionPrefix, bz)
}

/*
func (k Keeper) GetRewardsDistributed(ctx sdk.Context) sdk.Uint {
	var rewardExecution types.RewardExecution
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardExecutionPrefix)
	if bz == nil {
		return sdk.ZeroUint()
	}
	k.cdc.MustUnmarshal(bz, &rewardExecution)
	return rewardExecution.Distributed
}

func (k Keeper) SetRewardsDistributed(ctx sdk.Context, distributed sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	rewardsExecution := types.RewardExecution{
		Distributed: distributed,
	}
	bz := k.cdc.MustMarshal(&rewardsExecution)
	store.Set(types.RewardExecutionPrefix, bz)
}
*/

func (k Keeper) UseUnlockedLiquidity(ctx sdk.Context, lp types.LiquidityProvider, units sdk.Uint) error {
	// Ensure there is enough liquidity requested for unlock, and also passed lock period.
	// Reduce liquidity in one or more unlock records.
	// Remove unlock records with zero units remaining.
	params := k.GetParams(ctx)
	currentHeight := ctx.BlockHeight()
	lockPeriod := params.LiquidityRemovalLockPeriod

	unitsLeftToUse := units
	for _, record := range lp.Unlocks {
		if record.RequestHeight+int64(lockPeriod) <= currentHeight {
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

	if !unitsLeftToUse.IsZero() {
		return types.ErrBalanceNotAvailable
	}

	// prune records.
	var records []*types.LiquidityUnlock
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

func (k Keeper) PruneUnlockRecords(ctx sdk.Context, lp types.LiquidityProvider, lockPeriod, cancelPeriod uint64) {
	currentHeight := ctx.BlockHeight()

	var write bool
	var records []*types.LiquidityUnlock
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
		k.SetLiquidityProvider(ctx, &lp)
	}
}

func (k Keeper) GetPoolMultiplier(asset string, period *types.RewardPeriod) sdk.Dec {
	for _, m := range period.Multipliers {
		if strings.EqualFold(asset, m.Asset) {
			if m.Multiplier != nil && !m.Multiplier.IsNil() {
				return *m.Multiplier
			}
		}
	}

	return sdk.NewDec(1)
}
