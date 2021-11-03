package ibctransfer

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"

	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

// OnAcknowledgementMaybeConvert runs on acknowledgement from receiving chain, of an outgoing send.
// I.e: Receiving ack from recipient chain of our send.
// Case A: Recipient acknowledges our send, do nothing.
// Case B: Recipient acknowledges error of our send, refund sender their coins.
func OnAcknowledgementMaybeConvert(
	ctx sdk.Context,
	sdkTransferKeeper sctransfertypes.SDKTransferKeeper,
	whitelistKeeper tokenregistrytypes.Keeper,
	bankKeeper sdktransfertypes.BankKeeper,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := sdktransfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return  sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	var data sdktransfertypes.FungibleTokenPacketData
	if err := sdktransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	// OnAcknowledgementPacket responds to the the success or failure of a packet
	// acknowledgement written on the receiving chain. If the acknowledgement
	// was a success then nothing occurs. If the acknowledgement failed, then
	// the sender is refunded their tokens using the refundPacketToken function.
	if err := sdkTransferKeeper.OnAcknowledgementPacket(ctx, packet, data, ack); err != nil {
		return err
	}
	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdktransfertypes.EventTypePacket,
				sdk.NewAttribute(sdk.AttributeKeyModule, sdktransfertypes.ModuleName),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyReceiver, data.Receiver),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyDenom, data.Denom),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAck, fmt.Sprintf("%v", ack)),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	// if acknowledgement error then a refund was processed so we must check if conversion is necessary
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdktransfertypes.EventTypePacket,
				sdk.NewAttribute(sdk.AttributeKeyModule, sdktransfertypes.ModuleName),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyReceiver, data.Receiver),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyDenom, data.Denom),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAck, fmt.Sprintf("%v", ack)),
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAckError, resp.Error),
			),
		)
		registry := whitelistKeeper.GetRegistry(ctx)
		denomEntry, err := whitelistKeeper.GetEntry(registry, data.Denom)
		if err == nil && denomEntry.Decimals > 0 && denomEntry.UnitDenom != "" {
			convertToDenomEntry, err := whitelistKeeper.GetEntry(registry, denomEntry.UnitDenom)
			if err == nil && convertToDenomEntry.Decimals > denomEntry.Decimals {
				err := helpers.ExecConvForRefundCoins(ctx, bankKeeper, denomEntry, convertToDenomEntry, packet, data)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}
