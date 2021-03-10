package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// fee grants.
type RemoveFaucetFeeDecorator struct{}

func NewRemoveFacuetFeeDecorator() RemoveFaucetFeeDecorator {
	return RemoveFaucetFeeDecorator{}
}

var _ types.AnteDecorator = RemoveFaucetFeeDecorator{}

func (r RemoveFaucetFeeDecorator) AnteHandle(ctx types.Context, tx types.Tx, simulate bool, next types.AnteHandler) (newCtx types.Context, err error) {
	sigTx, ok := tx.(ante.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}
	msgs := sigTx.GetMsgs()
	fmt.Println("Executing ante Handler")
	for _, msg := range msgs {
		if msg.Type() == "request_coins" {
			//limit := ctx.GasMeter().Limit()
			fmt.Println("Inside Request coins")
			newDecCoin := types.NewDecCoinFromDec("rowan", types.NewDec(0))
			ctx = ctx.WithMinGasPrices(types.NewDecCoins(newDecCoin))
		}
	}
	return next(ctx, tx, simulate)
}

// NewRejectFeeGranterDecorator returns a new RejectFeeGranterDecorator.
