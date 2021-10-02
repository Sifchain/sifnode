package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	"github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
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
	registryEntry := srv.tokenRegistryKeeper.GetDenom(registry, msg.Token.Denom)
	if registryEntry == nil {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "denom is not whitelisted")
	}
	// disallow direct transfers of denom aliases
	if registryEntry.UnitDenom != "" && registryEntry.UnitDenom != registryEntry.Denom {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "transfers of denom aliases are not yet supported")
	}
	// check export permission
	if !srv.tokenRegistryKeeper.CheckDenomPermissions(registryEntry, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCEXPORT}) {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "denom cannot be exported")
	}
	// check if registry entry has an IBC counterparty conversion to process
	if registryEntry.IbcCounterpartyDenom != "" && registryEntry.IbcCounterpartyDenom != registryEntry.Denom {
		sendAsRegistryEntry := srv.tokenRegistryKeeper.GetDenom(registry, registryEntry.IbcCounterpartyDenom)
		if sendAsRegistryEntry != nil && sendAsRegistryEntry.Decimals > 0 && registryEntry.Decimals > sendAsRegistryEntry.Decimals {
			token, tokenConversion := helpers.ConvertCoinsForTransfer(msg, registryEntry, sendAsRegistryEntry)
			if token.Amount.Equal(sdk.NewInt(0)) || tokenConversion.Amount.Equal(sdk.NewInt(0)) {
				return nil, types.ErrAmountTooLowToConvert
			}
			if !token.Amount.IsUint64() || !tokenConversion.Amount.IsUint64() {
				return nil, types.ErrAmountTooLargeToSend
			}
			err := helpers.PrepareToSendConvertedCoins(goCtx, msg, token, tokenConversion, srv.bankKeeper)
			if err != nil {
				return nil, sdkerrors.Wrap(types.ErrConvertingToCounterpartyDenom, err.Error())
			}
			msg.Token = tokenConversion
		}
	}
	if !msg.Token.Amount.IsUint64() || msg.Token.Amount.Equal(sdk.NewInt(0)) {
		return nil, types.ErrAmountTooLargeToSend
	}
	return srv.sdkMsgServer.Transfer(goCtx, msg)
}
