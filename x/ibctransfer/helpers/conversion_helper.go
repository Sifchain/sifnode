package helpers

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"

	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

// ConvertCoinsForTransfer Converts the coins requested for transfer into an amount that should be deducted from requested denom,
// and the Coins that should be minted in the new denom.

// TODO only used in tests , remove this function completely
func ConvertCoinsForTransfer(msg *sdktransfertypes.MsgTransfer, sendRegistryEntry *tokenregistrytypes.RegistryEntry,
	sendAsRegistryEntry *tokenregistrytypes.RegistryEntry) (sdk.Coin, sdk.Coin) {
	// calculate the conversion difference and reduce precision
	po := uint64(sendRegistryEntry.Decimals - sendAsRegistryEntry.Decimals)
	decAmount := sdk.NewDecFromInt(msg.Token.Amount)
	convAmountDec := ReducePrecision(decAmount, po)
	convAmount := sdk.NewIntFromBigInt(convAmountDec.TruncateInt().BigInt())
	// create converted and Sifchain tokens with corresponding denoms and amounts
	convToken := sdk.NewCoin(sendRegistryEntry.IbcCounterpartyDenom, convAmount)
	// increase convAmount precision to ensure amount deducted from address is the same that gets sent
	tokenAmountDec := IncreasePrecision(sdk.NewDecFromInt(convAmount), po)
	tokenAmount := sdk.NewIntFromBigInt(tokenAmountDec.TruncateInt().BigInt())
	token := sdk.NewCoin(msg.Token.Denom, tokenAmount)
	return token, convToken
}

// PrepareToSendConvertedCoins moves outgoing tokens into the denom that will be sent via IBC.
// The requested tokens will be escrowed, and the new denom to send over IBC will be minted in the senders account.
func PrepareToSendConvertedCoins(goCtx context.Context, msg *sdktransfertypes.MsgTransfer, token sdk.Coin, convToken sdk.Coin, bankKeeper sctransfertypes.BankKeeper) error {
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

func IsRecvPacketAllowed(ctx sdk.Context, whitelistKeeper tokenregistrytypes.Keeper, packet channeltypes.Packet, data sdktransfertypes.FungibleTokenPacketData, mintedDenomEntry *tokenregistrytypes.RegistryEntry) bool {
	if sdktransfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		return true
	}
	return whitelistKeeper.CheckEntryPermissions(mintedDenomEntry, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT})
}

func GetMintedDenomFromPacket(packet channeltypes.Packet, data sdktransfertypes.FungibleTokenPacketData) string {
	if sdktransfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		denom := data.Denom[len(sdktransfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())):]
		denomTrace := sdktransfertypes.ParseDenomTrace(denom)
		if denomTrace.Path != "" {
			return denomTrace.IBCDenom()
		}
		return denom
	}
	return sdktransfertypes.ParseDenomTrace(sdktransfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel()) + data.Denom).IBCDenom()
}

func ConvertIncomingCoins(amount string, diff uint64) (sdk.Int, error) {
	intAmount, ok := sdk.NewIntFromString(amount)
	if !ok {
		return sdk.ZeroInt(), errors.New("Failed to convert to int from string")
	}
	return sdk.NewIntFromBigInt(IncreasePrecision(sdk.NewDecFromInt(intAmount), diff).TruncateInt().BigInt()), nil
}

func ExecConvForIncomingCoins(
	ctx sdk.Context,
	bankKeeper sdktransfertypes.BankKeeper,
	mintedDenomEntry *tokenregistrytypes.RegistryEntry,
	convertToDenomEntry *tokenregistrytypes.RegistryEntry,
	packet channeltypes.Packet,
	data sdktransfertypes.FungibleTokenPacketData,
) error {
	// decode the receiver address
	receiver, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}
	amount, ok := sdk.NewIntFromString(data.Amount)
	if !ok {
		return errors.New("Unable to get string amount")
	}
	incomingCoins := sdk.NewCoins(sdk.NewCoin(mintedDenomEntry.Denom, amount))
	// send ibcdenom coins from account to module
	err = bankKeeper.SendCoinsFromAccountToModule(ctx, receiver, sctransfertypes.ModuleName, incomingCoins)
	if err != nil {
		return err
	}
	// burn ibcdenom coins
	err = bankKeeper.BurnCoins(ctx, sctransfertypes.ModuleName, incomingCoins)
	if err != nil {
		return err
	}
	convAmount := amount
	finalCoins := sdk.NewCoins(sdk.NewCoin(convertToDenomEntry.Denom, convAmount))
	if convertToDenomEntry.Decimals > mintedDenomEntry.Decimals {
		diff := uint64(convertToDenomEntry.Decimals - mintedDenomEntry.Decimals)
		// This is the reduced precision xToken coming in , so we know for sure conversion to uint64 will not cause problems
		convAmount, err = ConvertIncomingCoins(data.Amount, diff)
		if err != nil {
			return err
		}
		finalCoins = sdk.NewCoins(sdk.NewCoin(convertToDenomEntry.Denom, convAmount))
	}
	// unescrow original tokens
	escrowAddress := sctransfertypes.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())

	if err := bankKeeper.SendCoins(ctx, escrowAddress, receiver, finalCoins); err != nil {
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
			sdk.NewAttribute(sctransfertypes.AttributeKeyPacketAmount, fmt.Sprintf("%v", data.Amount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeyPacketDenom, mintedDenomEntry.Denom),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertAmount, fmt.Sprintf("%v", convAmount)),
			sdk.NewAttribute(sctransfertypes.AttributeKeyConvertDenom, convertToDenomEntry.Denom),
		),
	)
	return nil
}

func IncreasePrecision(dec sdk.Dec, po uint64) sdk.Dec {
	p := sdk.NewDec(10).Power(po)
	return dec.MulTruncate(p)
}

func ReducePrecision(dec sdk.Dec, po uint64) sdk.Dec {
	p := sdk.NewDec(10).Power(po)
	return dec.QuoTruncate(p)
}
