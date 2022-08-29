//go:build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

// ------------------------------------------------------------------------------------------------------------------
// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func SwapOne(from types.Asset,
	sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	pmtpCurrentRunningRate, swapFeeRate sdk.Dec) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {

	X, Y, toRowan := pool.ExtractValues(to)

	liquidityFee := CalcLiquidityFee(toRowan, X, sentAmount, Y, swapFeeRate, pmtpCurrentRunningRate)
	priceImpact := calcPriceImpact(X, sentAmount)
	swapResult := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate, swapFeeRate)

	// NOTE: impossible... pre-pmtp at least
	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), types.Pool{}, types.ErrNotEnoughAssetTokens
	}

	pool.UpdateBalances(toRowan, X, sentAmount, Y, swapResult)

	return swapResult, liquidityFee, priceImpact, pool, nil
}

func GetSwapFee(sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	pmtpCurrentRunningRate, swapFeeRate sdk.Dec) sdk.Uint {
	X, Y, toRowan := pool.ExtractValues(to)

	swapResult := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate, swapFeeRate)

	if swapResult.GTE(Y) {
		return sdk.ZeroUint()
	}
	return swapResult
}
