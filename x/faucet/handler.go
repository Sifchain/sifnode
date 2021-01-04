package faucet

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/faucet/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

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
	bank := keeper.GetBankKeeper()
	supply := keeper.GetSupplyKeeper()

	ok, err := keeper.CanRequest(ctx, msg.Requester, msg.Coins)
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
	ok, err = keeper.ExecuteRequest(ctx, msg.Requester, msg.Coins)
	if !ok || err != nil {
		return nil, err
	}
	return &sdk.Result{}, nil
}

// Handle the add coins message and send coins from the signers account to the module account.
func handleMsgAddCoins(ctx sdk.Context, keeper Keeper, msg types.MsgAddCoins) (*sdk.Result, error) {
	bank := keeper.GetBankKeeper()
	err := bank.SendCoins(ctx, msg.Signer, types.GetFaucetModuleAddress(), msg.Coins)
	if err != nil {
		return nil, errors.Wrap(err, types.ErrorAddingTokens.Error())
	}
	return &sdk.Result{}, nil
}
