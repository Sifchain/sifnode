package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

type msgServer struct {
	keeper types.Keeper
}

func (m msgServer) Register(ctx context.Context, req *types.MsgRegister) (
	*types.MsgRegisterResponse, error) {

	addr, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, err
	}

	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}

	m.keeper.SetToken(sdk.UnwrapSDKContext(ctx), req.Entry)

	return &types.MsgRegisterResponse{}, nil
}

func (m msgServer) Deregister(ctx context.Context, req *types.MsgDeregister) (
	*types.MsgDeregisterResponse, error) {

	addr, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, err
	}

	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}

	m.keeper.RemoveToken(sdk.UnwrapSDKContext(ctx), req.Denom)

	return &types.MsgDeregisterResponse{}, nil

}

// NewMsgServerImpl returns an implementation of MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}
