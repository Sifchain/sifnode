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

func GetAssetSwapFeeParams(asset types.Asset, swapFeeParams *types.SwapFeeParams) (sdk.Dec, sdk.Uint) {

	tokenParams := swapFeeParams.TokenParams
	for _, p := range tokenParams {
		if types.StringCompare(asset.Symbol, p.Asset) {
			return p.SwapFeeRate, p.MinSwapFee
		}
	}

	return swapFeeParams.DefaultSwapFeeRate, sdk.ZeroUint()
}
