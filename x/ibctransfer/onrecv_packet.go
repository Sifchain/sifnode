package ibctransfer

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
)

func OnRecvPacketWhitelistConvert(
	ctx sdk.Context,
	sdkTransferKeeper sctransfertypes.SDKTransferKeeper,
	whitelistKeeper tokenregistrytypes.Keeper,
	bankKeeper transfertypes.BankKeeper,
	packet channeltypes.Packet,
) (*sdk.Result, []byte, error) {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	err := sdkTransferKeeper.OnRecvPacket(ctx, packet, data)
	if err != nil {
		acknowledgement := channeltypes.NewErrorAcknowledgement(err.Error())
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, acknowledgement.GetBytes(), nil
	}
	// Get the denom that will be minted by sdk transfer module,
	// so that it can be converted to the denom it should be stored as.
	// For a native token that has been returned, this will just be a base_denom,
	// which will be on the whitelist.
	mintedDenom := helpers.GetMintedDenomFromPacket(packet, data)
	registry := whitelistKeeper.GetRegistry(ctx)
	mintedDenomEntry := whitelistKeeper.GetDenom(registry, mintedDenom)
	if !helpers.IsRecvPacketAllowed(ctx, whitelistKeeper, packet, data, mintedDenomEntry) {
		acknowledgement := channeltypes.NewErrorAcknowledgement(
			sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom not whitelisted").Error(),
		)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				transfertypes.EventTypePacket,
				sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
				sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
				sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
				sdk.NewAttribute(transfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
				sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", false)),
			),
		)
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, acknowledgement.GetBytes(), nil
	}
	convertToDenomEntry := whitelistKeeper.GetDenom(registry, mintedDenomEntry.UnitDenom)
	if convertToDenomEntry != nil && convertToDenomEntry.Decimals > 0 && mintedDenomEntry.Decimals > 0 && convertToDenomEntry.Decimals > mintedDenomEntry.Decimals {
		err = helpers.ExecConvForIncomingCoins(ctx, bankKeeper, whitelistKeeper, mintedDenomEntry, convertToDenomEntry, packet, data)
		// Revert, although this may cause packet to be relayed again.
		if err != nil {
			return nil, nil, sdkerrors.Wrap(sctransfertypes.ErrConvertingToUnitDenom, err.Error())
		}
	}
	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			transfertypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
			sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
			sdk.NewAttribute(transfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", err == nil)),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, acknowledgement.GetBytes(), nil
}
