package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"strconv"
)

// NewHandler creates an sdk.Handler for all the clp type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreatePool:
			return handleMsgCreatePool(ctx, k, msg)
		case MsgDecommissionPool:
			return handleMsgDecommissionPool(ctx, k, msg)
		case MsgAddLiquidity:
			return handleMsgAddLiquidity(ctx, k, msg)
		case MsgRemoveLiquidity:
			return handleMsgRemoveLiquidity(ctx, k, msg)
		case MsgSwap:
			return handleMsgSwap(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDecommissionPool(ctx sdk.Context, keeper Keeper, msg MsgDecommissionPool) (*sdk.Result, error) {
	// Verify pool
	pool, err := keeper.GetPool(ctx, msg.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	if !keeper.ValidateAddress(ctx, msg.Signer) {
		return nil, errors.Wrap(types.ErrInvalid, "user does not have permission to decommission pool")
	}
	if pool.NativeAssetBalance.GTE(sdk.NewUint(uint64(keeper.GetParams(ctx).MinCreatePoolThreshold))) {
		return nil, types.ErrBalanceTooHigh
	}
	// Get all LP's for the pool
	lpList := keeper.GetLiqudityProvidersForAsset(ctx, pool.ExternalAsset)
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	// iterate over Lp list and refund them there tokens
	// Return both RWN and EXTERNAL ASSET
	for _, lp := range lpList {
		withdrawNativeAsset, withdrawExternalAsset, _, _ := CalculateWithdrawal(poolUnits, nativeAssetBalance.String(), externalAssetBalance.String(),
			lp.LiquidityProviderUnits.String(), sdk.NewInt(MaxWbasis).String(), sdk.ZeroInt())
		poolUnits = poolUnits.Sub(lp.LiquidityProviderUnits)
		nativeAssetBalance = nativeAssetBalance.Sub(withdrawNativeAsset)
		externalAssetBalance = externalAssetBalance.Sub(withdrawExternalAsset)
		withdrawNativeCoins := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(withdrawNativeAsset.Uint64()))
		withdrawExternalCoins := sdk.NewCoin(msg.Ticker, sdk.NewIntFromUint64(withdrawExternalAsset.Uint64()))
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
	MinThreshold := sdk.NewUint(uint64(keeper.GetParams(ctx).MinCreatePoolThreshold))
	if msg.NativeAssetAmount.LT(MinThreshold) { // Need to verify
		return nil, types.ErrTotalAmountTooLow
	}
	// Check if pool already exists
	if keeper.ExistsPool(ctx, msg.ExternalAsset.Ticker) {
		return nil, types.ErrUnableToCreatePool
	}

	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	poolUnits, lpunits, err := calculatePoolUnits(sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), nativeBalance, externalBalance)
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
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	newPoolUnits, lpUnits, err := calculatePoolUnits(pool.PoolUnits, pool.NativeAssetBalance, pool.ExternalAssetBalance, msg.NativeAssetAmount, msg.ExternalAssetAmount)
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
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}

	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(withdrawExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(withdrawNativeAssetAmount.Uint64()))

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
		swapResult, _, _, swappedPool, err := SwapOne(GetSettlementAsset(), swapAmount, msg.ExternalAsset, pool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(swapResult.Uint64()))
			swapAmountInCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(swapAmount.Uint64()))
			externalAssetCoin = externalAssetCoin.Add(swapCoin)
			nativeAssetCoin = nativeAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	if msg.Asymmetry.IsNegative() {
		swapResult, _, _, swappedPool, err := SwapOne(msg.ExternalAsset, swapAmount, GetSettlementAsset(), pool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(swapResult.Uint64()))
			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(swapAmount.Uint64()))

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
		liquidityFee sdk.Uint
		tradeSlip    sdk.Uint
	)
	liquidityFee = sdk.ZeroUint()
	tradeSlip = sdk.ZeroUint()
	sentAmount := msg.SentAmount
	sentAsset := msg.SentAsset
	receivedAsset := msg.ReceivedAsset
	// Get native asset
	nativeAsset := types.GetSettlementAsset()
	// If its one swap , this pool would be RWN:RWN ( Ex User sends RWN wants ETH)
	// If its two swap . this pool would be RWN:EXTERNAL1 ( Ex User sends ETH wants XCT , ETH is EXTERNAL1)
	//CASE 1 : RWN:ETH
	//CASE 2 : RWN:ETH
	inPool, err := keeper.GetPool(ctx, msg.SentAsset.Ticker)
	if err != nil {
		return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
	}
	// If its one swap , this pool would be RWN:EXTERNAL ( Ex User sends RWN wants ETH , ETH IS EXTERNAL )
	// If its two swap . this pool would be RWN:EXTERNAL2 ( Ex User sends ETH wants XCT , XCT is EXTERNAL2)
	//CASE 1 : RWN:ETH
	//CASE 2 : RWN:XCT
	outPool, err := keeper.GetPool(ctx, msg.ReceivedAsset.Ticker)
	if err != nil {
		return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
	}

	// Deducting Balance from the user , Sent Asset is the asset the user is sending to the Pool
	// Case 1 . Deducting his RWN and adding to RWN:ETH pool
	// Case 2 , Deduction his ETH and adding to RWN:ETH pool
	sentCoin := sdk.NewCoin(msg.SentAsset.Ticker, sdk.NewIntFromUint64(sentAmount.Uint64()))
	err = keeper.InitiateSwap(ctx, sentCoin, msg.Signer)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	// Check if its a two way swap, swapping non native fro non native .
	// If its one way we can skip this if condition and add balance to users account from outpool
	if msg.SentAsset != nativeAsset && msg.ReceivedAsset != nativeAsset {

		emitAmount, lp, ts, finalPool, err := SwapOne(sentAsset, sentAmount, nativeAsset, inPool)
		if err != nil {
			return nil, err
		}
		err = keeper.SetPool(ctx, finalPool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
		sentAmount = emitAmount
		sentAsset = nativeAsset
		liquidityFee = liquidityFee.Add(lp)
		tradeSlip = tradeSlip.Add(ts)
	}
	// Calculating amount user receives
	emitAmount, lp, ts, finalPool, err := SwapOne(sentAsset, sentAmount, receivedAsset, outPool)
	if err != nil {
		return nil, err
	}
	err = keeper.FinalizeSwap(ctx, sentAmount, finalPool, msg)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	liquidityFee = liquidityFee.Add(lp)
	tradeSlip = tradeSlip.Add(ts)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
			sdk.NewAttribute(types.AttributeKeyLiquidityFee, liquidityFee.String()),
			sdk.NewAttribute(types.AttributeKeyTradeSlip, tradeSlip.String()),
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
