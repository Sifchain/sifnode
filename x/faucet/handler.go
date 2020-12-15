package faucet

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func NewHandler(bk bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRequestCoins:
			return handleMsgRequestCoins(ctx, bk, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized faucet Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgRequestCoins(ctx sdk.Context, bk bank.Keeper, msg MsgRequestCoins) sdk.Result {
	_, _, err := bk.AddCoins(ctx, msg.Requester, msg.Coins)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	return sdk.Result{}
}
