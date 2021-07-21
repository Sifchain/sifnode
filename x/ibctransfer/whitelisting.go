package ibctransfer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	whitelisttypes "github.com/Sifchain/sifnode/x/whitelist/types"
)

func OnRecvPacketWhiteListed(
	ctx sdk.Context,
	sdkAppModule transfer.AppModule,
	whitelistKeeper whitelisttypes.Keeper,
	packet channeltypes.Packet,
) (*sdk.Result, []byte, error) {

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if !isRecvPacketAllowed(ctx, whitelistKeeper, packet, data) {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom not on whitelist")
	}

	return sdkAppModule.OnRecvPacket(ctx, packet)
}

func isRecvPacketAllowed(ctx sdk.Context, whitelistKeeper whitelisttypes.Keeper,
	packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) bool {

	isReturning := IsRecvPacketReturning(packet, data)

	denom := GetMintedDenomFromPacket(packet, data)
	isWhitelisted := IsWhitelisted(ctx, whitelistKeeper, denom)

	if isReturning || isWhitelisted {
		return true
	}

	return false
}

func GetMintedDenomFromPacket(packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) string {
	// Note: Code and comments taken from SDK transfer keeper,
	// used here only to determine the token that will be minted.

	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		// sender chain is not the source, unescrow tokens

		// remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := data.Denom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom := unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}

		return denom
	}

	// sender chain is the source, mint vouchers

	// since SendPacket did not prefix the denomination, we must prefix denomination here
	sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
	// NOTE: sourcePrefix contains the trailing "/"
	prefixedDenom := sourcePrefix + data.Denom

	// construct the denomination trace from the full raw denomination
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

	return denomTrace.IBCDenom()
}

func IsRecvPacketReturning(packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) bool {
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		// token originated on sifchain and is now being returned. This is allowed
		// For paths Sifchain -> X -> Sifchain return true
		// For paths Sifchain -> X -> Y -> Sifchain this condition is not triggered
		// No need to whitelist channel and port,
		// we assume tokens will come back using the same channel they used to go across.
		// If Sifchain and Chain X have two channels running between them,
		// and Token A uses channel 1 to go from sifchain to chain X . It needs to use channel 1 to come back.
		return true
	}

	return false
}

func IsWhitelisted(ctx sdk.Context, whitelistKeeper whitelisttypes.Keeper, denom string) bool {
	// In the case that token did not originate on sifchain,
	// allow if all the following conditions are met:
	//    a) Token should belong to whitelist
	//    b) Token should be a direct transfer it should not have any jumps
	//    c) The port and channel should have been whitelisted
	// All the above conditions can be a met by whitelisting the ibc/token that is minted on chain.

	return whitelistKeeper.IsDenomWhitelisted(ctx, denom)
}
