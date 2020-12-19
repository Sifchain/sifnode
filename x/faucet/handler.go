package faucet

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/faucet/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(bankKeeper types.BankKeeper, supply types.SupplyKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgRequestCoins:
			return handleMsgRequestCoins(ctx, bankKeeper, supplyKeeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgRequestCoins(ctx sdk.Context, bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper, msg MsgRequestCoins) (*sdk.Result, error) {
	err := bankKeeper.HasCoins(ctx, types.GetFaucetModuleAddress(), msg.Coins)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error()).Result()
	}
	keeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.Requester, msg.Coins)
	return nil, sdk.Result{}
}
