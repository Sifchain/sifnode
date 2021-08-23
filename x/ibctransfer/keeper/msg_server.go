package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"

	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
)

type msgServer struct {
	bankKeeper          bankkeeper.Keeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	sdkMsgServer        types.MsgServer
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(sdkMsgServer types.MsgServer, bankKeeper bankkeeper.Keeper, tokenRegistryKeeper tokenregistrytypes.Keeper) types.MsgServer {
	return &msgServer{
		sdkMsgServer:        sdkMsgServer,
		bankKeeper:          bankKeeper,
		tokenRegistryKeeper: tokenRegistryKeeper,
	}
}

var _ types.MsgServer = msgServer{}

// Transfer defines a rpc handler method for MsgTransfer.
func (srv msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !srv.tokenRegistryKeeper.CheckDenomPermissions(ctx, msg.Token.Denom, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCEXPORT}) {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "denom cannot be exported")
	}

	return srv.sdkMsgServer.Transfer(goCtx, msg)
}
