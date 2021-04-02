package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the ethbridge MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Lock(goCtx context.Context, msg *types.MsgLockRequest) (*types.MsgLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgLockResponse{}, nil

}
func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurnRequest) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgBurnResponse{}, nil

}
func (k msgServer) CreateEthBridgeClaim(goCtx context.Context, msg *types.MsgCreateEthBridgeClaimRequest) (*types.MsgCreateEthBridgeClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgCreateEthBridgeClaimResponse{}, nil

}
func (k msgServer) UpdateWhiteListValidator(goCtx context.Context, msg *types.MsgUpdateWhiteListValidatorRequest) (*types.MsgUpdateWhiteListValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdateWhiteListValidatorResponse{}, nil

}