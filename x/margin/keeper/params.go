package keeper

import (
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

func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	k.paramStore.SetParamSet(ctx, params)
}
