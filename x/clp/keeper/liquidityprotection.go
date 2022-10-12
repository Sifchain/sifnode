package keeper

import (
	"errors"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetLiquidityProtectionParams(ctx sdk.Context, params *types.LiquidityProtectionParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LiquidityProtectionParamsPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetLiquidityProtectionParams(ctx sdk.Context) *types.LiquidityProtectionParams {
	params := types.LiquidityProtectionParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LiquidityProtectionParamsPrefix)
	k.cdc.MustUnmarshal(bz, &params)
	return &params
}

// This method should only be called if buying or selling native asset.
// If sellNative is false then this method assumes that buyNative is true.
// The nativePrice should be in MaxRowanLiquidityThresholdAsset
// NOTE: this method panics if sellNative is true and the value of the sell amount
// is greater than the value of currentRowanLiquidityThreshold. Call IsBlockedByLiquidityProtection
// before if unsure.
func (k Keeper) MustUpdateLiquidityProtectionThreshold(ctx sdk.Context, sellNative bool, nativeAmount sdk.Uint, nativePrice sdk.Dec) {
	liquidityProtectionParams := k.GetLiquidityProtectionParams(ctx)
	maxRowanLiquidityThreshold := liquidityProtectionParams.MaxRowanLiquidityThreshold
	currentRowanLiquidityThreshold := k.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold

	if liquidityProtectionParams.IsActive {
		nativeValue := CalcRowanValue(nativeAmount, nativePrice)

		var updatedRowanLiquidityThreshold sdk.Uint
		if sellNative {
			if currentRowanLiquidityThreshold.LT(nativeValue) {
				panic(errors.New("expect sell native value to be less than currentRowanLiquidityThreshold"))
			} else {
				updatedRowanLiquidityThreshold = currentRowanLiquidityThreshold.Sub(nativeValue)
			}
		} else {
			// This is equivalent to currentRowanLiquidityThreshold := sdk.MinUint(currentRowanLiquidityThreshold.Add(nativeValue), maxRowanLiquidityThreshold)
			// except it prevents any overflows when adding the nativeValue
			// Assume that maxRowanLiquidityThreshold >= currentRowanLiquidityThreshold
			if maxRowanLiquidityThreshold.Sub(currentRowanLiquidityThreshold).LT(nativeValue) {
				updatedRowanLiquidityThreshold = maxRowanLiquidityThreshold
			} else {
				updatedRowanLiquidityThreshold = currentRowanLiquidityThreshold.Add(nativeValue)
			}
		}

		k.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, updatedRowanLiquidityThreshold)
	}
}

// Currently this calculates the native price on the fly
// Calculates the price of the native token in MaxRowanLiquidityThresholdAsset
func (k Keeper) GetNativePrice(ctx sdk.Context) (sdk.Dec, error) {
	liquidityProtectionParams := k.GetLiquidityProtectionParams(ctx)
	maxRowanLiquidityThresholdAsset := liquidityProtectionParams.MaxRowanLiquidityThresholdAsset

	if types.StringCompare(maxRowanLiquidityThresholdAsset, types.NativeSymbol) {
		return sdk.OneDec(), nil
	}
	pool, err := k.GetPool(ctx, maxRowanLiquidityThresholdAsset)
	if err != nil {
		return sdk.Dec{}, types.ErrMaxRowanLiquidityThresholdAssetPoolDoesNotExist
	}

	return CalcRowanSpotPrice(&pool, k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate)

}

// The nativePrice should be in MaxRowanLiquidityThresholdAsset
func (k Keeper) IsBlockedByLiquidityProtection(ctx sdk.Context, nativeAmount sdk.Uint, nativePrice sdk.Dec) bool {
	value := CalcRowanValue(nativeAmount, nativePrice)
	currentRowanLiquidityThreshold := k.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold
	return currentRowanLiquidityThreshold.LT(value)
}
