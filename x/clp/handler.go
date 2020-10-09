package clp

import (
	"fmt"
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
	asset := NewAsset("ROWAN", "ETHROWAN", "ETH")
	pool := NewPool(asset, 1000, 200, 8, "ethRowan")
	keeper.SetPool(ctx, pool)
	return &sdk.Result{}, nil
}

func handleMsgAddLiquidity(ctx sdk.Context, keeper Keeper, msg MsgAddLiquidity) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func handleMsgRemoveLiquidity(ctx sdk.Context, keeper Keeper, msg MsgRemoveLiquidity) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func handleMsgSwap(ctx sdk.Context, keeper Keeper, msg MsgSwap) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}
