package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CLPCalcSwap(ctx sdk.Context, sentAmount sdk.Uint, to types.Asset, pool types.Pool, marginEnabled bool) (sdk.Uint, error) {

	X, Y, toRowan := pool.ExtractValues(to)

	if marginEnabled {
		X, Y = pool.ExtractDebt(X, Y, toRowan)
	}

	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate

	swapResult := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate)

	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), types.ErrNotEnoughAssetTokens
	}

	return swapResult, nil
}
