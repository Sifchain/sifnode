package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	epochstypes "github.com/Sifchain/sifnode/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeforeEpochStart performs a no-op
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {}

// AfterEpochEnd distributes available rewards from rewards bucket to liquidity providers
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) {
	if !k.ShouldDistributeRewards(ctx, epochIdentifier) {
		return
	}

	rewardsEligibleLps, err := k.GetRewardsEligibleLiquidityProviders(ctx)
	if err != nil {
		ctx.Logger().Error(types.ErrUnableToGetRewardsEligibleLiquidityProviders.Error(), "error", err)
		return
	}

	for asset, assetLps := range rewardsEligibleLps {
		// get reward bucket for given asset
		rewardsBucket, found := k.GetRewardsBucket(ctx, asset.Symbol)
		if !found {
			continue
		}

		k.Logger(ctx).Info(fmt.Sprintf("rewards bucket not found for denom: %s", asset.Symbol))

		rewardShares := k.CalculateRewardShareForLiquidityProviders(ctx, assetLps)
		rewardAmounts := k.CalculateRewardAmountForLiquidityProviders(ctx, rewardShares, rewardsBucket.Amount)

		for i, lp := range assetLps {
			if k.ShouldDistributeRewardsToLPWallet(ctx) {
				err := k.DistributeLiquidityProviderRewards(ctx, lp, asset.Symbol, rewardAmounts[i])
				if err != nil {
					ctx.Logger().Error(types.ErrUnableToDistributeLPRewards.Error(), "error", err)
				}
			} else {
				err := k.AddRewardAmountToLiquidityPool(ctx, lp, asset, rewardAmounts[i])
				if err != nil {
					ctx.Logger().Error(types.ErrUnableToAddRewardAmountToLiquidityPool.Error(), "error", err)
				}
			}

			// increment lp reward amount
			lp.RewardAmount = lp.RewardAmount.Add(sdk.NewCoin(asset.Symbol, rewardAmounts[i]))

			// update the liquidity provider
			k.SetLiquidityProvider(ctx, lp)
		}

		// increment pool reward amount
		pool, err := k.GetPool(ctx, asset.Symbol)
		if err != nil {
			ctx.Logger().Error(types.ErrPoolDoesNotExist.Error(), "error", err)
			continue
		}
		pool.RewardAmountExternal = pool.RewardAmountExternal.Add(sdk.NewUintFromBigInt(rewardsBucket.Amount.BigInt()))
		err = k.SetPool(ctx, &pool)
		if err != nil {
			ctx.Logger().Error(types.ErrUnableToSetPool.Error(), "error", err)
			continue
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
