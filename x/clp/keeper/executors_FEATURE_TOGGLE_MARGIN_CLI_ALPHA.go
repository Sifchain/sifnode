//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func (k Keeper) ProcessRemoveLiquidityMsg(ctx sdk.Context, msg *types.MsgRemoveLiquidity) (sdk.Int, sdk.Int, sdk.Uint, error) {
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	_, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ExternalAsset.Symbol)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrTokenNotSupported
	}
	pool, err := k.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := k.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrLiquidityProviderDoesNotExist
	}
	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate

	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, pool.String())
	extRowanValue := CalculateWithdrawalRowanValue(withdrawExternalAssetAmount, types.GetSettlementAsset(), pool, pmtpCurrentRunningRate, marginEnabled)

	withdrawExternalAssetAmountInt, ok := k.ParseToInt(withdrawExternalAssetAmount.String())
	if !ok {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrUnableToParseInt
	}
	withdrawNativeAssetAmountInt, ok := k.ParseToInt(withdrawNativeAssetAmount.String())
	if !ok {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrUnableToParseInt
	}
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, withdrawExternalAssetAmountInt)
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, withdrawNativeAssetAmountInt)
	// Subtract Value from pool
	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	// Check if withdrawal makes pool too shallow , checking only for asymetric withdraw.
	if !msg.Asymmetry.IsZero() && (pool.ExternalAssetBalance.IsZero() || pool.NativeAssetBalance.IsZero()) {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), sdkerrors.Wrap(types.ErrPoolTooShallow, "pool balance nil before adjusting asymmetry")
	}
	// Swapping between Native and External based on Asymmetry
	if msg.Asymmetry.IsPositive() {
		marginEnabled := k.getMarginKeeper().IsPoolEnabled(ctx, pool.String())
		swapResult, _, _, swappedPool, err := SwapOne(types.GetSettlementAsset(), swapAmount, *msg.ExternalAsset, pool, pmtpCurrentRunningRate, marginEnabled)
		if err != nil {
			return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapResultInt, ok := k.ParseToInt(swapResult.String())
			if !ok {
				return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.ParseToInt(swapAmount.String())
			if !ok {
				return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrUnableToParseInt
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
			return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapInt, ok := k.ParseToInt(swapResult.String())
			if !ok {
				return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.ParseToInt(swapAmount.String())
			if !ok {
				return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, swapInt)
			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapAmountInt)
			nativeAssetCoin = nativeAssetCoin.Add(swapCoin)
			externalAssetCoin = externalAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	// Check and  remove Liquidity
	err = k.RemoveLiquidity(ctx, pool, externalAssetCoin, nativeAssetCoin, lp, lpUnitsLeft, poolOriginalEB, poolOriginalNB)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroUint(), sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, err.Error())
	}

	return nativeAssetCoin.Amount, externalAssetCoin.Amount, extRowanValue.Add(withdrawNativeAssetAmount), nil
}
