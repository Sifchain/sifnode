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
		feeInRowan := feeTx.GetFee()
		requiredRowan := feeInRowan.AmountOf(clpTypes.GetSettlementAsset().Symbol)
		coinsBalance := r.ck.GetBankKeeper().GetCoins(ctx, payer)
		userRowan := coinsBalance.AmountOf(clpTypes.GetSettlementAsset().Symbol)
		payerHasRowan := true
		if userRowan.LT(requiredRowan) {
			requiredRowan = requiredRowan.Sub(userRowan)
			payerHasRowan = false
			ctx.Logger().Info(fmt.Sprintf("\nUser Does not have enough rowan | Trying to swap  :%s for %s rowan ", msg.SentAsset, requiredRowan.String()))
		}
		if !payerHasRowan {
			err = EnrichPayerWithRowan(r.ck, ctx, msg, requiredRowan)
			if err != nil {
				return types.Context{}, err
			}
			ctx.Logger().Info(fmt.Sprintf("\nEnriched user %s with %s rowan : ", payer.String(), requiredRowan.String()))
		}
	default:
		return types.Context{}, errors.New("Unknown Swap type")
	}

	return next(ctx, tx, simulate)
}

func EnrichPayerWithRowan(ck keeper.Keeper, ctx types.Context, msg clpTypes.MsgSwap, requiredRowan types.Int) (err error) {
	pool, err := ck.GetPool(ctx, msg.SentAsset.Symbol)
	if err != nil {
		return
	}
	ex := pool.ExternalAssetBalance
	na := pool.NativeAssetBalance
	priceMultiplier := types.NewIntFromBigInt(ex.Quo(na).BigInt())
	cTokenSendCoin := types.NewCoins(types.NewCoin(msg.SentAsset.Symbol, priceMultiplier.Mul(requiredRowan)))
	rowanReceiveCoin := types.NewCoins(types.NewCoin(GetSettlementAsset().Symbol, requiredRowan))
	err = ck.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Signer, ModuleName, cTokenSendCoin)
	if err != nil {
		return
	}
	err = ck.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Signer, rowanReceiveCoin)
	if err != nil {
		return
	}
	return nil
}
