package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

// NewHandler creates an sdk.Handler for all the clp type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreatePool:
			return handleMsgCreatePool(ctx, k, msg)
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

func handleMsgCreatePool(ctx sdk.Context, keeper Keeper, msg MsgCreatePool) (*sdk.Result, error) {
	MinThreshold := keeper.GetParams(ctx).MinCreatePoolThreshold
	if (msg.ExternalAssetAmount + msg.NativeAssetAmount) < MinThreshold {
		return nil, types.TotalAmountTooLow
	}
	asset := msg.ExternalAsset
	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	poolUnits, lpunits := calculatePoolUnits(0, 0, 0, nativeBalance, externalBalance)
	pool := NewPool(asset, nativeBalance, externalBalance, poolUnits)
	lp := NewLiquidityProvider(asset, lpunits, msg.Signer.String())
	keeper.SetPool(ctx, pool)
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
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker, msg.ExternalAsset.SourceChain)
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
	keeper.SetPool(ctx, pool)
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
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker, msg.ExternalAsset.SourceChain)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		return nil, types.LiquidityProviderDoesNotExist
	}
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft := calculateWithdrawl(pool.PoolUnits,
		pool.NativeAssetBalance, pool.ExternalAssetBalance, lp.LiquidityProviderUnits,
		msg.WBasisPoints, msg.Asymmetry)
	pool.PoolUnits = pool.PoolUnits - lp.LiquidityProviderUnits + lpUnitsLeft
	pool.NativeAssetBalance = pool.NativeAssetBalance - withdrawNativeAssetAmount
	pool.ExternalAssetBalance = pool.ExternalAssetBalance - withdrawExternalAssetAmount
	keeper.SetPool(ctx, pool)

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

	return &sdk.Result{}, nil
}

func handleMsgSwap(ctx sdk.Context, keeper Keeper, msg MsgSwap) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

//------------------------------------------------------------------------------------------------------------------

func calculateWithdrawl(poolUnits uint, nativeAssetBalance uint,
	externalAssetBalance uint, lpUnits uint, wBasisPoints int, asymmetry int) (uint, uint, uint) {
	var (
		externalAssetPercent int
		nativeAssetBasis     uint
		externalAssetBasis   uint
	)
	if asymmetry < 0 {
		asymmetry = asymmetry * -1
		nativeAssetPercent := asymmetry + (asymmetry) ^ 2
		externalAssetPercent := 1 - nativeAssetPercent
		nativeAssetBasis = uint(nativeAssetPercent * wBasisPoints)
		externalAssetBasis = uint(externalAssetPercent * wBasisPoints)
	} else if asymmetry > 0 {
		externalAssetPercent = asymmetry + asymmetry ^ 2
		nativeAssetPercent := 1 - externalAssetPercent
		nativeAssetBasis = uint(nativeAssetPercent * wBasisPoints)
		externalAssetBasis = uint(externalAssetPercent * wBasisPoints)
	} else if asymmetry == 0 {
		externalAssetBasis = uint(wBasisPoints / 2) //ignoring decimals
		nativeAssetBasis = uint(wBasisPoints / 2)   //ignoring decimals
	}
	nativeAssetUnits := lpUnits / (10000 / nativeAssetBasis)
	externalAssetUnits := lpUnits / (10000 / externalAssetBasis)
	withdrawNativeAssetAmount := nativeAssetBalance / (poolUnits / nativeAssetUnits)
	withdrawExternalAssetAmount := externalAssetBalance / (poolUnits / externalAssetUnits)
	lpUnitsLeft := lpUnits - (nativeAssetUnits + externalAssetUnits)
	return withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft
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
