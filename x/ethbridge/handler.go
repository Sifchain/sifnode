//nolint:dupl
package ethbridge

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(k Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgCreateEthBridgeClaim:
			res, err := msgServer.CreateEthBridgeClaim(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBurn:
			res, err := msgServer.Burn(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgLock:
			res, err := msgServer.Lock(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateWhiteListValidator:
			res, err := msgServer.UpdateWhiteListValidator(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateCrossChainFeeReceiverAccount:
			res, err := msgServer.UpdateCrossChainFeeReceiverAccount(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRescueCrossChainFee:
			res, err := msgServer.RescueCrossChainFee(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSetFeeInfo:
			res, err := msgServer.SetFeeInfo(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized ethbridge message type: %v", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
