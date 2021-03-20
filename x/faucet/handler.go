package faucet

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/faucet/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

// NewHandler creates an sdk.Handler for all the faucet type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgRequestCoins:
			return handleMsgRequestCoins(ctx, keeper, msg)
		case types.MsgAddCoins:
			return handleMsgAddCoins(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle the incoming request message and distribute coins from the module to the requesters account.
// Will need to update this in the future with some distribution limitations
func handleMsgRequestCoins(ctx sdk.Context, keeper Keeper, msg types.MsgRequestCoins) (*sdk.Result, error) {
	if ctx.ChainID() != "sifchain" {
		bank := keeper.GetBankKeeper()
		supply := keeper.GetSupplyKeeper()

		ok, err := keeper.CanRequest(ctx, msg.Requester.String(), msg.Coins)
		if !ok || err != nil {
			return nil, err
		}
		ok = bank.HasCoins(ctx, types.GetFaucetModuleAddress(), msg.Coins)
		if !ok {
			return nil, types.NotEnoughBalance
		}
		err = supply.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.Requester, msg.Coins)
		if err != nil {
			return nil, errors.Wrap(err, types.ErrorRequestingTokens.Error())
		}
		ok, err = keeper.ExecuteRequest(ctx, msg.Requester.String(), msg.Coins)
		if !ok || err != nil {
			return nil, err
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeRequestCoins,
				sdk.NewAttribute(types.AttributeKeyFaucet, types.ModuleName),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Requester.String()),
			),
		})
		return &sdk.Result{Events: ctx.EventManager().Events()}, nil
	}
	return nil, nil
}

// Handle the add coins message and send coins from the signers account to the module account.
func handleMsgAddCoins(ctx sdk.Context, keeper Keeper, msg types.MsgAddCoins) (*sdk.Result, error) {
	if ctx.ChainID() != "sifchain" {
		supply := keeper.GetSupplyKeeper()
		err := supply.SendCoinsFromAccountToModule(ctx, msg.Signer, types.ModuleName, msg.Coins)
		if err != nil {
			return nil, errors.Wrap(err, types.ErrorAddingTokens.Error())
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeAddCoins,
				sdk.NewAttribute(types.AttributeKeyFaucet, types.ModuleName),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
			),
		})
		return &sdk.Result{Events: ctx.EventManager().Events()}, nil
	}
	return nil, nil
}
