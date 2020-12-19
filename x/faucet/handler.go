package faucet

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/Sifchain/sifnode/x/faucet/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(bankKeeper types.BankKeeper, supply types.SupplyKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgRequestCoins:
			return handleMsgRequestCoins(ctx, bankKeeper, supply, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgRequestCoins(ctx sdk.Context, bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper, msg types.MsgRequestCoins) (*sdk.Result, error) {
	ok := bankKeeper.HasCoins(ctx, types.GetFaucetModuleAddress(), msg.Coins)
	if !ok {
		return nil, errors.New("Not enough balance")
	}
	supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.Requester, msg.Coins)
	return &sdk.Result{},nil
}
