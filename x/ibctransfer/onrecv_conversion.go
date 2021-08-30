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

// OnRecvPacketMaybeConvert will receive a transfer (after whitelisting is checked),
// and if the receive was successful,
// it will be converted into the unit_denom of the denom minted by the IBC transfer module.
func OnRecvPacketMaybeConvert(
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

	// get result of transfer receive
	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	err := sdkTransferKeeper.OnRecvPacket(ctx, packet, data)

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

	if err != nil {
		acknowledgement = channeltypes.NewErrorAcknowledgement(err.Error())

		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, acknowledgement.GetBytes(), nil
	}

	// Incoming coins were successfully minted onto the chain,
	// check if conversion to another denom is required.
	if ShouldConvertIncomingCoins(ctx, whitelistKeeper, packet, data) {
		receievedCoins, finalCoins := GetConvForIncomingCoins(ctx, whitelistKeeper, packet, data)
		err = ExecConvForIncomingCoins(ctx, receievedCoins, finalCoins, bankKeeper, packet, data)
		if err != nil {
			// Revert,
			// although this may cause packet to be relayed again.
			return nil, nil, sdkerrors.Wrap(sctransfertypes.ErrConvertingToUnitDenom, err.Error())
		}
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, acknowledgement.GetBytes(), nil
}
