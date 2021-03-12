package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	clpTypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

type SwapFeeChangeDecorator struct {
	ck keeper.Keeper
}

func NewSwapFeeChangeDecorator(ck keeper.Keeper) SwapFeeChangeDecorator {
	return SwapFeeChangeDecorator{
		ck: ck,
	}
}

var _ types.AnteDecorator = SwapFeeChangeDecorator{}

func (r SwapFeeChangeDecorator) AnteHandle(ctx types.Context, tx types.Tx, simulate bool, next types.AnteHandler) (newCtx types.Context, err error) {

	feeTx, ok := tx.(ante.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	msg := feeTx.GetMsgs()[0]
	switch msg := msg.(type) {
	case clpTypes.MsgSwap:
		feeInRowan := feeTx.GetFee()
		payer := feeTx.FeePayer()
		coinsBalance := r.ck.GetBankKeeper().GetCoins(ctx, payer)
		payerHasRowan := true
		if coinsBalance.AmountOf(clpTypes.GetSettlementAsset().Symbol).IsZero() {
			payerHasRowan = false
		}
		if !payerHasRowan {
			_ = EnrichPayerWithRowan(r.ck, ctx, msg, payer, feeInRowan)
		}
	default:
		fmt.Println("Unreachable code :")
	}

	return next(ctx, tx, simulate)
}

func EnrichPayerWithRowan(ck keeper.Keeper, ctx types.Context, msg clpTypes.MsgSwap, payer types.Address, feeInRowan types.Coins) (err error) {
	_, err = ck.GetLiquidityProvider(ctx, msg.ReceivedAsset.Symbol, payer.String())
	if err != nil {
		return
	}
	requiredRowan := feeInRowan.AmountOf(clpTypes.GetSettlementAsset().Symbol)
	payerAccAddress, err := types.AccAddressFromHex(payer.String())
	if err != nil {
		return
	}
	swapMSG := clpTypes.MsgSwap{
		Signer:             payerAccAddress, // Same as msg.signer
		SentAsset:          msg.ReceivedAsset,
		ReceivedAsset:      clpTypes.GetSettlementAsset(),
		SentAmount:         types.Uint{},
		MinReceivingAmount: types.NewUintFromBigInt(requiredRowan.BigInt()),
	}
	_, err = handleMsgSwap(ctx, ck, swapMSG)
	if err != nil {
		return
	}
	return nil
}
