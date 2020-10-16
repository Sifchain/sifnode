package sifnode

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/sifnode/types"
	"github.com/Sifchain/sifnode/x/sifnode/keeper"
)

// Handle a message to delete name
func handleMsgDeleteUser(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteUser) (*sdk.Result, error) {
	if !k.UserExists(ctx, msg.ID) {
		// replace with ErrKeyNotFound for 0.39+
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msg.ID)
	}
	if !msg.Creator.Equals(k.GetUserOwner(ctx, msg.ID)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	k.DeleteUser(ctx, msg.ID)
	return &sdk.Result{}, nil
}
