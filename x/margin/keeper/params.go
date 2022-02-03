package keeper

import (
	"strings"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLeverageParam(ctx sdk.Context) sdk.Uint {
	var leverageMax sdk.Uint
	k.paramStore.Get(ctx, types.KeyLeverageMaxParam, &leverageMax)
	return leverageMax
}

func (k Keeper) GetInterestRateMax(ctx sdk.Context) sdk.Dec {
	var d sdk.Dec
	k.paramStore.Get(ctx, types.KeyInterestRateMaxParam, &d)
	return d
}

func (k Keeper) GetInterestRateMin(ctx sdk.Context) sdk.Dec {
	var d sdk.Dec
	k.paramStore.Get(ctx, types.KeyInterestRateMinParam, &d)
	return d
}

func (k Keeper) GetInterestRateIncrease(ctx sdk.Context) sdk.Dec {
	var d sdk.Dec
	k.paramStore.Get(ctx, types.KeyInterestRateIncreaseParam, &d)
	return d
}

func (k Keeper) GetInterestRateDecrease(ctx sdk.Context) sdk.Dec {
	var d sdk.Dec
	k.paramStore.Get(ctx, types.KeyInterestRateDecreaseParam, &d)
	return d
}

func (k Keeper) GetHealthGainFactor(ctx sdk.Context) sdk.Dec {
	var d sdk.Dec
	k.paramStore.Get(ctx, types.KeyHealthGainFactorParam, &d)
	return d
}

func (k Keeper) GetEpochLength(ctx sdk.Context) int64 {
	var d int64
	k.paramStore.Get(ctx, types.KeyEpochLengthParam, &d)
	return d
}

func (k Keeper) GetForceCloseThreshold(ctx sdk.Context) sdk.Dec {
	var d sdk.Dec
	k.paramStore.Get(ctx, types.KeyForceCloseThresholdParam, &d)
	return d
}

func (k Keeper) GetEnabledPools(ctx sdk.Context) []string {
	var pools []string
	k.paramStore.Get(ctx, types.KeyPoolsParam, &pools)
	return pools
}

func (k Keeper) SetEnabledPools(ctx sdk.Context, pools []string) {
	k.paramStore.Set(ctx, types.KeyPoolsParam, &pools)
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
	k.paramStore.SetParamSet(ctx, params)
}
