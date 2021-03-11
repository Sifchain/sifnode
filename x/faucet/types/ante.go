package types

import (
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

type RemoveFaucetFeeDecorator struct{}

func NewRemoveFaucetFeeDecorator() RemoveFaucetFeeDecorator {
	return RemoveFaucetFeeDecorator{}
}

var _ types.AnteDecorator = RemoveFaucetFeeDecorator{}

func (r RemoveFaucetFeeDecorator) AnteHandle(ctx types.Context, tx types.Tx, simulate bool, next types.AnteHandler) (newCtx types.Context, err error) {
	sigTx, ok := tx.(ante.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}
	msgs := sigTx.GetMsgs()
	for _, msg := range msgs {
		if msg.Type() == "request_coins" && ctx.ChainID() != "sifchain" && Profile == TESTNET {
			newDecCoin := types.NewDecCoinFromDec(clp.GetSettlementAsset().Symbol, types.NewDec(0))
			ctx = ctx.WithMinGasPrices(types.NewDecCoins(newDecCoin))
		}
	}
	return next(ctx, tx, simulate)
}
