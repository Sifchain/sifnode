//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"strings"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetMaxLeverageParam(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).LeverageMax
}

func (k Keeper) GetInterestRateMax(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateMax
}

func (k Keeper) GetInterestRateMin(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateMin
}

func (k Keeper) GetInterestRateIncrease(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateIncrease
}

func (k Keeper) GetInterestRateDecrease(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateDecrease
}

func (k Keeper) GetHealthGainFactor(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).HealthGainFactor
}

func (k Keeper) GetEpochLength(ctx sdk.Context) int64 {
	return k.GetParams(ctx).EpochLength
}

func (k Keeper) GetForceCloseThreshold(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).ForceCloseThreshold
}

func (k Keeper) GetRemovalQueueThreshold(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).RemovalQueueThreshold
}

func (k Keeper) GetMaxOpenPositions(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).MaxOpenPositions
}

func (k Keeper) GetEnabledPools(ctx sdk.Context) []string {
	return k.GetParams(ctx).Pools
}

func (k Keeper) SetEnabledPools(ctx sdk.Context, pools []string) {
	params := k.GetParams(ctx)
	params.Pools = pools
	k.SetParams(ctx, &params)
}

func (k Keeper) IsPoolEnabled(ctx sdk.Context, asset string) bool {
	pools := k.GetEnabledPools(ctx)
	for _, p := range pools {
		if strings.EqualFold(p, asset) {
			return true
		}
	}

	return false
}

func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ParamsPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsPrefix)
	if bz == nil {
		return *types.DefaultGenesis().Params
	}
	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return params
}
