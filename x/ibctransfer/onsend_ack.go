package ibctransfer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

// OnAcknowledgementMaybeConvert() runs on acknowledgement from receiving chain, of an outgoing send.
// I.e: Receiving ack from recipient chain of our send.
// Case A: Recipient acknowledges our send, do nothing.
// Case B: Recipient acknowledges error of our send, refund sender their coins.
func OnAcknowledgementMaybeConvert(
	ctx sdk.Context,
	sdkTransferKeeper sctransfertypes.SDKTransferKeeper,
	whitelistKeeper tokenregistrytypes.Keeper,
	bankKeeper transfertypes.BankKeeper,
	packet channeltypes.Packet,
	acknowledgement []byte,
) (*sdk.Result, error) {

	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
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
			transfertypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
			sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
			sdk.NewAttribute(transfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttribute(transfertypes.AttributeKeyAck, fmt.Sprintf("%v", ack)),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				transfertypes.EventTypePacket,
				sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	// if acknowledgement error then a refund was processed so we must check if conversion is necessary
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				transfertypes.EventTypePacket,
				sdk.NewAttribute(transfertypes.AttributeKeyAckError, resp.Error),
			),
		)
		// TODO: Always refund, not only if sender of ack is source.
		// TODO: Copy error / panic pattern from sdkkeeper.refundPacketToken in ExecConvForIncomingCoins.
		// TODO: Why does sdk transfer module use escrow vs minting in different scenarios here?
		// if sender is source check for conversion
		if transfertypes.SenderChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
			// if needs conversion, convert and send
			if ShouldConvertIncomingCoins(ctx, whitelistKeeper, packet, data) {
				ibcToken, convToken := GetConvForIncomingCoins(ctx, whitelistKeeper, packet, data)
				err := ExecConvForIncomingCoins(ctx, ibcToken, convToken, bankKeeper, packet, data)
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
