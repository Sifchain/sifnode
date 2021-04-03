package clp

import (
	"fmt"
	"strconv"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

// NewHandler creates an sdk.Handler for all the clp type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgCreatePool:
			res, err := msgServer.CreatePool(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDecommissionPool:
			res, err := msgServer.DecommissionPool(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddLiquidity:
			res, err := msgServer.AddLiquidity(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveLiquidity:
			res, err := msgServer.RemoveLiquidity(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSwap:
			res, err := msgServer.Swap(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDecommissionPool(ctx sdk.Context, keeper Keeper, msg MsgDecommissionPool) (*sdk.Result, error) {
	// Verify pool
	pool, err := keeper.GetPool(ctx, msg.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	if !keeper.ValidateAddress(ctx, msg.Signer) {
		return nil, errors.Wrap(types.ErrInvalid, "user does not have permission to decommission pool")
	}
	if pool.NativeAssetBalance.GTE(sdk.NewUintFromString(PoolThrehold)) {
		return nil, types.ErrBalanceTooHigh
	}
	// Get all LP's for the pool
	lpList := keeper.GetLiquidityProvidersForAsset(ctx, pool.ExternalAsset)
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	// iterate over Lp list and refund them there tokens
	// Return both RWN and EXTERNAL ASSET
	for _, lp := range lpList {
		withdrawNativeAsset, withdrawExternalAsset, _, _ := clpkeeper.CalculateAllAssetsForLP(pool, lp)
		poolUnits = poolUnits.Sub(lp.LiquidityProviderUnits)
		nativeAssetBalance = nativeAssetBalance.Sub(withdrawNativeAsset)
		externalAssetBalance = externalAssetBalance.Sub(withdrawExternalAsset)

		withdrawNativeAssetInt, ok := keeper.ParseToInt(withdrawNativeAsset.String())
		if !ok {
			return nil, types.ErrUnableToParseInt
		}
		withdrawExternalAssetInt, ok := keeper.ParseToInt(withdrawExternalAsset.String())
		if !ok {
			return nil, types.ErrUnableToParseInt
		}
		withdrawNativeCoins := sdk.NewCoin(GetSettlementAsset().Symbol, withdrawNativeAssetInt)
		withdrawExternalCoins := sdk.NewCoin(msg.Symbol, withdrawExternalAssetInt)
		refundingCoins := sdk.Coins{withdrawExternalCoins, withdrawNativeCoins}
		err := keeper.RemoveLiquidityProvider(ctx, refundingCoins, lp)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToRemoveLiquidityProvider, err.Error())
		}
	}
	// Pool should be empty at this point
	// Decommission the pool
	err = keeper.DecommissionPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToDecommissionPool, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDecommissionPool,
			sdk.NewAttribute(types.AttributeKeyPool, pool.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreatePool(ctx sdk.Context, keeper Keeper, msg MsgCreatePool) (*sdk.Result, error) {
	// Verify min threshold

	MinThreshold := sdk.NewUintFromString(PoolThrehold)

	if msg.NativeAssetAmount.LT(MinThreshold) { // Need to verify
		return nil, types.ErrTotalAmountTooLow
	}
	// Check if pool already exists
	if keeper.ExistsPool(ctx, msg.ExternalAsset.Symbol) {
		return nil, types.ErrUnableToCreatePool
	}

	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	poolUnits, lpunits, err := clpkeeper.CalculatePoolUnits(msg.ExternalAsset.Symbol, sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), nativeBalance, externalBalance)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}
	// Create Pool
	pool, err := keeper.CreatePool(ctx, poolUnits, msg)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Create Liquidity Provider
	lp := keeper.CreateLiquidityProvider(ctx, msg.ExternalAsset, lpunits, msg.Signer)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyPool, pool.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			types.EventTypeCreateLiquidityProvider,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgAddLiquidity(ctx sdk.Context, keeper Keeper, msg MsgAddLiquidity) (*sdk.Result, error) {
	// Get pool
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}

	newPoolUnits, lpUnits, err := clpkeeper.CalculatePoolUnits(
		msg.ExternalAsset.Symbol,
		pool.PoolUnits,
		pool.NativeAssetBalance,
		pool.ExternalAssetBalance,
		msg.NativeAssetAmount,
		msg.ExternalAssetAmount)
	if err != nil {
		return nil, err
	}

	// Get lp , if lp doesnt exist create lp

	lp, err := keeper.AddLiquidity(ctx, msg, pool, newPoolUnits, lpUnits)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToAddLiquidity, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveLiquidity(ctx sdk.Context, keeper Keeper, msg MsgRemoveLiquidity) (*sdk.Result, error) {
	// Get pool
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer.String())
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}
	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	withdrawExternalAssetAmountInt, ok := keeper.ParseToInt(withdrawExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	withdrawNativeAssetAmountInt, ok := keeper.ParseToInt(withdrawNativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, withdrawExternalAssetAmountInt)
	nativeAssetCoin := sdk.NewCoin(GetSettlementAsset().Symbol, withdrawNativeAssetAmountInt)

	// Subtract Value from pool
	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	// Check if withdrawal makes pool too shallow , checking only for asymetric withdraw.

	if !msg.Asymmetry.IsZero() && (pool.ExternalAssetBalance.IsZero() || pool.NativeAssetBalance.IsZero()) {
		return nil, errors.Wrap(types.ErrPoolTooShallow, "pool balance nil before adjusting asymmetry")
	}

	// Swapping between Native and External based on Asymmetry
	if msg.Asymmetry.IsPositive() {
		swapResult, _, _, swappedPool, err := clpkeeper.SwapOne(GetSettlementAsset(), swapAmount, msg.ExternalAsset, pool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapResultInt, ok := keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapResultInt)
			swapAmountInCoin := sdk.NewCoin(GetSettlementAsset().Symbol, swapAmountInt)
			externalAssetCoin = externalAssetCoin.Add(swapCoin)
			nativeAssetCoin = nativeAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	if msg.Asymmetry.IsNegative() {
		swapResult, _, _, swappedPool, err := clpkeeper.SwapOne(msg.ExternalAsset, swapAmount, GetSettlementAsset(), pool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapInt, ok := keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(GetSettlementAsset().Symbol, swapInt)
			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapAmountInt)

			nativeAssetCoin = nativeAssetCoin.Add(swapCoin)
			externalAssetCoin = externalAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	// Check and  remove Liquidity
	err = keeper.RemoveLiquidity(ctx, pool, externalAssetCoin, nativeAssetCoin, lp, lpUnitsLeft, poolOriginalEB, poolOriginalNB)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToRemoveLiquidity, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSwap(ctx sdk.Context, keeper Keeper, msg MsgSwap) (*sdk.Result, error) {
	var (
		priceImpact sdk.Uint
	)
	logger := keeper.Logger(ctx)

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
	err := errors.New("Swap Error")
	// If sending rowan ,deduct directly from the Native balance  instead of fetching from rowan pool
	if msg.SentAsset != types.GetSettlementAsset() {
		inPool, err = keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	}
	logger.Error("error:", err)
	sentAmountInt, ok := keeper.ParseToInt(sentAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	sentCoin := sdk.NewCoin(msg.SentAsset.Symbol, sentAmountInt)
	err = keeper.InitiateSwap(ctx, sentCoin, msg.Signer)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	// Check if its a two way swap, swapping non native fro non native .
	// If its one way we can skip this if condition and add balance to users account from outpool

	if msg.SentAsset != nativeAsset && msg.ReceivedAsset != nativeAsset {
		emitAmount, lp, ts, finalPool, err := clpkeeper.SwapOne(sentAsset, sentAmount, nativeAsset, inPool)
		if err != nil {
			return nil, err
		}
		err = keeper.SetPool(ctx, finalPool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
		sentAmount = emitAmount
		sentAsset = nativeAsset
		priceImpact = priceImpact.Add(ts)
		liquidityFeeNative = liquidityFeeNative.Add(lp)
	}
	// If receiving  rowan , add directly to  Native balance  instead of fetching from rowan pool
	if msg.ReceivedAsset == types.GetSettlementAsset() {
		outPool, err = keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	} else {
		outPool, err = keeper.GetPool(ctx, msg.ReceivedAsset.Symbol)
		if err != nil {
			return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
		}
	}

	// Calculating amount user receives
	emitAmount, lp, ts, finalPool, err := clpkeeper.SwapOne(sentAsset, sentAmount, receivedAsset, outPool)
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
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
			),
		})
		return &sdk.Result{Events: ctx.EventManager().Events()}, types.ErrReceivedAmountBelowExpected
	}

	err = keeper.FinalizeSwap(ctx, emitAmount.String(), finalPool, msg)

	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	if liquidityFeeNative.GT(sdk.ZeroUint()) {
		liquidityFeeExternal = liquidityFeeExternal.Add(lp)
		firstSwapFeeInOutputAsset := clpkeeper.GetSwapFee(liquidityFeeNative, msg.ReceivedAsset, outPool)
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
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
