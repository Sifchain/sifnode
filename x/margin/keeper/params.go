package keeper

import (
	"errors"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k Keeper) GetPoolOpenThreshold(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).PoolOpenThreshold
}

func (k Keeper) GetRemovalQueueThreshold(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).RemovalQueueThreshold
}

func (k Keeper) GetForceCloseFundPercentage(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).ForceCloseFundPercentage
}

func (k Keeper) GetForceCloseFundAddress(ctx sdk.Context) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(k.GetParams(ctx).ForceCloseFundAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (k Keeper) GetIncrementalInterestPaymentFundPercentage(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).IncrementalInterestPaymentFundPercentage
}

func (k Keeper) GetIncrementalInterestPaymentFundAddress(ctx sdk.Context) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(k.GetParams(ctx).IncrementalInterestPaymentFundAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (k Keeper) GetMaxOpenPositions(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).MaxOpenPositions
}

func (k Keeper) GetIncrementalInterestPaymentEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).IncrementalInterestPaymentEnabled
}
func (k Keeper) GetSafetyFactor(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SafetyFactor
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
		if types.StringCompare(p, asset) {
			return true
		}
	}

	return false
}

func (k Keeper) IsPoolClosed(ctx sdk.Context, asset string) bool {
	params := k.GetParams(ctx)
	for _, p := range params.ClosedPools {
		if types.StringCompare(p, asset) {
			return true
		}
	}

	return false
}

func (k Keeper) GetSqModifier(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SqModifier
}

func (k Keeper) IsWhitelistingEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).WhitelistingEnabled
}

func (k Keeper) IsRowanCollateralEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).RowanCollateralEnabled
}

func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	err := ValidateParams(params)
	if err != nil {
		panic(err)
	}
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

func ValidateParams(params *types.Params) error {
	if params.SqModifier.IsNil() || params.SqModifier.IsZero() {
		return sdkerrors.Wrap(errors.New("invalid valid"), "sq modifier must be > 0")
	}

	if params.LeverageMax.IsNegative() {
		return sdkerrors.Wrap(errors.New("invalid value"), "leverage max must be >= 0")
	}

	return nil
}
