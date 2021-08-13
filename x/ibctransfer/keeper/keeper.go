package keeper

import (
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktransferkeeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
)

// Keeper defines the IBC fungible transfer keeper
type Keeper struct {
	storeKey            sdk.StoreKey
	sdktransferkeeper   sdktransferkeeper.Keeper
	tokenregistrykeeper tokenregistrytypes.Keeper
	bankKeeper          types.BankKeeper
}

// var _ types.MsgServer = Keeper{}

// // MsgServer is the server API for Msg service.
// type MsgServer interface {
// 	// Transfer defines a rpc handler method for MsgTransfer.
// 	Transfer(sdk.Context, *types.MsgTransfer) (*types.MsgTransferResponse, error)
// }

// type msgServer struct {
// 	Keeper
// }

// // NewMsgServerImpl returns an implementation of the bank MsgServer interface
// // for the provided Keeper.
// func NewMsgServerImpl(keeper Keeper) types.MsgServer {
// 	return &msgServer{Keeper: keeper}
// }

// var _ types.MsgServer = msgServer{}

// Transfer defines a rpc handler method for MsgTransfer.
func (k Keeper) WrappedTransfer(ctx sdk.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	// get token registry entry for sent token
	registryEntry := k.tokenregistrykeeper.GetDenom(ctx, msg.Token.Denom)
	// check if registry entry has an IBC decimal field
	if registryEntry.IBCDecimals != nil && registryEntry.Decimals > 10 {
		//
	}

	// k.bankKeeper.

	return k.sdktransferkeeper.Transfer(ctx, msg)
}
