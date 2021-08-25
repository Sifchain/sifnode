package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/Sifchain/sifnode/x/ibctransfer/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"

	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
)

type msgServer struct {
	bankKeeper          bankkeeper.Keeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	sdkMsgServer        sdktransfertypes.MsgServer
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(sdkMsgServer sdktransfertypes.MsgServer, bankKeeper bankkeeper.Keeper, tokenRegistryKeeper tokenregistrytypes.Keeper) sdktransfertypes.MsgServer {
	return &msgServer{
		sdkMsgServer:        sdkMsgServer,
		bankKeeper:          bankKeeper,
		tokenRegistryKeeper: tokenRegistryKeeper,
	}
}

var _ sdktransfertypes.MsgServer = msgServer{}

// Transfer defines a rpc handler method for MsgTransfer.
func (srv msgServer) Transfer(goCtx context.Context, msg *sdktransfertypes.MsgTransfer) (*sdktransfertypes.MsgTransferResponse, error) {
	// get token registry entry for sent token
	registryEntry := srv.tokenRegistryKeeper.GetDenom(sdk.UnwrapSDKContext(goCtx), msg.Token.Denom)
	// check if registry entry has an IBC counter party conversion to process
	if registryEntry.IbcCounterPartyDenom != "" && registryEntry.IbcCounterPartyDenom != registryEntry.Denom {
		sendAsRegistryEntry := srv.tokenRegistryKeeper.GetDenom(sdk.UnwrapSDKContext(goCtx), registryEntry.IbcCounterPartyDenom)
		if registryEntry.Decimals > sendAsRegistryEntry.Decimals {
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

func PrepareToSendConvertedCoins(goCtx context.Context, msg *sdktransfertypes.MsgTransfer, token sdk.Coin, convToken sdk.Coin, bankKeeper bankkeeper.Keeper) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	// Deduct requested denom so it can be converted to the denom that will be sent out
	err = bankKeeper.SendCoinsFromAccountToModule(ctx, sender, sdktransfertypes.ModuleName, sdk.NewCoins(token))
	if err != nil {
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
			tokenregistrytypes.EventTypeConvertTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, sdktransfertypes.ModuleName),
			sdk.NewAttribute(tokenregistrytypes.AttributeKeySentAmount, fmt.Sprintf("%v", token.Amount)),
			sdk.NewAttribute(tokenregistrytypes.AttributeKeySentDenom, token.Denom),
			sdk.NewAttribute(tokenregistrytypes.AttributeKeyConvertAmount, fmt.Sprintf("%v", convToken.Amount)),
			sdk.NewAttribute(tokenregistrytypes.AttributeKeyConvertDenom, convToken.Denom),
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
