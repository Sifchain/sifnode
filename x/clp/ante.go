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
		fmt.Println("Balance Before : ", r.ck.GetBankKeeper().GetCoins(ctx, payer))
		coinsBalance := r.ck.GetBankKeeper().GetCoins(ctx, payer)
		payerHasRowan := false
		if coinsBalance.AmountOf(clpTypes.GetSettlementAsset().Symbol).IsZero() {
			payerHasRowan = false
		}
		if !payerHasRowan {
			err = EnrichPayerWithRowan(r.ck, ctx, msg, payer, feeInRowan)
			if err != nil {
				fmt.Println("Error :", err)
			}
		}
		fmt.Println("Balance After : ", r.ck.GetBankKeeper().GetCoins(ctx, payer))
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
	pool, err := ck.GetPool(ctx, msg.ReceivedAsset.Symbol)
	if err != nil {
		return
	}
	// x - amount of cToken
	// X - cToken Balance
	// Y - Rowan Balance
	// S - amount of Rowan
	sendAmount, err := keeper.ReverseSwap(pool.ExternalAssetBalance, pool.NativeAssetBalance, types.Uint(requiredRowan))
	if err != nil {
		return
	}

	swapMSG := clpTypes.MsgSwap{
		Signer:             msg.Signer,
		SentAsset:          msg.ReceivedAsset,
		ReceivedAsset:      clpTypes.GetSettlementAsset(),
		SentAmount:         sendAmount,
		MinReceivingAmount: types.NewUintFromBigInt(requiredRowan.BigInt()),
	}
	res, err := handleMsgSwap(ctx, ck, swapMSG)
	fmt.Println("Swap Result :", res, err)
	if err != nil {
		return
	}
	return nil
}
