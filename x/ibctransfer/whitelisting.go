package ibctransfer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func OnRecvPacketWhiteListed(
	ctx sdk.Context,
	sdkAppModule transfer.AppModule,
	whitelistKeeper tokenregistrytypes.Keeper,
	bankKeeper types.BankKeeper,
	packet channeltypes.Packet,
) (*sdk.Result, []byte, error) {

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if !isRecvPacketAllowed(ctx, whitelistKeeper, packet, data) {
		acknowledgement := channeltypes.NewErrorAcknowledgement(
			sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "denom not on whitelist").Error(),
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				transfertypes.EventTypePacket,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName), // is this ethbridge on purpose?
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
	// get result of transfer receive
	recvResult, resBytes, err := sdkAppModule.OnRecvPacket(ctx, packet)
	// if no error and packet is returning and needs conversion: convert
	if err != nil && IsRecvPacketReturning(packet, data) && checkRecvConvert(ctx, whitelistKeeper, packet, data) {

		ibcToken, convToken := convertRecvDenom(ctx, whitelistKeeper, packet, data)
		recvResult, err = sendConvertRecvDenom(ctx, ibcToken, convToken, bankKeeper, data)

	}

	// otherwise return
	return recvResult, resBytes, err
}

func checkRecvConvert(ctx sdk.Context, whitelistKeeper tokenregistrytypes.Keeper,
	packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) bool {
	denom := GetMintedDenomFromPacket(packet, data)
	// get token registry entry for received token
	registryEntry := whitelistKeeper.GetIBCDenom(ctx, denom)
	return registryEntry.IbcDenom != "" && registryEntry.Decimals > registryEntry.IbcDecimals
}

func convertRecvDenom(ctx sdk.Context, whitelistKeeper tokenregistrytypes.Keeper,
	packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) (sdk.Coin, sdk.Coin) {
	denom := GetMintedDenomFromPacket(packet, data)
	// get token registry entry for received token
	registryEntry := whitelistKeeper.GetIBCDenom(ctx, denom)
	// check if registry entry has an IBC decimal field
	// calculate conversion
	po := registryEntry.Decimals - registryEntry.IbcDecimals
	decAmount := sdk.NewDecFromInt(sdk.NewIntFromUint64(data.Amount))
	convAmountDec := IncreasePrecision(decAmount, po)
	convAmount := sdk.NewIntFromBigInt(convAmountDec.RoundInt().BigInt())
	convToken := sdk.NewCoin(registryEntry.Denom, convAmount)
	ibcToken := sdk.NewCoin(denom, sdk.NewIntFromUint64(data.Amount))

	return ibcToken, convToken
}
func sendConvertRecvDenom(ctx sdk.Context, ibcToken sdk.Coin, convToken sdk.Coin, bankKeeper types.BankKeeper, data transfertypes.FungibleTokenPacketData) (*sdk.Result, error) {
	// denom := GetMintedDenomFromPacket(packet, data)
	// // get token registry entry for received token
	// registryEntry := whitelistKeeper.GetIBCDenom(ctx, denom)
	// // check if registry entry has an IBC decimal field
	// // calculate conversion
	// po := registryEntry.Decimals - registryEntry.IbcDecimals
	// decAmount := sdk.NewDecFromInt(sdk.NewIntFromUint64(data.Amount))
	// convAmountDec := IncreasePrecision(sdk.NewDecFromBigInt(decAmount.BigInt()), po)
	// convAmount := sdk.NewIntFromBigInt(convAmountDec.RoundInt().BigInt())

	// convToken := sdk.NewCoin(registryEntry.Denom, convAmount)
	// // decode the receiver address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return nil, err
	}
	// send ibcdenom coins from account to module
	// ibcToken := sdk.NewCoin(denom, sdk.NewIntFromUint64(data.Amount))
	err = bankKeeper.SendCoinsFromAccountToModule(ctx, receiver, transfertypes.ModuleName, sdk.NewCoins(ibcToken))
	if err != nil {
		return nil, err
	}
	// send coins from module account to address
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, receiver, sdk.NewCoins(convToken))
	if err != nil {
		return nil, err
	}
	// burn ibcdenom coins
	err = bankKeeper.BurnCoins(ctx, transfertypes.ModuleName, sdk.NewCoins(ibcToken))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			transfertypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, convToken.Denom),
			sdk.NewAttribute(transfertypes.AttributeKeyAmount, fmt.Sprintf("%d", convToken.Amount)),
			// sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", err == nil)),  ??
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.Mul(p)
}

func isRecvPacketAllowed(ctx sdk.Context, whitelistKeeper tokenregistrytypes.Keeper,
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

func IsWhitelisted(ctx sdk.Context, whitelistKeeper tokenregistrytypes.Keeper, denom string) bool {
	// In the case that token did not originate on sifchain,
	// allow if all the following conditions are met:
	//    a) Token should belong to whitelist
	//    b) Token should be a direct transfer it should not have any jumps
	//    c) The port and channel should have been whitelisted
	// All the above conditions can be a met by whitelisting the ibc/token that is minted on chain.

	return whitelistKeeper.IsDenomWhitelisted(ctx, denom)
}
