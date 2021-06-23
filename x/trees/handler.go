package trees

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/trees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateTree:
			return handleMsgCreateTree(ctx, keeper, msg)
		case types.MsgBuyTree:
			return handleMsgBuyTree(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateTree(ctx sdk.Context, keeper Keeper, msg types.MsgCreateTree) (*sdk.Result, error) {
	ctx.Logger().Error("Inside handler msgtree")
	id, err := keeper.CreateTree(ctx, msg)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return nil, errors.Wrapf(types.ErrInvalid, "unable to Create Tree ")
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateTree,
			sdk.NewAttribute(types.AttributeKeyTree, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Seller.String()),
			sdk.NewAttribute(types.AttributeKeyTreeID, id),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBuyTree(ctx sdk.Context, keeper Keeper, msg types.MsgBuyTree) (*sdk.Result, error) {
	ctx.Logger().Error("Inside handler msgtree")
	// id, err := keeper.CreateTree(ctx, msg)
	// if err != nil {
	// 	ctx.Logger().Error(err.Error())
	// 	return nil, errors.Wrapf(types.ErrInvalid, "unable to Create Tree ")
	// }
	id, err := keeper.CreateLimitOrder(ctx, msg)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalid, "unable to Create Limit Order ")
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBuyTree,
			sdk.NewAttribute(types.AttributeKeyTree, types.ModuleName),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, id),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
