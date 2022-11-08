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

func (m msgServer) Register(ctx context.Context, req *types.MsgRegister) (*types.MsgRegisterResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, err
	}
	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), types.AdminType_TOKENREGISTRY, addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}
	m.keeper.SetToken(sdk.UnwrapSDKContext(ctx), req.Entry)
	return &types.MsgRegisterResponse{}, nil
}

func (m msgServer) SetRegistry(ctx context.Context, req *types.MsgSetRegistry) (*types.MsgSetRegistryResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, err
	}
	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), types.AdminType_TOKENREGISTRY, addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}
	m.keeper.SetRegistry(sdk.UnwrapSDKContext(ctx), *req.Registry)
	return &types.MsgSetRegistryResponse{}, nil
}

func (m msgServer) Deregister(ctx context.Context, req *types.MsgDeregister) (*types.MsgDeregisterResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.From)
	if err != nil {
		return nil, err
	}
	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), types.AdminType_TOKENREGISTRY, addr) {
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
