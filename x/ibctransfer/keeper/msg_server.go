package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"

	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
)

type msgServer struct {
	bankKeeper          types.BankKeeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	sdkMsgServer        types.MsgServer
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(sdkMsgServer types.MsgServer, bankKeeper types.BankKeeper, tokenRegistryKeeper tokenregistrytypes.Keeper) sdktransfertypes.MsgServer {
	return &msgServer{
		sdkMsgServer:        sdkMsgServer,
		bankKeeper:          bankKeeper,
		tokenRegistryKeeper: tokenRegistryKeeper,
	}
}

var _ sdktransfertypes.MsgServer = msgServer{}

// Transfer defines a rpc handler method for MsgTransfer.
func (srv msgServer) Transfer(goCtx context.Context, msg *sdktransfertypes.MsgTransfer) (*sdktransfertypes.MsgTransferResponse, error) {
	// Check export permission
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !srv.tokenRegistryKeeper.CheckDenomPermissions(ctx, msg.Token.Denom, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCEXPORT}) {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "denom cannot be exported")
	}

	// get token registry entry for sent token
	registryEntry := srv.tokenRegistryKeeper.GetDenom(sdk.UnwrapSDKContext(goCtx), msg.Token.Denom)
	// disallow direct transfers of denom aliases
	if registryEntry.UnitDenom != "" && registryEntry.UnitDenom != registryEntry.Denom {
		return nil, sdkerrors.Wrap(tokenregistrytypes.ErrPermissionDenied, "transfers of denom aliases are not yet supported")
	}

	// check if registry entry has an IBC counter party conversion to process
	if registryEntry.IbcCounterPartyDenom != "" && registryEntry.IbcCounterPartyDenom != registryEntry.Denom {
		sendAsRegistryEntry := srv.tokenRegistryKeeper.GetDenom(sdk.UnwrapSDKContext(goCtx), registryEntry.IbcCounterPartyDenom)
		if sendAsRegistryEntry.Decimals != 0 && registryEntry.Decimals > sendAsRegistryEntry.Decimals {
			token, tokenConversion := ConvertCoinsForTransfer(goCtx, msg, registryEntry, sendAsRegistryEntry)
			if token.Amount.Equal(sdk.NewInt(0)) {
				return nil, types.ErrAmountTooLowToConvert
			}
			err := PrepareToSendConvertedCoins(goCtx, msg, token, tokenConversion, srv.bankKeeper)
			if err != nil {
				return nil, err
			}
			msg.Token = tokenConversion
		}
	}

	return srv.sdkMsgServer.Transfer(goCtx, msg)
}

// Converts the coins requested for transfer into an amount that should be deducted from requested denom,
// and the Coins that should be minted in the new denom.
func ConvertCoinsForTransfer(goCtx context.Context, msg *sdktransfertypes.MsgTransfer, sendRegistryEntry tokenregistrytypes.RegistryEntry, sendAsRegistryEntry tokenregistrytypes.RegistryEntry) (sdk.Coin, sdk.Coin) {
	// calculate the conversion difference and reduce precision
	po := sendRegistryEntry.Decimals - sendAsRegistryEntry.Decimals
	decAmount := sdk.NewDecFromInt(msg.Token.Amount)
	convAmountDec := ReducePrecision(decAmount, po)

	convAmount := sdk.NewIntFromBigInt(convAmountDec.TruncateInt().BigInt())
	// create converted and sifchain tokens with corresponding denoms and amounts
	convToken := sdk.NewCoin(sendRegistryEntry.IbcCounterPartyDenom, convAmount)
	// increase convAmount precision to ensure amount deducted from address is the same that gets sent
	tokenAmountDec := IncreasePrecision(sdk.NewDecFromInt(convAmount), po)
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
	escrowAddress := sdktransfertypes.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)

	// escrow requested denom so it can be converted to the denom that will be sent out. It fails if balance insufficient.
	if err := bankKeeper.SendCoins(
		ctx, sender, escrowAddress, sdk.NewCoins(token),
	); err != nil {
		return err
	}

	// Mint into module account the new coins of the denom that will be sent via IBC
	err = bankKeeper.MintCoins(ctx, sdktransfertypes.ModuleName, sdk.NewCoins(convToken))
	if err != nil {
		return err
	}
	// Send minted coins (from module account) to senders address
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, sdktransfertypes.ModuleName, sender, sdk.NewCoins(convToken))
	if err != nil {
		return err
	}
	// Record conversion event, sender and coins
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeConvertTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, sdktransfertypes.ModuleName),
			sdk.NewAttribute(types.AttributeKeySentAmount, fmt.Sprintf("%v", token.Amount)),
			sdk.NewAttribute(types.AttributeKeySentDenom, token.Denom),
			sdk.NewAttribute(types.AttributeKeyConvertAmount, fmt.Sprintf("%v", convToken.Amount)),
			sdk.NewAttribute(types.AttributeKeyConvertDenom, convToken.Denom),
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
