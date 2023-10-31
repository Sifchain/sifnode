package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdktransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
)

type msgServer struct {
	bankKeeper          types.BankKeeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	sdkMsgServer        types.MsgServer
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(sdkMsgServer types.MsgServer, bankKeeper types.BankKeeper, tokenRegistryKeeper tokenregistrytypes.Keeper) sdktransfertypes.MsgServer {
	return &msgServer{
		sdkMsgServer:        sdkMsgServer,
		bankKeeper:          bankKeeper,
		tokenRegistryKeeper: tokenRegistryKeeper,
	}
}

var _ sdktransfertypes.MsgServer = msgServer{}

// Transfer defines a rpc handler method for MsgTransfer.
func (srv msgServer) Transfer(goCtx context.Context, msg *sdktransfertypes.MsgTransfer) (*sdktransfertypes.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	registry := srv.tokenRegistryKeeper.GetRegistry(ctx)
	registryEntry, err := srv.tokenRegistryKeeper.GetEntry(registry, msg.Token.Denom)
	if err != nil {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "denom is not whitelisted")
	}
	// disallow direct transfers of denom aliases
	if registryEntry.UnitDenom != "" && registryEntry.UnitDenom != registryEntry.Denom {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "transfers of denom aliases are not yet supported")
	}
	// check export permission
	if !srv.tokenRegistryKeeper.CheckEntryPermissions(registryEntry, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCEXPORT}) {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "denom cannot be exported")
	}
	if msg.Token.Amount.LTE(sdk.NewInt(0)) {
		return nil, types.ErrAmountTooLowToConvert
	}

	return srv.sdkMsgServer.Transfer(goCtx, msg)
}
