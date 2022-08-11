//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package clp_test

import (
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SwapOne(ctx sdk.Context,
	k clpkeeper.Keeper,
	sentAsset types.Asset,
	sentAmount sdk.Uint,
	nativeAsset types.Asset,
	inPool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {
	marginEnabled := k.GetMarginKeeper().IsPoolEnabled(ctx, inPool.String())
	return clpkeeper.SwapOne(sentAsset, sentAmount, nativeAsset, inPool, pmtpCurrentRunningRate, marginEnabled)
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_GetSwapFee(ctx sdk.Context,
	k clpkeeper.Keeper,
	ReceivedAsset *types.Asset,
	liquidityFeeNative sdk.Uint,
	outPool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) sdk.Uint {
	marginEnabled := k.GetMarginKeeper().IsPoolEnabled(ctx, outPool.String())
	return clpkeeper.GetSwapFee(liquidityFeeNative, *ReceivedAsset, outPool, pmtpCurrentRunningRate, marginEnabled)
}
