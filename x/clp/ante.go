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
		// Fee Payer is always the msg signer for swap transactions
		if !payer.Equals(msg.Signer) {
			return types.Context{}, errors.New("Fee Payer and MSG Signer are not the same ")
		}
		// User is sending rowan , that means they already have rowan balance
		if msg.SentAsset.Equals(GetSettlementAsset()) {
			return types.Context{}, nil
		}
		// The amount of Fee user agreed to pay
		feeInRowan := feeTx.GetFee()
		requiredRowan := feeInRowan.AmountOf(clpTypes.GetSettlementAsset().Symbol)
		coinsBalance := r.ck.GetBankKeeper().GetCoins(ctx, payer)
		userRowan := coinsBalance.AmountOf(clpTypes.GetSettlementAsset().Symbol)

		// Check if user does not have enough
		if userRowan.LT(requiredRowan) {
			// Add the remaining amount of rowan to the users balance
			requiredRowan = requiredRowan.Sub(userRowan)
			err = EnrichPayerWithRowan(r.ck, ctx, msg, requiredRowan)
			if err != nil {
				return types.Context{}, errors.Wrap(clpTypes.ErrUnableEnrichUser, err.Error())
			}
		}

	default:
		return types.Context{}, errors.New("Unknown Swap type") // Unreachable code
	}

	return next(ctx, tx, simulate)
}

func EnrichPayerWithRowan(ck keeper.Keeper, ctx types.Context, msg clpTypes.MsgSwap, requiredRowan types.Int) (err error) {
	// Derive rowan price from sent asset pool
	pool, err := ck.GetPool(ctx, msg.SentAsset.Symbol)
	if err != nil {
		return
	}
	ex := pool.ExternalAssetBalance
	na := pool.NativeAssetBalance
	// Derive price of cToken relative to rowan
	priceMultiplier := types.NewIntFromBigInt((ex.Quo(na).Mul(types.NewUintFromString(TxFeeMultiplier))).BigInt())
	cTokenSendCoin := types.NewCoins(types.NewCoin(msg.SentAsset.Symbol, priceMultiplier.Mul(requiredRowan)))
	rowanReceiveCoin := types.NewCoins(types.NewCoin(GetSettlementAsset().Symbol, requiredRowan))
	// Send to module first to avoid deficit
	err = ck.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, msg.Signer, ModuleName, cTokenSendCoin)
	if err != nil {
		return
	}
	err = ck.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, msg.Signer, rowanReceiveCoin)
	if err != nil {
		return
	}
	ctx.Logger().Info(fmt.Sprintf("\nEnriched user %s | Swapped %s for %s : ", msg.Signer.String(), cTokenSendCoin.String(), rowanReceiveCoin.String()))
	return nil
}
