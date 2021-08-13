package keeper

import (
	"context"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	sdktransferkeeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
)

// Keeper defines the IBC fungible transfer keeper
type Keeper struct {
	sdkTransferKeeper   sdktransferkeeper.Keeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	bankKeeper          types.BankKeeper
}

var _ types.MsgServer = msgServer{}

type msgServer struct {
	bankKeeper          bankkeeper.Keeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	sdkMsgServer        types.MsgServer
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl( /* bankKeeper, tokenRegistryKeeper */ ) types.MsgServer {
	return &msgServer{}
}

// var _ types.MsgServer = msgServer{}

// Transfer defines a rpc handler method for MsgTransfer.
func (srv msgServer) Transfer(ctx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	// get token registry entry for sent token
	//registryEntry := srv.tokenRegistryKeeper.GetDenom(sdk.UnwrapSDKContext(ctx), msg.Token.Denom)
	// check if registry entry has an IBC decimal field
	//if registryEntry.IBCDecimals != nil && registryEntry.Decimals > 10 {
	//
	//}

	// k.bankKeeper.

	return srv.sdkMsgServer.Transfer(ctx, msg)
}
