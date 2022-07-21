package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CLPCalcSwap(ctx sdk.Context, sentAmount sdk.Uint, to types.Asset, pool types.Pool) (sdk.Uint, error) {

	normalizationFactor, adjustExternalToken, err := k.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	X, x, Y, toRowan := SetInputs(sentAmount, to, pool)

	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate

	swapResult, _ := CalcSwapResult(toRowan, normalizationFactor, adjustExternalToken, X, x, Y, pmtpCurrentRunningRate)

	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), types.ErrNotEnoughAssetTokens
	}

	return swapResult, nil
}
