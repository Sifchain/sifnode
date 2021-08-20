package keeper

import (
	"context"
	"fmt"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

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
		return nil, fmt.Errorf("Token cannot be exported : %s  , does not have permission", msg.Token.Denom)
	}

	return srv.sdkMsgServer.Transfer(goCtx, msg)
}
