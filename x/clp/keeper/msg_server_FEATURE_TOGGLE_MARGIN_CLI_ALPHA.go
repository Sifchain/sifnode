//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"context"

	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

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

	//  ensure requested removal amount is less than available - what is already on the queue
	lpQueuedUnits := k.GetRemovalQueueUnitsForLP(ctx, lp)
	msgUnits := ConvWBasisPointsToUnits(lp.LiquidityProviderUnits, msg.WBasisPoints)
	if msgUnits.GT(lp.LiquidityProviderUnits.Sub(lpQueuedUnits)) {
		return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, fmt.Sprintf("WithdrawUnits %s greater than total LP units %s minus queued removals", msgUnits, lp.LiquidityProviderUnits))
	}

	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	normalizationFactor, adjustExternalToken := k.GetNormalizationFactor(eAsset.Decimals)
	extRowanValue, err := CalculateWithdrawalRowanValue(withdrawExternalAssetAmount, types.GetSettlementAsset(), pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.UseUnlockedLiquidity(ctx, lp, lp.LiquidityProviderUnits.Sub(lpUnitsLeft), false)
	if err != nil {
		return nil, err
	}

	futurePool := pool
	futurePool.NativeAssetBalance = futurePool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	futurePool.ExternalAssetBalance = futurePool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	if k.GetMarginKeeper().CalculatePoolHealth(&futurePool).LT(k.GetMarginKeeper().GetRemovalQueueThreshold(ctx)) {
		k.QueueRemoval(ctx, msg, extRowanValue.Add(withdrawExternalAssetAmount))
		return nil, types.ErrQueued
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
		normalizationFactor, adjustExternalToken := k.GetNormalizationFactor(eAsset.Decimals)
		swapResult, _, _, swappedPool, err := SwapOne(types.GetSettlementAsset(), swapAmount, *msg.ExternalAsset, pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)
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
		normalizationFactor, adjustExternalToken := k.GetNormalizationFactor(eAsset.Decimals)
		swapResult, _, _, swappedPool, err := SwapOne(*msg.ExternalAsset, swapAmount, types.GetSettlementAsset(), pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)
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

func (k msgServer) RemoveLiquidityUnits(goCtx context.Context, msg *types.MsgRemoveLiquidityUnits) (*types.MsgRemoveLiquidityUnitsResponse, error) {
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
	//  ensure requested removal amount is less than available - what is already on the queue
	lpQueuedUnits := k.GetRemovalQueueUnitsForLP(ctx, lp)
	if msg.WithdrawUnits.GT(lp.LiquidityProviderUnits.Sub(lpQueuedUnits)) {
		return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, fmt.Sprintf("WithdrawUnits %s greater than total LP units %s minus queued removals", msg.WithdrawUnits, lp.LiquidityProviderUnits))
	}

	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	// Prune pools
	params := k.GetRewardsParams(ctx)
	k.PruneUnlockRecords(ctx, &lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)

	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft := CalculateWithdrawalFromUnits(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WithdrawUnits)

	normalizationFactor, adjustExternalToken := k.GetNormalizationFactor(eAsset.Decimals)
	extRowanValue, err := CalculateWithdrawalRowanValue(withdrawExternalAssetAmount, types.GetSettlementAsset(), pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.UseUnlockedLiquidity(ctx, lp, lp.LiquidityProviderUnits.Sub(lpUnitsLeft), false)
	if err != nil {
		return nil, err
	}

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
		return nil, types.ErrQueued
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
	return &types.MsgRemoveLiquidityUnitsResponse{}, nil
}

func (k msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(eAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	// Get pool
	pool, err := k.Keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	normalizationFactor, adjustExternalToken := k.GetNormalizationFactor(eAsset.Decimals)
	symmetryThreshold := k.GetSymmetryThreshold(ctx)
	newPoolUnits, lpUnits, err := CalculatePoolUnits(
		pool.PoolUnits,
		pool.NativeAssetBalance,
		pool.ExternalAssetBalance,
		msg.NativeAssetAmount,
		msg.ExternalAssetAmount,
		normalizationFactor,
		adjustExternalToken,
		symmetryThreshold)
	if err != nil {
		return nil, err
	}
	// Get lp , if lp doesnt exist create lp
	lp, err := k.Keeper.AddLiquidity(ctx, msg, pool, newPoolUnits, lpUnits)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToAddLiquidity, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyUnits, lpUnits.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

	if k.GetRemovalQueue(ctx).Count > 0 {
		k.ProcessRemovalQueue(ctx, msg, newPoolUnits)
	}

	return &types.MsgAddLiquidityResponse{}, nil
}
