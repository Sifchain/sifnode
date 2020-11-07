package clp

import (
	"errors"
	"github.com/Sifchain/sifnode/x/clp/types"
)

//------------------------------------------------------------------------------------------------------------------
// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func SwapOne(from Asset, sentAmount uint, to Asset, pool Pool) (uint, uint, uint, Pool, error) {

	var X uint
	var Y uint

	if to == GetSettlementAsset() {
		Y = pool.NativeAssetBalance
		X = pool.ExternalAssetBalance
	} else {
		X = pool.NativeAssetBalance
		Y = pool.ExternalAssetBalance
	}
	x := sentAmount
	liquidityFee := calcLiquidityFee(X, x, Y)
	tradeSlip := calcTradeSlip(X, x)
	swapResult := calcSwapResult(X, x, Y)
	if swapResult >= Y {
		return 0, 0, 0, Pool{}, types.ErrNotEnoughAssetTokens
	}
	if from == GetSettlementAsset() {
		pool.NativeAssetBalance = X + x
		pool.ExternalAssetBalance = Y - swapResult
	} else {
		pool.ExternalAssetBalance = X + x
		pool.NativeAssetBalance = Y - swapResult
	}

	return swapResult, liquidityFee, tradeSlip, pool, nil
}

// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func CalculateWithdrawal(poolUnits uint, nativeAssetBalance uint,
	externalAssetBalance uint, lpUnits uint, wBasisPoints int, asymmetry int) (uint, uint, uint, uint) {
	poolUnitsF := float64(poolUnits)
	nativeAssetBalanceF := float64(nativeAssetBalance)
	externalAssetBalanceF := float64(externalAssetBalance)
	lpUnitsF := float64(lpUnits)
	wBasisPointsF := float64(wBasisPoints)
	asymmetryF := float64(asymmetry)

	unitsToClaim := lpUnitsF / (10000 / (wBasisPointsF))
	withdrawExternalAssetAmount := externalAssetBalanceF / (poolUnitsF / unitsToClaim)
	withdrawNativeAssetAmount := nativeAssetBalanceF / (poolUnitsF / unitsToClaim)

	swapAmount := 0.0
	//if asymmetry is positive we need to swap from native to external
	if asymmetry > 0 {
		unitsToSwap := (unitsToClaim) / (10000 / (asymmetryF))
		swapAmount = (nativeAssetBalanceF) / (poolUnitsF / unitsToSwap)
	}
	//if asymmetry is negative we need to swap from external to native
	if asymmetry < 0 {
		unitsToSwap := (unitsToClaim) / (10000 / (-1 * asymmetryF))
		swapAmount = (externalAssetBalanceF) / (poolUnitsF / unitsToSwap)
	}
	//if asymmetry is 0 we don't need to swap

	lpUnitsLeft := lpUnitsF - unitsToClaim
	if withdrawNativeAssetAmount < 0 {
		withdrawNativeAssetAmount = 0
	}
	if withdrawExternalAssetAmount < 0 {
		withdrawExternalAssetAmount = 0
	}
	if lpUnitsLeft < 0 {
		lpUnitsLeft = 0
	}
	if swapAmount < 0 {
		swapAmount = 0
	}

	return uint(withdrawNativeAssetAmount), uint(withdrawExternalAssetAmount), uint(lpUnitsLeft), uint(swapAmount)
}

// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md

//native asset balance  : currently in pool before adding
//external asset balance : currently in pool before adding
//native asset to added  : the amount the user sends
//external asset amount to be added : the amount the user sends

// r = native asset added;
// a = external asset added
// R = native Balance (before)
// A = external Balance (before)
// P = existing Pool Units
// slipAdjustment = (1 - ABS((R a - r A)/((2 r + R) (a + A))))
// units = ((P (a R + A r))/(2 A R))*slidAdjustment

func calculatePoolUnits(oldPoolUnits uint, nativeAssetBalance uint, externalAssetBalance uint,
	nativeAssetAmount uint, externalAssetAmount uint) (uint, uint, error) {
	if nativeAssetBalance+nativeAssetAmount == 0 {
		return 0, 0, errors.New("total Native in the pool is zero")
	}
	if externalAssetBalance+externalAssetAmount == 0 {
		return 0, 0, errors.New("total External in the pool is zero")
	}
	if nativeAssetBalance == 0 || externalAssetBalance == 0 {
		return nativeAssetAmount, externalAssetAmount, nil
	}
	P := float64(oldPoolUnits)
	R := float64(nativeAssetBalance)
	A := float64(externalAssetBalance)
	r := float64(nativeAssetAmount)
	a := float64(externalAssetAmount)

	// (2 r + R) (a + A)
	slipAdjDenominator := (2*r + R) * (a + A)
	// (R a - r A)/((2 r + R) (a + A))
	slipAd := (R*a - r*A) / slipAdjDenominator
	var slipAdjustment float64
	//ABS((R a - r A)/((2 r + R) (a + A)))
	if slipAd < 0 {
		slipAdjustment = -1.0 * slipAd
	}
	// (1 - ABS((R a - r A)/((2 r + R) (a + A))))
	slipAdjustment = 1 - slipAdjustment

	// ((P (a R + A r))
	numerator := P * (a*R + A*r)
	// 2AR
	denominator := 2 * A * R
	quotient := uint(numerator / denominator)
	lpUnits := quotient * uint(slipAdjustment)
	newPoolUnit := uint(P) + lpUnits
	return newPoolUnit, lpUnits, nil
}

func calcLiquidityFee(X, x, Y uint) uint {
	return (x * x * Y) / ((x + X) * (x + X))
}

func calcTradeSlip(X, x uint) uint {
	return x * (2*X + x) / (X * X)
}

func calcSwapResult(X, x, Y uint) uint {
	return (x * X * Y) / ((x + X) * (x + X))
}
