package ibctransfer

import (
	"fmt"

	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
)

func OnRecvPacketEnforceWhitelist(
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
	if !IsRecvPacketAllowed(ctx, whitelistKeeper, packet, data) {
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
	// Executes the actual receive, with potential conversion.
	return OnRecvPacketMaybeConvert(ctx, sdkTransferKeeper, whitelistKeeper, bankKeeper, packet)
}

func IsRecvPacketAllowed(ctx sdk.Context, whitelistKeeper tokenregistrytypes.Keeper, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) bool {
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		return true
	}
	denom := GetMintedDenomFromPacket(packet, data)
	registry := whitelistKeeper.GetDenomWhitelist(ctx)
	entry := whitelistKeeper.GetDenom(registry, denom)
	if entry == nil {
		return false
	}
	return whitelistKeeper.CheckDenomPermissions(entry, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT})
}

func GetMintedDenomFromPacket(packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) string {
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		denom := data.Denom[len(transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())):]
		denomTrace := transfertypes.ParseDenomTrace(denom)
		if denomTrace.Path != "" {
			return denomTrace.IBCDenom()
		}
		return denom
	} else {
		return transfertypes.ParseDenomTrace(transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel()) + data.Denom).IBCDenom()
	}
}
