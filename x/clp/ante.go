package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	clpTypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/pkg/errors"
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
		payer := feeTx.FeePayer()
		if !payer.Equals(msg.Signer) {
			return types.Context{}, errors.New("Fee Payer and MSG Signer are not the same ")
		}
		feeInRowan := feeTx.GetFee() // --fees 2000000rowan
		requiredRowan := feeInRowan.AmountOf(clpTypes.GetSettlementAsset().Symbol)
		coinsBalance := r.ck.GetBankKeeper().GetCoins(ctx, payer)
		userRowan := coinsBalance.AmountOf(clpTypes.GetSettlementAsset().Symbol)
		fmt.Println("Balance Before : ", r.ck.GetBankKeeper().GetCoins(ctx, payer))

		payerHasRowan := true
		if userRowan.LT(requiredRowan) {
			requiredRowan = requiredRowan.Sub(userRowan)
			payerHasRowan = false
			fmt.Printf("User Does not have enough rowan , trying to swap  :%s rowan ", requiredRowan.String())
		}
		if !payerHasRowan {
			err = EnrichPayerWithRowan(r.ck, ctx, msg, requiredRowan)
			if err != nil {
				fmt.Println("Error :", err)
				return types.Context{}, err
			}
			ctx.Logger().Info(fmt.Sprintf("Enriched user %s with %s rowan : ", payer.String(), requiredRowan.String()))
		}
		fmt.Println("Balance After : ", r.ck.GetBankKeeper().GetCoins(ctx, payer))
	default:
		return types.Context{}, errors.New("Unknown Swap type")
	}

	return next(ctx, tx, simulate)
}

func EnrichPayerWithRowan(ck keeper.Keeper, ctx types.Context, msg clpTypes.MsgSwap, requiredRowan types.Int) (err error) {
	pool, err := ck.GetPool(ctx, msg.ReceivedAsset.Symbol)
	if err != nil {
		return
	}
	ex := pool.ExternalAssetBalance
	na := pool.NativeAssetBalance
	//priceMultiplier := ex.Quo(na)
	//// Send to module account priceMultiplier * requiredRowan amount of cToken
	//// Send to user from module account requiredRowan of rowan
	//
	//
	//// subtract ceth
	//// add rowan
	//
	//// x - amount of cToken
	//// X - cToken Balance
	//// Y - Rowan Balance
	//// S - amount of Rowan
	//sendAmount, err := keeper.ReverseSwap(msg.ReceivedAsset.Symbol, pool.ExternalAssetBalance, pool.NativeAssetBalance, types.Uint(requiredRowan))
	//if err != nil {
	//	return
	//}
	//
	//swapMSG := clpTypes.MsgSwap{
	//	Signer:             msg.Signer,
	//	SentAsset:          msg.ReceivedAsset,
	//	ReceivedAsset:      clpTypes.GetSettlementAsset(),
	//	SentAmount:         sendAmount,
	//	MinReceivingAmount: types.NewUintFromBigInt(requiredRowan.BigInt()),
	//}
	//res, err := handleMsgSwap(ctx, ck, swapMSG)
	//fmt.Println("Swap Result :", res, err)
	if err != nil {
		return
	}
	return nil
}
