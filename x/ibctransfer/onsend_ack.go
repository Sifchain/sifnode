package ibctransfer

import (
	"fmt"

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
		// if needs conversion, convert and send
		incomingCoins, finalCoins := GetConvForRefundCoins(ctx, whitelistKeeper, packet, data)
		if incomingCoins != nil && finalCoins != nil {
			err := ExecConvForRefundCoins(ctx, incomingCoins, finalCoins, bankKeeper, packet, data)
			if err != nil {
				return nil, err
			}
			return &sdk.Result{
				Events: ctx.EventManager().Events().ToABCIEvents(),
			}, nil
		}
	}
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

// GetConvForRefundCoins returns 1) the coins that are being received via IBC,
// which need to be deducted from that denom when converting to final denom,
// and 2) the coins that need to be added to the final denom.
func GetConvForRefundCoins(
	ctx sdk.Context,
	whitelistKeeper tokenregistrytypes.Keeper,
	packet channeltypes.Packet,
	data sdktransfertypes.FungibleTokenPacketData,
) (*sdk.Coin, *sdk.Coin) {
	// we don't need to manipulate the denom because the data and packet was created on this chain
	denom := data.Denom
	wl := whitelistKeeper.GetDenomWhitelist(ctx)
	// get token registry entry for received denom
	denomEntry := whitelistKeeper.GetDenom(wl, denom)
	// convert to unit_denom
	if denomEntry == nil || (denomEntry.Decimals == 0 || denomEntry.UnitDenom == "") {
		// noop, should prevent getting here.
		return nil, nil
	}
	convertToDenomEntry := whitelistKeeper.GetDenom(wl, denomEntry.UnitDenom)
	if convertToDenomEntry == nil || convertToDenomEntry.Decimals <= denomEntry.Decimals {
		return nil, nil
	}
	// get the token amount from the packet data
	decAmount := sdk.NewDecFromInt(sdk.NewIntFromUint64(data.Amount))
	// Calculate the conversion difference for increasing precision.
	po := convertToDenomEntry.Decimals - denomEntry.Decimals
	convAmountDec := IncreasePrecision(decAmount, po)
	convAmount := sdk.NewIntFromBigInt(convAmountDec.TruncateInt().BigInt())
	// create converted and ibc tokens with corresponding denoms and amounts
	convertToCoins := sdk.NewCoin(convertToDenomEntry.Denom, convAmount)
	mintedCoins := sdk.NewCoin(denom, sdk.NewIntFromUint64(data.Amount))
	return &mintedCoins, &convertToCoins
}

func ExecConvForRefundCoins(
	ctx sdk.Context,
	incomingCoins *sdk.Coin,
	finalCoins *sdk.Coin,
	bankKeeper sdktransfertypes.BankKeeper,
	packet channeltypes.Packet,
	data sdktransfertypes.FungibleTokenPacketData,
) error {
	// decode the receiver address
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}
	// send ibcdenom coins from account to module
	err = bankKeeper.SendCoinsFromAccountToModule(ctx, sender, sctransfertypes.ModuleName, sdk.NewCoins(*incomingCoins))
	if err != nil {
		return err
	}
	// unescrow original tokens
	escrowAddress := sctransfertypes.GetEscrowAddress(packet.GetSourcePort(), packet.GetSourceChannel())
	if err := bankKeeper.SendCoins(ctx, escrowAddress, sender, sdk.NewCoins(*finalCoins)); err != nil {
		// NOTE: this error is only expected to occur given an unexpected bug or a malicious
		// counterparty module. The bug may occur in bank or any part of the code that allows
		// the escrow address to be drained. A malicious counterparty module could drain the
		// escrow address by allowing more tokens to be sent back then were escrowed.
		return sdkerrors.Wrap(err, "unable to unescrow original tokens")
	}
	// burn ibcdenom coins
	err = bankKeeper.BurnCoins(ctx, sctransfertypes.ModuleName, sdk.NewCoins(*incomingCoins))
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sctransfertypes.EventTypeConvertRefund,
			sdk.NewAttribute(sdk.AttributeKeyModule, sctransfertypes.ModuleName),
			sdk.NewAttribute(sctransfertypes.AttributeKeyPacketAmount, fmt.Sprintf("%v", incomingCoins.Amount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeyPacketDenom, incomingCoins.Denom),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertAmount, fmt.Sprintf("%v", finalCoins.Amount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertDenom, finalCoins.Denom),
		),
	)
	return nil
}
