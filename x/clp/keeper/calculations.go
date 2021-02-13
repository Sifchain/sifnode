package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

//------------------------------------------------------------------------------------------------------------------
// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func SwapOne(from types.Asset, sentAmount sdk.Uint, to types.Asset, pool types.Pool) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {

	var X sdk.Uint
	var Y sdk.Uint

	if to == types.GetSettlementAsset() {
		Y = pool.NativeAssetBalance
		X = pool.ExternalAssetBalance
	} else {
		X = pool.NativeAssetBalance
		Y = pool.ExternalAssetBalance
	}
	x := sentAmount
	liquidityFee := calcLiquidityFee(X, x, Y)
	priceImpact := calcPriceImpact(X, x)
	swapResult := calcSwapResult(X, x, Y)
	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), types.Pool{}, types.ErrNotEnoughAssetTokens
	}
	if from == types.GetSettlementAsset() {
		pool.NativeAssetBalance = X.Add(x)
		pool.ExternalAssetBalance = Y.Sub(swapResult)
	} else {
		pool.ExternalAssetBalance = X.Add(x)
		pool.NativeAssetBalance = Y.Sub(swapResult)
	}

	return swapResult, liquidityFee, priceImpact, pool, nil
}

// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func CalculateWithdrawal(poolUnits sdk.Uint, nativeAssetBalance string,
	externalAssetBalance string, lpUnits string, wBasisPoints string, asymmetry sdk.Int) (sdk.Uint, sdk.Uint, sdk.Uint, sdk.Uint) {
	poolUnitsF := sdk.NewDecFromBigInt(poolUnits.BigInt())

	nativeAssetBalanceF, err := sdk.NewDecFromStr(nativeAssetBalance)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetBalance, err))
	}
	externalAssetBalanceF, err := sdk.NewDecFromStr(externalAssetBalance)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetBalance, err))
	}
	lpUnitsF, err := sdk.NewDecFromStr(lpUnits)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", lpUnits, err))
	}
	wBasisPointsF, err := sdk.NewDecFromStr(wBasisPoints)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", wBasisPoints, err))
	}
	asymmetryF, err := sdk.NewDecFromStr(asymmetry.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", asymmetry.String(), err))
	}
	denominator := sdk.NewDec(10000).Quo(wBasisPointsF)
	unitsToClaim := lpUnitsF.Quo(denominator)
	withdrawExternalAssetAmount := externalAssetBalanceF.Quo(poolUnitsF.Quo(unitsToClaim))
	withdrawNativeAssetAmount := nativeAssetBalanceF.Quo(poolUnitsF.Quo(unitsToClaim))

	swapAmount := sdk.NewDec(0)
	//if asymmetry is positive we need to swap from native to external
	if asymmetry.IsPositive() {
		unitsToSwap := unitsToClaim.Quo(sdk.NewDec(10000).Quo(asymmetryF.Abs()))
		swapAmount = nativeAssetBalanceF.Quo(poolUnitsF.Quo(unitsToSwap))
	}
	//if asymmetry is negative we need to swap from external to native
	if asymmetry.IsNegative() {
		unitsToSwap := unitsToClaim.Quo(sdk.NewDec(10000).Quo(asymmetryF.Abs()))
		swapAmount = externalAssetBalanceF.Quo(poolUnitsF.Quo(unitsToSwap))
	}
	//if asymmetry is 0 we don't need to swap

	lpUnitsLeft := lpUnitsF.Sub(unitsToClaim)
	return sdk.NewUintFromBigInt(withdrawNativeAssetAmount.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(withdrawExternalAssetAmount.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(lpUnitsLeft.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(swapAmount.RoundInt().BigInt())
}

// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md

//native asset balance  : currently in pool before adding
//external asset balance : currently in pool before adding
//native asset to added  : the amount the user sends
//external asset amount to be added : the amount the user sends

// R = native Balance (before)
// A = external Balance (before)
// r = native asset added;
// a = external asset added
// P = existing Pool Units
// slipAdjustment = (1 - ABS((R a - r A)/((2 r + R) (a + A))))
// units = ((P (a R + A r))/(2 A R))*slidAdjustment

func GetLen(str string) int64 {
	return int64(len(str))
}

func ValidatePoolUnit(oldPoolUnits, nativeAssetBalance, externalAssetBalance,
	nativeAssetAmount, externalAssetAmount sdk.Uint) bool {

	minValue := sdk.NewUintFromString("1000000000")
	// No token is added
	if nativeAssetAmount.IsZero() && externalAssetAmount.IsZero() {
		return false
	}
	// Check all values are within range
	if !oldPoolUnits.IsZero() && oldPoolUnits.LT(minValue) {
		return false
	}
	if !nativeAssetBalance.IsZero() && nativeAssetBalance.LT(minValue) {
		return false
	}
	if !externalAssetBalance.IsZero() && externalAssetBalance.LT(minValue) {
		return false
	}
	if !nativeAssetAmount.IsZero() && nativeAssetAmount.LT(minValue) {
		return false
	}
	if !externalAssetAmount.IsZero() && externalAssetAmount.LT(minValue) {
		return false
	}
	return true
}
func CalculatePoolUnits(oldPoolUnits, nativeAssetBalance, externalAssetBalance,
	nativeAssetAmount, externalAssetAmount sdk.Uint) (sdk.Uint, sdk.Uint, error) {
	// refactor this to use ValidInputs
	if !ValidatePoolUnit(oldPoolUnits, nativeAssetBalance, externalAssetBalance,
		nativeAssetAmount, externalAssetAmount) {
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAmountTooLow
	}

	minLen := GetLen(oldPoolUnits.String())
	if GetLen(nativeAssetAmount.String()) < minLen {
		minLen = GetLen(nativeAssetAmount.String())
	}
	if GetLen(externalAssetAmount.String()) < minLen {
		minLen = GetLen(externalAssetAmount.String())
	}
	if GetLen(nativeAssetAmount.String()) < minLen {
		minLen = GetLen(nativeAssetAmount.String())
	}
	if GetLen(externalAssetAmount.String()) < minLen {
		minLen = GetLen(externalAssetAmount.String())
	}

	if nativeAssetBalance.Add(nativeAssetAmount).IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), errors.Wrap(errors.ErrInsufficientFunds, nativeAssetAmount.String())
	}
	if externalAssetBalance.Add(externalAssetAmount).IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), errors.Wrap(errors.ErrInsufficientFunds, externalAssetAmount.String())
	}
	if nativeAssetBalance.IsZero() || externalAssetBalance.IsZero() {
		return nativeAssetAmount, nativeAssetAmount, nil
	}
	P, err := sdk.NewDecFromStr(oldPoolUnits.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", oldPoolUnits.String(), err))
	}
	R, err := sdk.NewDecFromStr(nativeAssetBalance.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetBalance.String(), err))
	}
	A, err := sdk.NewDecFromStr(externalAssetBalance.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetBalance.String(), err))
	}
	r, err := sdk.NewDecFromStr(nativeAssetAmount.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetAmount.String(), err))
	}
	a, err := sdk.NewDecFromStr(externalAssetAmount.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetAmount.String(), err))
	}

	slipAdjDenominator := (r.MulInt64(2).Add(R)).Mul(a.Add(A))
	// ABS((R a - r A)/((2 r + R) (a + A)))
	var slipAdjustment sdk.Dec
	if R.Mul(a).GT(r.Mul(A)) {
		slipAdjustment = R.Mul(a).Sub(r.Mul(A)).Quo(slipAdjDenominator)
	} else {
		slipAdjustment = r.Mul(A).Sub(R.Mul(a)).Quo(slipAdjDenominator)
	}
	// (1 - ABS((R a - r A)/((2 r + R) (a + A))))
	slipAdjustment = sdk.NewDec(1).Sub(slipAdjustment)

	P = ReducePrecision(P, minLen)
	R = ReducePrecision(R, minLen)
	A = ReducePrecision(A, minLen)
	a = ReducePrecision(a, minLen)
	r = ReducePrecision(r, minLen)

	// ((P (a R + A r))

	numerator := P.Mul(a.Mul(R).Add(A.Mul(r)))
	// 2AR
	denominator := sdk.NewDec(2).Mul(A).Mul(R)
	stakeUnits := numerator.Quo(denominator).Mul(slipAdjustment)
	P = IncreasePrecision(P, minLen)
	newPoolUnit := P.Add(stakeUnits)
	newPoolUnit = IncreasePrecision(newPoolUnit, minLen)
	stakeUnits = IncreasePrecision(stakeUnits, minLen)
	return sdk.NewUintFromBigInt(newPoolUnit.RoundInt().BigInt()), sdk.NewUintFromBigInt(stakeUnits.RoundInt().BigInt()), nil
}

// Add validations for X,x,Y
//( x^2 * Y ) / ( x + X )^2
func calcLiquidityFee(X, x, Y sdk.Uint) sdk.Uint {
	// if inputs are outside range return error
	if !ValidateInputs([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint() //error
	}
	//if any input is 0 return 0
	if !ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint()
	}
	//if !ValidateInputs([]sdk.Uint{x,Y}){
	//	return sdk.ZeroUint() //error
	//}
	//if X.IsZero() && x.IsZero() {
	//	return sdk.ZeroUint()
	//}
	//n := x.Mul(x).Mul(Y)
	//d := x.Add(X)
	//de := d.Mul(d)
	//return n.Quo(de)
	d := x.Add(X)
	denom := d.Mul(d)
	return (x.Mul(x).Mul(Y)).Quo(denom)
}
func calcSwapResult(X, x, Y sdk.Uint) sdk.Uint {
	// if inputs are outside range return error
	if !ValidateInputs([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint()
	}
	//if any input is 0 return 0
	if !ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint()
	}
	d := x.Add(X)
	denom := d.Mul(d)
	return (x.Mul(X).Mul(Y)).Quo(denom)
}

//( x^2 * Y ) / ( x + X )^2

func calcPriceImpact(X, x sdk.Uint) sdk.Uint {
	// if inputs are outside range return error
	if !ValidateInputs([]sdk.Uint{X, x}) {
		return sdk.ZeroUint()
	}
	if (X.IsZero() && x.IsZero()) || x.IsZero() {
		return sdk.ZeroUint()
	}
	denom := x.Add(X)
	return x.Quo(denom)
}

func CalculateAllAssetsForLP(pool types.Pool, lp types.LiquidityProvider) (sdk.Uint, sdk.Uint, sdk.Uint, sdk.Uint) {
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	return CalculateWithdrawal(poolUnits, nativeAssetBalance.String(), externalAssetBalance.String(),
		lp.LiquidityProviderUnits.String(), sdk.NewInt(types.MaxWbasis).String(), sdk.ZeroInt())
}

func ValidateInputs(inputs []sdk.Uint) bool {
	minValue := sdk.NewUintFromString("1000000000")
	for _, val := range inputs {
		if !val.IsZero() && val.LT(minValue) {
			return false
		}
	}
	return true
}

func ValidateZero(inputs []sdk.Uint) bool {
	for _, val := range inputs {
		if val.IsZero() {
			return false
		}
	}
	return true
}

func ReducePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.Quo(p)
}

func IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.Mul(p)
}
