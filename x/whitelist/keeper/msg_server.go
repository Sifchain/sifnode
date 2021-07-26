package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/whitelist/types"
)

type msgServer struct {
	keeper types.Keeper
}

func (m msgServer) UpdateWhitelist(ctx context.Context, req *types.MsgUpdateWhitelist) (
	*types.MsgUpdateWhitelistResponse, error) {

	addr, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, err
	}
	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}

	m.keeper.SetDenom(sdk.UnwrapSDKContext(ctx), req.From, req.Decimals)

	return &types.MsgUpdateWhitelistResponse{}, nil
}

// NewMsgServerImpl returns an implementation of MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}
