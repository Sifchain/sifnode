package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPmtpRateParams(ctx sdk.Context, params types.PmtpRateParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PmtpRateParamsPrefix, k.cdc.MustMarshal(&params))
}

func (k Keeper) GetPmtpRateParams(ctx sdk.Context) types.PmtpRateParams {
	params := types.PmtpRateParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PmtpRateParamsPrefix)
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) SetPmtpBlockRate(ctx sdk.Context, blockRate sdk.Dec) {
	currentParams := k.GetPmtpRateParams(ctx)
	currentParams.PmtpPeriodBlockRate = blockRate
	k.SetPmtpRateParams(ctx, currentParams)
}

func (k Keeper) SetPmtpCurrentRunningRate(ctx sdk.Context, runningRate sdk.Dec) {
	currentParams := k.GetPmtpRateParams(ctx)
	currentParams.PmtpCurrentRunningRate = runningRate
	k.SetPmtpRateParams(ctx, currentParams)
}

func (k Keeper) SetPmtpInterPolicyRate(ctx sdk.Context, interPolicyRate sdk.Dec) {
	currentParams := k.GetPmtpRateParams(ctx)
	currentParams.PmtpInterPolicyRate = interPolicyRate
	k.SetPmtpRateParams(ctx, currentParams)
}
