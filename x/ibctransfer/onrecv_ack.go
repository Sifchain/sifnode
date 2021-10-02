package ibctransfer

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

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
) (*sdk.Result, error) {
	var ack channeltypes.Acknowledgement
	if err := sdktransfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	var data sdktransfertypes.FungibleTokenPacketData
	if err := sdktransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	// OnAcknowledgementPacket responds to the the success or failure of a packet
	// acknowledgement written on the receiving chain. If the acknowledgement
	// was a success then nothing occurs. If the acknowledgement failed, then
	// the sender is refunded their tokens using the refundPacketToken function.
	if err := sdkTransferKeeper.OnAcknowledgementPacket(ctx, packet, data, ack); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdktransfertypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, sdktransfertypes.ModuleName),
			sdk.NewAttribute(sdktransfertypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(sdktransfertypes.AttributeKeyDenom, data.Denom),
			sdk.NewAttribute(sdktransfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttribute(sdktransfertypes.AttributeKeyAck, fmt.Sprintf("%v", ack)),
		),
	)
	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdktransfertypes.EventTypePacket,
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	// if acknowledgement error then a refund was processed so we must check if conversion is necessary
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdktransfertypes.EventTypePacket,
				sdk.NewAttribute(sdktransfertypes.AttributeKeyAckError, resp.Error),
			),
		)
		denom := data.Denom
		registry := whitelistKeeper.GetRegistry(ctx)
		denomEntry := whitelistKeeper.GetEntry(registry, denom)
		if denomEntry != nil && denomEntry.Decimals > 0 && denomEntry.UnitDenom != "" {
			convertToDenomEntry := whitelistKeeper.GetEntry(registry, denomEntry.UnitDenom)
			if convertToDenomEntry != nil && convertToDenomEntry.Decimals > denomEntry.Decimals {
				err := helpers.ExecConvForRefundCoins(ctx, bankKeeper, whitelistKeeper, denomEntry, convertToDenomEntry, packet, data)
				if err != nil {
					return nil, err
				}
				return &sdk.Result{
					Events: ctx.EventManager().Events().ToABCIEvents(),
				}, nil
			}
		}
	}
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
