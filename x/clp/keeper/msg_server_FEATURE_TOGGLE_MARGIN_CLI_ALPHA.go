//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"fmt"

	"context"

	"strconv"
	"strings"

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

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var (
		priceImpact sdk.Uint
	)
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	sAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.SentAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	rAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ReceivedAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(sAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(rAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	if k.tokenRegistryKeeper.CheckEntryPermissions(sAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_DISABLE_SELL}) {
		return nil, tokenregistrytypes.ErrNotAllowedToSellAsset
	}
	if k.tokenRegistryKeeper.CheckEntryPermissions(rAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_DISABLE_BUY}) {
		return nil, tokenregistrytypes.ErrNotAllowedToBuyAsset
	}

	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate

	liquidityProtectionParams := k.GetLiquidityProtectionParams(ctx)
	maxRowanLiquidityThreshold := liquidityProtectionParams.MaxRowanLiquidityThreshold
	maxRowanLiquidityThresholdAsset := liquidityProtectionParams.MaxRowanLiquidityThresholdAsset
	currentRowanLiquidityThreshold := k.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold
	var (
		sentValue sdk.Uint
	)

	// if liquidity protection is active and selling rowan
	if liquidityProtectionParams.IsActive && strings.EqualFold(sAsset.Denom, types.NativeSymbol) {
		if strings.EqualFold(maxRowanLiquidityThresholdAsset, types.NativeSymbol) {
			sentValue = msg.SentAmount
		} else {
			pool, err := k.GetPool(ctx, maxRowanLiquidityThresholdAsset)
			if err != nil {
				return nil, types.ErrMaxRowanLiquidityThresholdAssetPoolDoesNotExist
			}
			sentValue, err = CalcRowanValue(&pool, pmtpCurrentRunningRate, msg.SentAmount)

			if err != nil {
				return nil, err
			}
		}

		if currentRowanLiquidityThreshold.LT(sentValue) {
			return nil, types.ErrReachedMaxRowanLiquidityThreshold
		}
	}

	liquidityFeeNative := sdk.ZeroUint()
	liquidityFeeExternal := sdk.ZeroUint()
	totalLiquidityFee := sdk.ZeroUint()
	priceImpact = sdk.ZeroUint()
	sentAmount := msg.SentAmount
	sentAsset := msg.SentAsset
	receivedAsset := msg.ReceivedAsset
	// Get native asset
	nativeAsset := types.GetSettlementAsset()
	inPool, outPool := types.Pool{}, types.Pool{}
	// If sending rowan ,deduct directly from the Native balance  instead of fetching from rowan pool
	if !msg.SentAsset.Equals(types.GetSettlementAsset()) {
		inPool, err = k.Keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	}
	sentAmountInt, ok := k.Keeper.ParseToInt(sentAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	sentCoin := sdk.NewCoin(msg.SentAsset.Symbol, sentAmountInt)
	err = k.Keeper.InitiateSwap(ctx, sentCoin, accAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	// Check if its a two way swap, swapping non native fro non native .
	// If its one way we can skip this if condition and add balance to users account from outpool
	if !msg.SentAsset.Equals(nativeAsset) && !msg.ReceivedAsset.Equals(nativeAsset) {
		marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, inPool.String())
		emitAmount, lp, ts, finalPool, err := SwapOne(*sentAsset, sentAmount, nativeAsset, inPool, pmtpCurrentRunningRate, marginEnabled)
		if err != nil {
			return nil, err
		}
		err = k.Keeper.SetPool(ctx, &finalPool)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
		sentAmount = emitAmount
		sentAsset = &nativeAsset
		priceImpact = priceImpact.Add(ts)
		liquidityFeeNative = liquidityFeeNative.Add(lp)
	}
	// If receiving  rowan , add directly to  Native balance  instead of fetching from rowan pool
	if msg.ReceivedAsset.Equals(types.GetSettlementAsset()) {
		outPool, err = k.Keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	} else {
		outPool, err = k.Keeper.GetPool(ctx, msg.ReceivedAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
		}
	}
	// Calculating amount user receives
	marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, inPool.String())
	emitAmount, lp, ts, finalPool, err := SwapOne(*sentAsset, sentAmount, *receivedAsset, outPool, pmtpCurrentRunningRate, marginEnabled)
	if err != nil {
		return nil, err
	}
	if emitAmount.LT(msg.MinReceivingAmount) {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSwapFailed,
				sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
				sdk.NewAttribute(types.AttributeKeyThreshold, msg.MinReceivingAmount.String()),
				sdk.NewAttribute(types.AttributeKeyInPool, inPool.String()),
				sdk.NewAttribute(types.AttributeKeyOutPool, outPool.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
		})
		return &types.MsgSwapResponse{}, types.ErrReceivedAmountBelowExpected
	}
	// todo nil pointer deref test
	err = k.Keeper.FinalizeSwap(ctx, emitAmount.String(), finalPool, *msg)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	if liquidityFeeNative.GT(sdk.ZeroUint()) {
		liquidityFeeExternal = liquidityFeeExternal.Add(lp)
		marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, outPool.String())
		firstSwapFeeInOutputAsset := GetSwapFee(liquidityFeeNative, *msg.ReceivedAsset, outPool, pmtpCurrentRunningRate, marginEnabled)
		totalLiquidityFee = liquidityFeeExternal.Add(firstSwapFeeInOutputAsset)
	} else {
		totalLiquidityFee = liquidityFeeNative.Add(lp)
	}
	priceImpact = priceImpact.Add(ts)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
			sdk.NewAttribute(types.AttributeKeyLiquidityFee, totalLiquidityFee.String()),
			sdk.NewAttribute(types.AttributeKeyPriceImpact, priceImpact.String()),
			sdk.NewAttribute(types.AttributeKeyInPool, inPool.String()),
			sdk.NewAttribute(types.AttributeKeyOutPool, outPool.String()),
			sdk.NewAttribute(types.AttributePmtpBlockRate, k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate.String()),
			sdk.NewAttribute(types.AttributePmtpCurrentRunningRate, pmtpCurrentRunningRate.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

	if liquidityProtectionParams.IsActive {
		// if sell rowan
		if strings.EqualFold(sAsset.Denom, types.NativeSymbol) {
			// we know that sentValue < currentRowanLiquidityThreshold so we can do the
			// substitution knowing it won't panic
			currentRowanLiquidityThreshold = currentRowanLiquidityThreshold.Sub(sentValue)
			k.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, currentRowanLiquidityThreshold)
		}

		// if buy rowan
		if strings.EqualFold(rAsset.Denom, types.NativeSymbol) {
			var emitValue sdk.Uint
			if strings.EqualFold(maxRowanLiquidityThresholdAsset, types.NativeSymbol) {
				emitValue = emitAmount
			} else {
				pool, err := k.GetPool(ctx, maxRowanLiquidityThresholdAsset)
				if err != nil {
					return nil, types.ErrMaxRowanLiquidityThresholdAssetPoolDoesNotExist
				}
				emitValue, err = CalcRowanValue(&pool, pmtpCurrentRunningRate, emitAmount)

				if err != nil {
					return nil, err
				}
			}

			// This is equivalent to currentRowanLiquidityThreshold := sdk.MinUint(currentRowanLiquidityThreshold.Add(emitValue), maxRowanLiquidityThreshold)
			// except it prevents any overflows when adding the emitValue
			if maxRowanLiquidityThreshold.Sub(currentRowanLiquidityThreshold).LT(emitValue) {
				currentRowanLiquidityThreshold = maxRowanLiquidityThreshold
			} else {
				currentRowanLiquidityThreshold = currentRowanLiquidityThreshold.Add(emitValue)
			}

			k.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, currentRowanLiquidityThreshold)
		}
	}

	return &types.MsgSwapResponse{}, nil
}

func (k msgServer) RemoveLiquidity(goCtx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(eAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	pool, err := k.Keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := k.Keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer)

	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}

	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	// Prune pools
	params := k.GetRewardsParams(ctx)
	k.PruneUnlockRecords(ctx, &lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)

	if !msg.Asymmetry.IsZero() {
		return nil, types.ErrAsymmetricRemove
	}

	err = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_VerifyEnoughWBasisPointsAvailableForLP(ctx, k, msg, lp)

	if err != nil {
		return nil, err
	}

	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	err = k.Keeper.UseUnlockedLiquidity(ctx, lp, lp.LiquidityProviderUnits.Sub(lpUnitsLeft), false)
	if err != nil {
		return nil, err
	}

	err = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_QueueRemovalWithWBasisPoints(ctx, k, msg, lp, pool, withdrawNativeAssetAmount, withdrawExternalAssetAmount, eAsset, pmtpCurrentRunningRate)
	if err != nil {
		return nil, err
	}

	withdrawExternalAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	withdrawNativeAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawNativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, withdrawExternalAssetAmountInt)
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, withdrawNativeAssetAmountInt)
	// Subtract Value from pool
	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	// Check if withdrawal makes pool too shallow , checking only for asymetric withdraw.
	if !msg.Asymmetry.IsZero() && (pool.ExternalAssetBalance.IsZero() || pool.NativeAssetBalance.IsZero()) {
		return nil, sdkerrors.Wrap(types.ErrPoolTooShallow, "pool balance nil before adjusting asymmetry")

	}
	// Swapping between Native and External based on Asymmetry
	if msg.Asymmetry.IsPositive() {
		marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, pool.String())
		swapResult, _, _, swappedPool, err := SwapOne(types.GetSettlementAsset(), swapAmount, *msg.ExternalAsset, pool, pmtpCurrentRunningRate, marginEnabled)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapResultInt, ok := k.Keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.Keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapResultInt)
			swapAmountInCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, swapAmountInt)
			externalAssetCoin = externalAssetCoin.Add(swapCoin)
			nativeAssetCoin = nativeAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	if msg.Asymmetry.IsNegative() {
		marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, pool.String())
		swapResult, _, _, swappedPool, err := SwapOne(*msg.ExternalAsset, swapAmount, types.GetSettlementAsset(), pool, pmtpCurrentRunningRate, marginEnabled)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapInt, ok := k.Keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.Keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, swapInt)
			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapAmountInt)
			nativeAssetCoin = nativeAssetCoin.Add(swapCoin)
			externalAssetCoin = externalAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	// Check and  remove Liquidity
	err = k.Keeper.RemoveLiquidity(ctx, pool, externalAssetCoin, nativeAssetCoin, lp, lpUnitsLeft, poolOriginalEB, poolOriginalNB)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyUnits, lp.LiquidityProviderUnits.Sub(lpUnitsLeft).String()),
			sdk.NewAttribute(types.AttributePmtpBlockRate, k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate.String()),
			sdk.NewAttribute(types.AttributePmtpCurrentRunningRate, pmtpCurrentRunningRate.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgRemoveLiquidityResponse{}, nil
}
