package ibctransfer

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"

	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

// OnRecvPacketWhitelistConvert receives a transfer, check if the denom is whitelisted, and converts it
// to match unit_denom decimals if conversion is needed.
func OnRecvPacketWhitelistConvert(
	ctx sdk.Context,
	sdkTransferKeeper sctransfertypes.SDKTransferKeeper,
	whitelistKeeper tokenregistrytypes.Keeper,
	bankKeeper transfertypes.BankKeeper,
	packet channeltypes.Packet,
) channeltypes.Acknowledgement {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		acknowledgement := channeltypes.NewErrorAcknowledgement(err)
		return acknowledgement
	}
	err := sdkTransferKeeper.OnRecvPacket(ctx, packet, data)
	if err != nil {
		acknowledgement := channeltypes.NewErrorAcknowledgement(err)
		return acknowledgement
	}
	// Get the denom that will be minted by sdk transfer module,
	// so that it can be converted to the denom it should be stored as.
	// For a native token that has been returned, this will just be a base_denom,
	// which will be on the whitelist.
	mintedDenom := helpers.GetMintedDenomFromPacket(packet, data)
	registry := whitelistKeeper.GetRegistry(ctx)
	mintedDenomEntry, err := whitelistKeeper.GetEntry(registry, mintedDenom)
	if err != nil || !helpers.IsRecvPacketAllowed(ctx, whitelistKeeper, packet, data, mintedDenomEntry) {
		acknowledgement := channeltypes.NewErrorAcknowledgement(
			sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom not whitelisted"),
		)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				transfertypes.EventTypePacket,
				sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
				sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
				sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
				sdk.NewAttribute(transfertypes.AttributeKeyAmount, data.Amount),
				sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", false)),
			),
		)
		return acknowledgement
	}
	// TODO Add entries fpr Non-X versions of tokens to tokenRegistry
	convertToDenomEntry, err := whitelistKeeper.GetEntry(registry, mintedDenomEntry.UnitDenom)
	if err == nil && convertToDenomEntry.Decimals > 0 && mintedDenomEntry.Decimals > 0 && convertToDenomEntry.Decimals > mintedDenomEntry.Decimals {
		err = helpers.ExecConvForIncomingCoins(ctx, bankKeeper, mintedDenomEntry, convertToDenomEntry, packet, data)
		// Revert, although this may cause packet to be relayed again.
		if err != nil {
			acknowledgement := channeltypes.NewErrorAcknowledgement(
				sdkerrors.Wrapf(sctransfertypes.ErrConvertingToUnitDenom, err.Error()),
			)
			return acknowledgement
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			transfertypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
			sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
			sdk.NewAttribute(transfertypes.AttributeKeyAmount, data.Amount),
			sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", err == nil)),
		),
	)
	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	return acknowledgement
}
