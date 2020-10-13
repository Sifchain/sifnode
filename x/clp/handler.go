package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		),
		sdk.NewEvent(
			types.EventTypeCreateLiquidityProvider,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
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

func handleMsgAddLiquidity(ctx sdk.Context, keeper Keeper, msg MsgAddLiquidity) (*sdk.Result, error) {
	createNewLP := false
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker, msg.ExternalAsset.SourceChain)
	if err != nil {
		return nil, err
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
	return &sdk.Result{}, nil
}

func handleMsgSwap(ctx sdk.Context, keeper Keeper, msg MsgSwap) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}
