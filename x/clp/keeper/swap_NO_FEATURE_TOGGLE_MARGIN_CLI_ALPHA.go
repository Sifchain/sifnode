//go:build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CLPCalcSwap(ctx sdk.Context, sentAmount sdk.Uint, to types.Asset, pool types.Pool, marginEnabled bool) (sdk.Uint, error) {

	X, Y, toRowan := pool.ExtractValues(to)

	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	swapFeeRate := k.GetSwapFeeRate(ctx).SwapFeeRate

	swapResult := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate, swapFeeRate)

	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), types.ErrNotEnoughAssetTokens
	}

	return swapResult, nil
}
