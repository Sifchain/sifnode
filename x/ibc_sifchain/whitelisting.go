package ibc_sifchain

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/ibc_sifchain/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/pkg/errors"
)

func OnRecvPacketWhiteListed(
	k Keeper,
	ctx sdk.Context,
	packet channeltypes.Packet,
	cdc codec.BinaryMarshaler,
) (*sdk.Result, []byte, error) {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	ok, err := isWhitelisted(ctx, packet, data, cdc)
	if !ok || err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom not on whitelist")
	}

	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	err = k.OnRecvPacket(ctx, packet, data)
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

func isWhitelisted(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData, cdc codec.BinaryMarshaler) (bool, error) {
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		// token originated on sifchain and is now being returned. This is allowed
		// For paths Sifchain -> X -> Sifchain return true
		// For paths Sifchain -> X -> Y -> Sifchain this condition is not triggered
		// No need to whitelist channel and port, We assume tokens will come back using the same channel they used to go across.
		// If Sifchain and Chain X have two channels running between them , and Token A uses channel 1 to go from sifchain to chain X . It needs to use channel 1 to come back.
		fmt.Printf("Returning to source | Denom : %s , SourcePort : %s , SourceChannel : %s ", data.Denom, packet.SourcePort, packet.SourceChannel)
		return true, nil
	}
	// Token did not originate on sifchain
	// In this case allow if all the conditions are met
	//    a) Token should belong to whitelist
	//    b) Token should be a direct transfer it should not have any jumps
	//    c) The port and channel should have been whitelisted
	// All the above conditions can be a met by whitelisting the trace hash of the token .
	whitelist, err := keeper.GetWhiteList(ctx, cdc)
	if err != nil {
		return false, errors.New("Whitelist not present")
	}
	denomTrace := transfertypes.ParseDenomTrace(data.Denom)
	if !whitelist.DenomWhitelist[denomTrace.Hash().String()] {
		return false, errors.New("< Token Channel Port > not present in whitelist")
	}
	fmt.Printf("Received whitelisted token %s , Hash %s ", data.Denom, denomTrace.Hash().String())
	return true, nil
}
