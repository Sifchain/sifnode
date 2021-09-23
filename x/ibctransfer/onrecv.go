package ibctransfer

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	"github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

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
	registry := whitelistKeeper.GetDenomWhitelist(ctx)
	mintedDenom := GetMintedDenomFromPacket(packet, data)
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
	convAmountDec := IncreasePrecision(decAmount, po)
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
	err = bankKeeper.SendCoinsFromAccountToModule(ctx, receiver, types.ModuleName, sdk.NewCoins(*incomingCoins))
	if err != nil {
		return err
	}
	// burn ibcdenom coins
	err = bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(*incomingCoins))
	if err != nil {
		return err
	}
	// unescrow original tokens
	escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())
	if err := bankKeeper.SendCoins(ctx, escrowAddress, receiver, sdk.NewCoins(*finalCoins)); err != nil {
		// NOTE: this error is only expected to occur given an unexpected bug or a malicious
		// counterparty module. The bug may occur in bank or any part of the code that allows
		// the escrow address to be drained. A malicious counterparty module could drain the
		// escrow address by allowing more tokens to be sent back then were escrowed.
		return sdkerrors.Wrap(err, "unable to unescrow original tokens")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeConvertReceived,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyPacketAmount, fmt.Sprintf("%v", incomingCoins.Amount)),
			sdk.NewAttribute(types.AttributeKeyPacketDenom, incomingCoins.Denom),
			sdk.NewAttribute(types.AttributeKeyConvertAmount, fmt.Sprintf("%v", finalCoins.Amount)),
			sdk.NewAttribute(types.AttributeKeyConvertDenom, finalCoins.Denom),
		),
	)
	return nil
}

func IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.MulTruncate(p)
}
