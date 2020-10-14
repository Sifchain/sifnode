package sifnode

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/sifnode/keeper"
	"github.com/Sifchain/sifnode/x/sifnode/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		// case MsgCreateEthBridgeClaim:
		// 	return handleMsgCreateEthBridgeClaim(ctx, cdc, bridgeKeeper, msg)
		// case MsgBurn:
		// 	return handleMsgBurn(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		// case MsgLock:
		// 	return handleMsgLock(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		// this line is used by starport scaffolding
		default:
			fmt.Print(msg)
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
