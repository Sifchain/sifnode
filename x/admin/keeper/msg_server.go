package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/admin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	keeper Keeper
}

func (m msgServer) SetParams(ctx context.Context, msg *types.MsgSetParams) (*types.MsgSetParamsResponse, error) {
	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), types.AdminType_ADMIN, addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}

	m.keeper.SetParams(sdk.UnwrapSDKContext(ctx), msg.Params)
	return &types.MsgSetParamsResponse{}, nil
}

func (m msgServer) AddAccount(ctx context.Context, msg *types.MsgAddAccount) (*types.MsgAddAccountResponse, error) {
	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), types.AdminType_ADMIN, addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}

	m.keeper.SetAdminAccount(sdk.UnwrapSDKContext(ctx), msg.Account)

	return &types.MsgAddAccountResponse{}, nil
}

func (m msgServer) RemoveAccount(ctx context.Context, msg *types.MsgRemoveAccount) (*types.MsgRemoveAccountResponse, error) {
	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	if !m.keeper.IsAdminAccount(sdk.UnwrapSDKContext(ctx), types.AdminType_ADMIN, addr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unauthorised signer")
	}

	m.keeper.RemoveAdminAccount(sdk.UnwrapSDKContext(ctx), msg.Account)

	return &types.MsgRemoveAccountResponse{}, nil
}

// NewMsgServerImpl returns an implementation of MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}
