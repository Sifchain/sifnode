package ibc_transfer

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	trasfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
)

func OnRecvPacketWhiteListed(
	k Keeper,
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, []byte, error) {
	var data trasfertypes.FungibleTokenPacketData
	if !isWhitelisted(data.Denom) {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom not on whitelist")
	}
	if err := trasfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	err := k.OnRecvPacket(ctx, packet, data)
	if err != nil {
		acknowledgement = channeltypes.NewErrorAcknowledgement(err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			trasfertypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, trasfertypes.ModuleName),
			sdk.NewAttribute(trasfertypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(trasfertypes.AttributeKeyDenom, data.Denom),
			sdk.NewAttribute(trasfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttribute(trasfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", err != nil)),
		),
	)

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, acknowledgement.GetBytes(), nil
}

func isWhitelisted(denom string) bool {
	return true
}
