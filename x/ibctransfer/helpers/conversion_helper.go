package helpers

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	"github.com/Sifchain/sifnode/x/ibctransfer/types"
	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

// ConvertCoinsForTransfer Converts the coins requested for transfer into an amount that should be deducted from requested denom,
// and the Coins that should be minted in the new denom.
func ConvertCoinsForTransfer(msg *sdktransfertypes.MsgTransfer, sendRegistryEntry *tokenregistrytypes.RegistryEntry,
	sendAsRegistryEntry *tokenregistrytypes.RegistryEntry) (sdk.Coin, sdk.Coin) {
	// calculate the conversion difference and reduce precision
	po := sendRegistryEntry.Decimals - sendAsRegistryEntry.Decimals
	decAmount := sdk.NewDecFromInt(msg.Token.Amount)
	convAmountDec := ReducePrecision(decAmount, int64(po))
	convAmount := sdk.NewIntFromBigInt(convAmountDec.TruncateInt().BigInt())
	// create converted and Sifchain tokens with corresponding denoms and amounts
	convToken := sdk.NewCoin(sendRegistryEntry.IbcCounterpartyDenom, convAmount)
	// increase convAmount precision to ensure amount deducted from address is the same that gets sent
	tokenAmountDec := IncreasePrecision(sdk.NewDecFromInt(convAmount), int64(po))
	tokenAmount := sdk.NewIntFromBigInt(tokenAmountDec.TruncateInt().BigInt())
	token := sdk.NewCoin(msg.Token.Denom, tokenAmount)
	return token, convToken
}

// PrepareToSendConvertedCoins moves outgoing tokens into the denom that will be sent via IBC.
// The requested tokens will be escrowed, and the new denom to send over IBC will be minted in the senders account.
func PrepareToSendConvertedCoins(goCtx context.Context, msg *sdktransfertypes.MsgTransfer, token sdk.Coin, convToken sdk.Coin, bankKeeper types.BankKeeper) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	// create the escrow address for the tokens
	escrowAddress := sctransfertypes.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)
	// escrow requested denom so it can be converted to the denom that will be sent out. It fails if balance insufficient.
	if err = bankKeeper.SendCoins(ctx, sender, escrowAddress, sdk.NewCoins(token)); err != nil {
		return err
	}
	convCoins := sdk.NewCoins(convToken)
	// Mint into module account the new coins of the denom that will be sent via IBC
	err = bankKeeper.MintCoins(ctx, sctransfertypes.ModuleName, convCoins)
	if err != nil {
		return err
	}
	// Send minted coins (from module account) to senders address
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, sctransfertypes.ModuleName, sender, convCoins)
	if err != nil {
		return err
	}
	// Record conversion event, sender and coins
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sctransfertypes.EventTypeConvertTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, sctransfertypes.ModuleName),
			sdk.NewAttribute(sctransfertypes.AttributeKeySentAmount, fmt.Sprintf("%v", token.Amount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeySentDenom, token.Denom),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertAmount, fmt.Sprintf("%v", convToken.Amount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertDenom, convToken.Denom),
		),
	)
	return nil
}

func GetMintedDenomFromPacket(packet channeltypes.Packet, data sdktransfertypes.FungibleTokenPacketData) string {
	if sdktransfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		denom := data.Denom[len(sdktransfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())):]
		denomTrace := sdktransfertypes.ParseDenomTrace(denom)
		if denomTrace.Path != "" {
			return denomTrace.IBCDenom()
		}
		return denom
	} else {
		return sdktransfertypes.ParseDenomTrace(sdktransfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel()) + data.Denom).IBCDenom()
	}
}

// GetConvForIncomingCoins returns 1) the coins that are being received via IBC,
// which need to be deducted from that denom when converting to final denom,
// and 2) the coins that need to be added to the final denom.
func GetConvForIncomingCoins(
	ctx sdk.Context,
	whitelistKeeper tokenregistrytypes.Keeper,
	packet channeltypes.Packet,
	data sdktransfertypes.FungibleTokenPacketData,
) (*sdk.Coin, *sdk.Coin) {
	// Get the denom that will be minted by sdk transfer module,
	// so that it can be converted to the denom it should be stored as.
	// For a native token that has been returned, this will just be a base_denom,
	// which will be on the whitelist.
	mintedDenom := GetMintedDenomFromPacket(packet, data)
	registry := whitelistKeeper.GetDenomWhitelist(ctx)
	// get token registry entry for received denom
	mintedDenomEntry := whitelistKeeper.GetDenom(registry, mintedDenom)
	// convert to unit_denom
	if mintedDenomEntry == nil {
		// noop, should prevent getting here.
		return nil, nil
	}
	convertToDenomEntry := whitelistKeeper.GetDenom(registry, mintedDenomEntry.UnitDenom)
	if convertToDenomEntry == nil {
		// noop, should prevent getting here.
		return nil, nil
	}
	// get the token amount from the packet data
	// Calculate the conversion difference for increasing precision.
	po := convertToDenomEntry.Decimals - mintedDenomEntry.Decimals
	if po <= 0 {
		// Shortcut to prevent crash if po <= 0
		return nil, nil
	}
	decAmount := sdk.NewDecFromInt(sdk.NewIntFromUint64(data.Amount))
	convAmountDec := IncreasePrecision(decAmount, int64(po))
	convAmount := sdk.NewIntFromBigInt(convAmountDec.TruncateInt().BigInt())
	// create converted and ibc tokens with corresponding denoms and amounts
	convertToCoins := sdk.NewCoin(convertToDenomEntry.Denom, convAmount)
	mintedCoins := sdk.NewCoin(mintedDenom, sdk.NewIntFromUint64(data.Amount))
	return &mintedCoins, &convertToCoins
}

func ExecConvForIncomingCoins(
	ctx sdk.Context,
	incomingCoins *sdk.Coin,
	finalCoins *sdk.Coin,
	bankKeeper sdktransfertypes.BankKeeper,
	packet channeltypes.Packet,
	data sdktransfertypes.FungibleTokenPacketData,
) error {
	// decode the receiver address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}
	// send ibcdenom coins from account to module
	err = bankKeeper.SendCoinsFromAccountToModule(ctx, receiver, sctransfertypes.ModuleName, sdk.NewCoins(*incomingCoins))
	if err != nil {
		return err
	}
	// burn ibcdenom coins
	err = bankKeeper.BurnCoins(ctx, sctransfertypes.ModuleName, sdk.NewCoins(*incomingCoins))
	if err != nil {
		return err
	}
	// unescrow original tokens
	escrowAddress := sctransfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	if err := bankKeeper.SendCoins(ctx, escrowAddress, receiver, sdk.NewCoins(*finalCoins)); err != nil {
		// NOTE: this error is only expected to occur given an unexpected bug or a malicious
		// counterparty module. The bug may occur in bank or any part of the code that allows
		// the escrow address to be drained. A malicious counterparty module could drain the
		// escrow address by allowing more tokens to be sent back then were escrowed.
		return sdkerrors.Wrap(err, "unable to unescrow original tokens")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sctransfertypes.EventTypeConvertReceived,
			sdk.NewAttribute(sdk.AttributeKeyModule, sctransfertypes.ModuleName),
			sdk.NewAttribute(sctransfertypes.AttributeKeyPacketAmount, fmt.Sprintf("%v", incomingCoins.Amount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeyPacketDenom, incomingCoins.Denom),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertAmount, fmt.Sprintf("%v", finalCoins.Amount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertDenom, finalCoins.Denom),
		),
	)
	return nil
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
	convAmountDec := IncreasePrecision(decAmount, int64(po))
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

func IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.MulTruncate(p)
}

func ReducePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.QuoTruncate(p)
}
