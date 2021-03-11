package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
)

type SwapFeeChangeDecorator struct{}

func NewSwapFeeChangeDecorator() SwapFeeChangeDecorator {
	return SwapFeeChangeDecorator{}
}

var _ types.AnteDecorator = SwapFeeChangeDecorator{}

func (r SwapFeeChangeDecorator) AnteHandle(ctx types.Context, tx types.Tx, simulate bool, next types.AnteHandler) (newCtx types.Context, err error) {
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case MsgSwap:
		currentGasPrice := ctx.MinGasPrices().AmountOf(GetSettlementAsset().Symbol)
		multiplier, _ := types.NewDecFromStr("0.001")
		currentGasPrice = currentGasPrice.Mul(multiplier)
		fmt.Println("Current Gas Price in rowan", currentGasPrice)
		fmt.Println("Want to charge gas in :", msg.ReceivedAsset)
		newDecCoin := types.NewDecCoinFromDec(msg.ReceivedAsset.Symbol, currentGasPrice)
		ctx = ctx.WithMinGasPrices(types.NewDecCoins(newDecCoin))
	default:
		fmt.Println("Unreachable code :")
	}

	return next(ctx, tx, simulate)
}
