package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Exit if the received asset amount is not above the minimum
func (k Keeper) CheckSentAmount(sentAmount sdk.Uint, minReceivingAmount sdk.Uint, to types.Asset) error {
	if sentAmount.LT(minReceivingAmount) {
		return types.ErrSwapAmountTooSmall
	}
	return nil
}

// Exit if the pool balance is not above the threshold
func (k Keeper) CheckPoolHealth(ctx sdk.Context, pool types.Pool) error {
	if !k.IsPoolHealthy(ctx, pool) {
		return types.ErrPoolNotHealthy
	}
	return nil
}

// Exit if the swap fee is not the expected one
func (k Keeper) CheckSwapFee(ctx sdk.Context, pool types.Pool, to types.Asset, marginEnabled bool, expectedSwapFee sdk.Dec) error {
	from, _ := pool.GetPoolAsset(to)
	swapFeeRate := k.GetSwapFeeRate(ctx, from, marginEnabled)
	if swapFeeRate.Abs().Sub(expectedSwapFee).GTE(sdk.MustNewDecFromStr("0.000000000000000001")) {
		return fmt.Errorf("swap fee is not the expected one, got %s, expected %s", swapFeeRate.String(), expectedSwapFee.String())
	}
	return nil
}

// Price impact affects the final price of the asset, good to have a limit
func (k Keeper) CheckPriceImpact(ctx sdk.Context, pool types.Pool, to types.Asset, sentAmount sdk.Uint) error {
	_, Y, toRowan := pool.ExtractValues(to)
	X, _ := pool.ExtractDebt(Y, Y, toRowan)
	priceImpact := CalcPriceImpact(X, sentAmount)
	// TODO: should be a parameter
	if priceImpact.GTE(sdk.MustNewDecFromStr("0.1")) {
		return fmt.Errorf("price impact is too high, got %s", priceImpact.String())
	}
	return nil
}
