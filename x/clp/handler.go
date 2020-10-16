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

	pool, err := keeper.GetPool(ctx, msg.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	if pool.ExternalAssetBalance+pool.NativeAssetBalance > keeper.GetParams(ctx).MinCreatePoolThreshold {
		return nil, types.ErrBalanceTooHigh
	}
	lpList := keeper.GetLiqudityProvidersForAsset(ctx, pool.ExternalAsset)
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	for _, lp := range lpList {
		withdrawNativeAsset, withdrawExternalAsset, _ := calculateWithdrawl(poolUnits, nativeAssetBalance, externalAssetBalance, lp.LiquidityProviderUnits, 10000, 1)
		poolUnits = poolUnits - lp.LiquidityProviderUnits
		nativeAssetBalance = nativeAssetBalance - withdrawNativeAsset
		externalAssetBalance = externalAssetBalance - withdrawExternalAsset
		//send withdrawNativeAsset to liquidityProvider.lpAddress
		//send withdrawExternalAsset to liquidityProvider.lpAddress
		keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
	}
	err = keeper.DestroyPool(ctx, pool.ExternalAsset.Ticker)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToDestroyPool, err.Error())
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
	MinThreshold := keeper.GetParams(ctx).MinCreatePoolThreshold
	if (msg.ExternalAssetAmount + msg.NativeAssetAmount) < MinThreshold { // Need to verify
		return nil, types.ErrTotalAmountTooLow
	}
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker)
	if err == nil {
		return nil, types.ErrUnableToCreatePool
	}
	asset := msg.ExternalAsset
	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	poolUnits, lpunits := calculatePoolUnits(0, 0, 0, nativeBalance, externalBalance)
	pool, err = NewPool(asset, nativeBalance, externalBalance, poolUnits)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}
	lp := NewLiquidityProvider(asset, lpunits, msg.Signer.String())
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	keeper.SetLiquidityProvider(ctx, lp)
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
	createNewLP := false
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	newPoolUnits, lpUnits := calculatePoolUnits(pool.PoolUnits, pool.NativeAssetBalance, pool.ExternalAssetBalance, msg.NativeAssetAmount, msg.ExternalAssetAmount)
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		createNewLP = true
	}

	pool.PoolUnits = newPoolUnits
	pool.NativeAssetBalance = pool.NativeAssetBalance + msg.NativeAssetAmount
	pool.ExternalAssetBalance = pool.ExternalAssetBalance + msg.ExternalAssetAmount
	if createNewLP {
		lp := NewLiquidityProvider(msg.ExternalAsset, lpUnits, msg.Signer.String())
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeCreateLiquidityProvider,
				sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
		})
	} else {
		lp.LiquidityProviderUnits = lp.LiquidityProviderUnits + lpUnits
	}
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	keeper.SetLiquidityProvider(ctx, lp)
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
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft := calculateWithdrawl(pool.PoolUnits,
		pool.NativeAssetBalance, pool.ExternalAssetBalance, lp.LiquidityProviderUnits,
		msg.WBasisPoints, msg.Asymmetry)
	pool.PoolUnits = pool.PoolUnits - lp.LiquidityProviderUnits + lpUnitsLeft
	pool.NativeAssetBalance = pool.NativeAssetBalance - withdrawNativeAssetAmount
	pool.ExternalAssetBalance = pool.ExternalAssetBalance - withdrawExternalAssetAmount
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	if lpUnitsLeft == 0 {
		keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
	} else {
		lp.LiquidityProviderUnits = lpUnitsLeft
		keeper.SetLiquidityProvider(ctx, lp)
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
		liquidityFee uint
		tradeSlip    uint
	)
	sentAmount := msg.SentAmount
	sentAsset := msg.SentAsset
	receivedAsset := msg.ReceivedAsset
	nativeAsset := types.GetSettlementAsset()
	if msg.SentAsset != nativeAsset && msg.ReceivedAsset != nativeAsset {
		inPool, err := keeper.GetPool(ctx, msg.SentAsset.Ticker)
		if err != nil {
			return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
		emitAmount, lp, ts, err := swapOne(ctx, keeper, sentAsset, sentAmount, nativeAsset, inPool)
		if err != nil {
			return nil, err
		}
		sentAmount = emitAmount
		sentAsset = nativeAsset
		liquidityFee = liquidityFee + lp
		tradeSlip = tradeSlip + ts
	}

	outPool, err := keeper.GetPool(ctx, msg.ReceivedAsset.Ticker)
	if err != nil {
		return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
	}
	emitAmount, lp, ts, err := swapOne(ctx, keeper, sentAsset, sentAmount, receivedAsset, outPool)
	if err != nil {
		return nil, err
	}
	liquidityFee = liquidityFee + lp
	tradeSlip = tradeSlip + ts
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeySwapAmount, strconv.FormatInt(int64(emitAmount), 10)),
			sdk.NewAttribute(types.AttributeKeyLiquidityFee, strconv.FormatInt(int64(liquidityFee), 10)),
			sdk.NewAttribute(types.AttributeKeyTradeSlip, strconv.FormatInt(int64(tradeSlip), 10)),
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

//------------------------------------------------------------------------------------------------------------------

func swapOne(ctx sdk.Context, keeper Keeper, from Asset, sentAmount uint, to Asset, pool Pool) (uint, uint, uint, error) {

	var X uint
	var Y uint
	if to == GetNativeAsset() {
		Y = pool.NativeAssetBalance
		X = pool.ExternalAssetBalance
	} else {
		X = pool.NativeAssetBalance
		Y = pool.ExternalAssetBalance
	}
	x := sentAmount
	liquidityFee := calcLiquidityFee(X, x, Y)
	tradeSlip := calcTradeSlip(X, x)
	swapResult := calcSwapResult(X, x, Y)
	if swapResult >= Y {
		return 0, 0, 0, types.ErrNotEnoughAssetTokens
	}
	if from == GetNativeAsset() {
		pool.NativeAssetBalance = X + x
		pool.ExternalAssetBalance = Y - swapResult
	} else {
		pool.ExternalAssetBalance = X + x
		pool.NativeAssetBalance = Y - swapResult
	}
	err := keeper.SetPool(ctx, pool)
	if err != nil {
		return 0, 0, 0, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	return swapResult, liquidityFee, tradeSlip, nil
}

func calculateWithdrawl(poolUnits uint, nativeAssetBalance uint,
	externalAssetBalance uint, lpUnits uint, wBasisPoints int, asymmetry int) (uint, uint, uint) {
	var (
		nativeAssetUnits            int
		withdrawNativeAssetAmount   int
		externalAssetUnits          int
		withdrawExternalAssetAmount int
	)

	unitsToClaim := int(lpUnits) / (10000 / wBasisPoints)
	if asymmetry == 0 {
		externalAssetUnits = unitsToClaim
		nativeAssetUnits = unitsToClaim
	} else {
		externalAssetUnits = unitsToClaim + (unitsToClaim / (10000 / asymmetry))
		nativeAssetUnits = unitsToClaim - (unitsToClaim / (10000 / asymmetry))
	}
	if externalAssetUnits == 0 {
		withdrawExternalAssetAmount = 0
	} else {
		withdrawExternalAssetAmount = int(externalAssetBalance) / (int(poolUnits) / externalAssetUnits)
	}
	if nativeAssetUnits == 0 {
		withdrawNativeAssetAmount = 0
	} else {
		withdrawNativeAssetAmount = int(nativeAssetBalance) / (int(poolUnits) / nativeAssetUnits)
	}
	lpUnitsLeft := int(lpUnits) - unitsToClaim
	if withdrawNativeAssetAmount < 0 {
		withdrawNativeAssetAmount = 0
	}
	if withdrawExternalAssetAmount < 0 {
		withdrawExternalAssetAmount = 0
	}
	if lpUnitsLeft < 0 {
		lpUnitsLeft = 0
	}
	return uint(withdrawNativeAssetAmount), uint(withdrawExternalAssetAmount), uint(lpUnitsLeft)
}

func calculatePoolUnits(oldPoolUnits uint, nativeAssetBalance uint, externalAssetBalance uint,
	nativeAssetAmount uint, externalAssetAmount uint) (uint, uint) {
	R := nativeAssetBalance + nativeAssetAmount
	A := externalAssetBalance + externalAssetAmount
	r := nativeAssetAmount
	a := externalAssetAmount
	lpUnits := ((R + A) * (r*A + R*a)) / (4 * R * A)
	poolUnits := oldPoolUnits + lpUnits
	return poolUnits, lpUnits
}

func calcLiquidityFee(X, x, Y uint) uint {
	return (x * x * Y) / ((x + X) * (x + X))
}

func calcTradeSlip(X, x uint) uint {
	return x * (2*X + x) / (X * X)
}

func calcSwapResult(X, x, Y uint) uint {
	return (x * X * Y) / ((x + X) * (x + X))
}
