package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSwapFeeParams(ctx sdk.Context, params *types.SwapFeeParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.SwapFeeParamsPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetSwapFeeParams(ctx sdk.Context) types.SwapFeeParams {
	params := types.SwapFeeParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SwapFeeParamsPrefix)
	if bz == nil {
		return types.SwapFeeParams{DefaultSwapFeeRate: sdk.NewDecWithPrec(3, 3)} //0.003
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) GetSwapFeeRate(ctx sdk.Context, asset types.Asset) sdk.Dec {

	params := k.GetSwapFeeParams(ctx)

	tokenParams := params.TokenParams
	for _, p := range tokenParams {
		if types.StringCompare(asset.Symbol, p.Asset) {
			return p.SwapFeeRate
		}
	}

	return params.DefaultSwapFeeRate
}
