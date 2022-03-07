package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

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

	keeper.PruneRewardPeriods(ctx, params)

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

func (k Keeper) DistributeDepthRewards(ctx sdk.Context, period *types.RewardPeriod, pools []*types.Pool) error {
	distributed := k.GetRewardsDistributed(ctx)
	remaining := period.Allocation.Sub(distributed)
	periodLength := period.EndBlock - period.StartBlock
	blockDistribution := remaining.QuoUint64(periodLength)

	if remaining.IsZero() || blockDistribution.IsZero() {
		return nil
	}

	totalDepth := sdk.ZeroUint()
	for _, pool := range pools {
		totalDepth = totalDepth.Add(pool.NativeAssetBalance)
	}

	for _, pool := range pools {
		weight := pool.NativeAssetBalance.Quo(totalDepth)
		poolDistribution := blockDistribution.Mul(weight)
		if poolDistribution.GT(remaining) {
			poolDistribution = remaining
		}
		rewardCoins := sdk.NewCoins(sdk.NewCoin(types.GetSettlementAsset().Symbol, sdk.NewIntFromUint64(poolDistribution.Uint64())))
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, rewardCoins)
		if err != nil {
			return err
		}
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(poolDistribution)
		remaining = remaining.Sub(poolDistribution)
		distributed = distributed.Add(poolDistribution)
		err = k.SetPool(ctx, pool)
		if err != nil {
			return err
		}
	}

	k.SetRewardsDistributed(ctx, distributed)

	return nil
}

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
