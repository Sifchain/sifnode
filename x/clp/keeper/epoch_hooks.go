package keeper

import (
	epochstypes "github.com/Sifchain/sifnode/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeforeEpochStart performs a no-op
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {}

// AfterEpochEnd distributes available rewards from rewards bucket to liquidity pools
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) {
	if !k.ShouldDistributeRewards(ctx, epochIdentifier) {
		return
	}
	rewardsEligibleLps, err := k.GetRewardsEligibleLiquidityProviders(ctx)
	if err != nil {
		ctx.Logger().Error("unable to get rewards eligible liquidity providers", "error", err)
		return
	}
	for asset, assetLps := range rewardsEligibleLps {
		// get reward bucket for given asset
		rewardsBucket, found := k.GetRewardsBucket(ctx, asset.Symbol)
		if !found {
			ctx.Logger().Error("unable to get rewards bucket", "asset", asset.Symbol)
			continue
		}
		_ = rewardsBucket
		for _, lp := range assetLps {
			err := k.DistributeLiquidityProviderRewards(ctx, *lp)
			if err != nil {
				ctx.Logger().Error("unable to distribute liquidity provider rewards", "error", err)
			}
		}
	}
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for commitments keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// BeforeEpochStart implements EpochHooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

// AfterEpochEnd implements EpochHooks
func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
