//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_ProcessRemovelQueue(ctx sdk.Context, k msgServer, msg *types.MsgAddLiquidity, newPoolUnits sdk.Uint) {
	if k.GetRemovalQueue(ctx, msg.ExternalAsset.Symbol).Count > 0 {
		k.ProcessRemovalQueue(ctx, msg, newPoolUnits)
	}
}

//  ensure requested removal amount is less than available - what is already on the queue
func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_VerifyEnoughWithdrawUnitsAvailableForLP(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidityUnits, lp types.LiquidityProvider) error {
	lpQueuedUnits := k.GetRemovalQueueUnitsForLP(ctx, lp)
	if msg.WithdrawUnits.GT(lp.LiquidityProviderUnits.Sub(lpQueuedUnits)) {
		return sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, fmt.Sprintf("WithdrawUnits %s greater than total LP units %s minus queued removals", msg.WithdrawUnits, lp.LiquidityProviderUnits))
	}
	return nil
}

//  ensure requested removal amount is less than available - what is already on the queue
func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_VerifyEnoughWBasisPointsAvailableForLP(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidity, lp types.LiquidityProvider) error {
	lpQueuedUnits := k.GetRemovalQueueUnitsForLP(ctx, lp)
	msgUnits := ConvWBasisPointsToUnits(lp.LiquidityProviderUnits, msg.WBasisPoints)
	if msgUnits.GT(lp.LiquidityProviderUnits.Sub(lpQueuedUnits)) {
		return sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, fmt.Sprintf("WithdrawUnits %s greater than total LP units %s minus queued removals", msgUnits, lp.LiquidityProviderUnits))
	}
	return nil
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_QueueRemovalWithWithdrawUnits(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidityUnits, lp types.LiquidityProvider, pool types.Pool, withdrawNativeAssetAmount, withdrawExternalAssetAmount sdk.Uint, eAsset *tokenregistrytypes.RegistryEntry, pmtpCurrentRunningRate sdk.Dec) error {
	marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, pool.String())
	extRowanValue := CalculateWithdrawalRowanValue(withdrawExternalAssetAmount, types.GetSettlementAsset(), pool, pmtpCurrentRunningRate, marginEnabled)

	futurePool := pool
	futurePool.NativeAssetBalance = futurePool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	futurePool.ExternalAssetBalance = futurePool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	if k.GetMarginKeeper().CalculatePoolHealth(&futurePool).LT(k.GetMarginKeeper().GetRemovalQueueThreshold(ctx)) {
		k.QueueRemoval(ctx, &types.MsgRemoveLiquidity{
			Signer:        msg.Signer,
			ExternalAsset: msg.ExternalAsset,
			WBasisPoints:  ConvUnitsToWBasisPoints(lp.LiquidityProviderUnits, msg.WithdrawUnits),
			Asymmetry:     sdk.ZeroInt(),
		}, extRowanValue.Add(withdrawNativeAssetAmount))
		return types.ErrQueued
	}
	return nil
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_QueueRemovalWithWBasisPoints(ctx sdk.Context, k msgServer, msg *types.MsgRemoveLiquidity, lp types.LiquidityProvider, pool types.Pool, withdrawNativeAssetAmount, withdrawExternalAssetAmount sdk.Uint, eAsset *tokenregistrytypes.RegistryEntry, pmtpCurrentRunningRate sdk.Dec) error {
	marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, pool.String())
	extRowanValue := CalculateWithdrawalRowanValue(withdrawExternalAssetAmount, types.GetSettlementAsset(), pool, pmtpCurrentRunningRate, marginEnabled)

	futurePool := pool
	futurePool.NativeAssetBalance = futurePool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	futurePool.ExternalAssetBalance = futurePool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	if k.GetMarginKeeper().CalculatePoolHealth(&futurePool).LT(k.GetMarginKeeper().GetRemovalQueueThreshold(ctx)) {
		k.QueueRemoval(ctx, msg, extRowanValue.Add(withdrawExternalAssetAmount))
		return types.ErrQueued
	}

	return nil
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SwapOne(ctx sdk.Context,
	k msgServer,
	sentAsset types.Asset,
	sentAmount sdk.Uint,
	nativeAsset types.Asset,
	inPool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {
	marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, inPool.String())
	return SwapOne(sentAsset, sentAmount, nativeAsset, inPool, pmtpCurrentRunningRate, marginEnabled)
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_GetSwapFee(ctx sdk.Context,
	k msgServer,
	ReceivedAsset *types.Asset,
	liquidityFeeNative sdk.Uint,
	outPool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) sdk.Uint {
	marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, outPool.String())
	return GetSwapFee(liquidityFeeNative, *ReceivedAsset, outPool, pmtpCurrentRunningRate, marginEnabled)
}
