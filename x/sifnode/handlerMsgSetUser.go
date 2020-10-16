package sifnode

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/sifnode/types"
	"github.com/Sifchain/sifnode/x/sifnode/keeper"
)

func handleMsgSetUser(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetUser) (*sdk.Result, error) {
	var user = types.User{
		Creator: msg.Creator,
		ID:      msg.ID,
    	Name: msg.Name,
    	Email: msg.Email,
	}
	if !msg.Creator.Equals(k.GetUserOwner(ctx, msg.ID)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	k.SetUser(ctx, user)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
