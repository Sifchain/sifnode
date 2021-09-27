package ibctransfer

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

// OnTimeoutPacket() refunds the sender since the original packet sent was
// never received and has been timed out.
//
// sdkkeeper.refundPacketToken will unescrow and send back the tokens back to sender
// if the sending chain was the source chain. Otherwise, the sent tokens
// were burnt in the original send so new tokens are minted and sent to
// the sending address.
//

// OnTimeoutMaybeConvert potentially does a conversion from the denom that was sent out,
// back to the original unit_denom.
func OnTimeoutMaybeConvert(
	ctx sdk.Context,
	sdkTransferKeeper sctransfertypes.SDKTransferKeeper,
	whitelistKeeper tokenregistrytypes.Keeper,
	bankKeeper transfertypes.BankKeeper,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	// refund tokens
	if err := sdkTransferKeeper.OnTimeoutPacket(ctx, packet, data); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			transfertypes.EventTypeTimeout,
			sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
			sdk.NewAttribute(transfertypes.AttributeKeyRefundReceiver, data.Sender),
			sdk.NewAttribute(transfertypes.AttributeKeyRefundDenom, data.Denom),
			sdk.NewAttribute(transfertypes.AttributeKeyRefundAmount, fmt.Sprintf("%d", data.Amount)),
		),
	)
	registry := whitelistKeeper.GetDenomWhitelist(ctx)
	denomEntry := whitelistKeeper.GetDenom(registry, data.Denom)
	if denomEntry != nil && denomEntry.Decimals > 0 && denomEntry.UnitDenom != "" {
		convertToDenomEntry := whitelistKeeper.GetDenom(registry, denomEntry.UnitDenom)
		if convertToDenomEntry != nil && convertToDenomEntry.Decimals > denomEntry.Decimals {
			diff := uint64(convertToDenomEntry.Decimals - denomEntry.Decimals)
			convAmount := helpers.ConvertIncomingCoins(ctx, whitelistKeeper, data.Amount, diff)
			incomingCoins := sdk.NewCoin(data.Denom, sdk.NewIntFromUint64(data.Amount))
			finalCoins := sdk.NewCoin(convertToDenomEntry.Denom, convAmount)
			err := helpers.ExecConvForRefundCoins(ctx, &incomingCoins, &finalCoins, bankKeeper, packet, data)
			if err != nil {
				return nil, err
			}
			return &sdk.Result{
				Events: ctx.EventManager().Events().ToABCIEvents(),
			}, nil
		}
	}
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
