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
	if err := trasfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	if !isWhitelisted(k, ctx, packet, data) {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom not on whitelist")
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

func isWhitelisted(k Keeper, ctx sdk.Context, packet channeltypes.Packet, data trasfertypes.FungibleTokenPacketData) bool {
	denomTrace := trasfertypes.ParseDenomTrace(data.Denom)
	return checkWhiteListMap(denomTrace.Hash().String())
}

func checkWhiteListMap(checkingHash string) bool {
	whitelist := make(map[string]bool)
	whitelist["E0263CEED41F926DCE9A805F0358074873E478B515A94DF202E6B69E29DA6178"] = true
	return whitelist[checkingHash]
}
