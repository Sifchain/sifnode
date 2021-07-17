package ibc_sifchain

import (
	"fmt"
	whitelisttypes "github.com/Sifchain/sifnode/x/whitelist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
)

func OnRecvPacketWhiteListed(
	k Keeper,
	ctx sdk.Context,
	packet channeltypes.Packet,
	whitelistKeeper whitelisttypes.Keeper,
) (*sdk.Result, []byte, error) {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	denom, ok := GetMintedDenomFromPacket(packet, data)
	if !(ok && whitelistKeeper.GetDenom(ctx, denom).IsWhitelisted) {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "Denom not on whitelist")
	}

	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	err := k.OnRecvPacket(ctx, packet, data)
	if err != nil {
		acknowledgement = channeltypes.NewErrorAcknowledgement(err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			transfertypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
			sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
			sdk.NewAttribute(transfertypes.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", err != nil)),
		),
	)

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, acknowledgement.GetBytes(), nil
}

func GetMintedDenomFromPacket(packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) (string, bool) {
	// Note: Code and comments taken from SDK transfer keeper,
	// used here only to determine the token that will be minted.

	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {

		return "", true
	}

	// sender chain is the source, mint vouchers

	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + data.Denom

	// construct the denomination trace from the full raw denomination
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

	return denomTrace.IBCDenom(), false
}
