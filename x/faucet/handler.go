package faucet

import (
	"fmt"
	"reflect"

	"github.com/Sifchain/sifnode/x/faucet/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(bankKeeper types.BankKeeper, supply types.SupplyKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRequestCoins:
			return handleMsgRequestCoins(ctx, bankKeeper, supplyKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized faucet Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Steps for handling message
// Check if account has coins via bank module
// Send coins to the requester
func handleMsgRequestCoins(ctx sdk.Context, bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper, msg MsgRequestCoins) sdk.Result {
	err := bankKeeper.HasCoins(ctx, types.GetFaucetModuleAddress, msg.Coins)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	keeper.SendCoins(ctx, types.GetFaucetModuleAddress, msg.Requester, msg.Coins)
	return sdk.Result{}
}
