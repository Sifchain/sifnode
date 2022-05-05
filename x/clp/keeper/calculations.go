package keeper

import (
	"fmt"
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
)

//------------------------------------------------------------------------------------------------------------------
// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func SwapOne(from types.Asset,
	sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {

	X, Y, toRowan := pool.ExtractValues(to)

	liquidityFee := CalcLiquidityFee(X, sentAmount, Y)
	priceImpact := calcPriceImpact(X, sentAmount)
	swapResult := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate)

	// NOTE: impossible... pre-pmtp at least
	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), types.Pool{}, types.ErrNotEnoughAssetTokens
	}

	pool.UpdateBalances(toRowan, X, sentAmount, Y, swapResult)

	return swapResult, liquidityFee, priceImpact, pool, nil
}

func CalcSwapPrice(from types.Asset,
	sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	normalizationFactor sdk.Dec,
	adjustExternalToken bool,
	pmtpCurrentRunningRate sdk.Dec) sdk.Dec {

	X, Y, toRowan := pool.ExtractValues(to)
	swapResult := CalcSwapPriceResult(toRowan, normalizationFactor, adjustExternalToken, X, sentAmount, Y, pmtpCurrentRunningRate)

	return swapResult
}

func CalcSwapPmtp(toRowan bool, y, pmtpCurrentRunningRate sdk.Dec) sdk.Dec {
	// if pmtpCurrentRunningRate.IsNil() {
	// 	if toRowan {
	// 		return y.Quo(sdk.NewDec(1))
	// 	}
	// 	return y.Mul(sdk.NewDec(1))
	// }
	if toRowan {
		return y.Quo(sdk.NewDec(1).Add(pmtpCurrentRunningRate))
	}
	return y.Mul(sdk.NewDec(1).Add(pmtpCurrentRunningRate))
}

func GetSwapFee(sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	pmtpCurrentRunningRate sdk.Dec) sdk.Uint {
	X, Y, toRowan := pool.ExtractValues(to)
	swapResult := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate)

	if swapResult.GTE(Y) {
		return sdk.ZeroUint()
	}
	return swapResult
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
func CalculateWithdrawalFromUnits(poolUnits sdk.Uint, nativeAssetBalance string,
	externalAssetBalance string, lpUnits string, withdrawUnits sdk.Uint) (sdk.Uint, sdk.Uint, sdk.Uint) {
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
	withdrawUnitsF, err := sdk.NewDecFromStr(withdrawUnits.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", withdrawUnits, err))
	}

	withdrawExternalAssetAmount := externalAssetBalanceF.Quo(poolUnitsF.Quo(withdrawUnitsF))
	withdrawNativeAssetAmount := nativeAssetBalanceF.Quo(poolUnitsF.Quo(withdrawUnitsF))

	//if asymmetry is 0 we don't need to swap
	lpUnitsLeft := lpUnitsF.Sub(withdrawUnitsF)

	return sdk.NewUintFromBigInt(withdrawNativeAssetAmount.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(withdrawExternalAssetAmount.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(lpUnitsLeft.RoundInt().BigInt())
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
// slipAdjustment = (1 - ABS((R a - r A)/((r + R) (a + A))))
// units = ((P (a R + A r))/(2 A R))*slidAdjustment

func CalculatePoolUnits(oldPoolUnits, nativeAssetBalance, externalAssetBalance, nativeAssetAmount,
	externalAssetAmount sdk.Uint, symmetryThreshold sdk.Dec) (sdk.Uint, sdk.Uint, error) {

	if nativeAssetAmount.IsZero() && externalAssetAmount.IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAmountTooLow
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

	slipAdjustment, RTimesa, rTimesA := calculateSlipAdjustment(nativeAssetBalance.BigInt(), externalAssetBalance.BigInt(),
		nativeAssetAmount.BigInt(), externalAssetAmount.BigInt())

	one := big.NewRat(1, 1)
	one.Sub(one, &slipAdjustment)

	symmetryThresholdRat := decToRat(&symmetryThreshold)
	if one.Cmp(&symmetryThresholdRat) == 1 { // this is: if a > b
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAsymmetricAdd
	}

	stakeUnits := calculateStakeUnits(oldPoolUnits.BigInt(), nativeAssetBalance.BigInt(),
		externalAssetBalance.BigInt(), nativeAssetAmount.BigInt(), &RTimesa, &rTimesA, &slipAdjustment)

	var newPoolUnit big.Int
	newPoolUnit.Add(oldPoolUnits.BigInt(), &stakeUnits)

	return sdk.NewUintFromBigInt(&newPoolUnit), sdk.NewUintFromBigInt(&stakeUnits), nil
}

// units = ((P (a R + A r))/(2 A R))*slidAdjustment
func calculateStakeUnits(P, R, A, r, RTimesa, rTimesA *big.Int, slipAdslipAdjustment *big.Rat) big.Int {
	var add, numerator big.Int
	add.Add(RTimesa, rTimesA)
	numerator.Mul(P, &add)

	var denominator big.Int
	denominator.Mul(big.NewInt(2), A)
	denominator.Mul(&denominator, R)

	var n, d, stakeUnits big.Rat
	n.SetInt(&numerator)
	d.SetInt(&denominator)
	stakeUnits.Quo(&n, &d)
	stakeUnits.Mul(&stakeUnits, slipAdslipAdjustment)

	num := stakeUnits.Num()
	denom := stakeUnits.Denom()
	num.Quo(num, denom)

	return *num
}

// slipAdjustment = (1 - ABS((R a - r A)/((r + R) (a + A))))
func calculateSlipAdjustment(R, A, r, a *big.Int) (big.Rat, big.Int, big.Int) {
	var denominator, rPlusR, aPlusA big.Int
	rPlusR.Add(r, R)
	aPlusA.Add(a, A)
	denominator.Mul(&rPlusR, &aPlusA)

	var RTimesa, rTimesA, nominator big.Int
	RTimesa.Mul(R, a)
	rTimesA.Mul(r, A)
	nominator.Sub(&RTimesa, &rTimesA)

	var one, nom, denom, slipAdjustment big.Rat
	one.SetInt64(1)

	nom.SetInt(&nominator)
	denom.SetInt(&denominator)

	slipAdjustment.Quo(&nom, &denom)
	slipAdjustment.Abs(&slipAdjustment)
	slipAdjustment.Sub(&one, &slipAdjustment)

	return slipAdjustment, RTimesa, rTimesA
}

func CalcLiquidityFee(X, x, Y sdk.Uint) sdk.Uint {
	if IsAnyZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint()
	}

	Xb := X.BigInt()
	xb := x.BigInt()
	Yb := Y.BigInt()

	var sq, n, s, d, fee big.Int

	sq.Mul(xb, xb)  // sq = x**2
	n.Mul(&sq, Yb)  // n = x**2 * Y
	s.Add(Xb, xb)   // s = x + X
	d.Mul(&s, &s)   // d = (x + X)**2
	fee.Quo(&n, &d) // fee = n / d = (x**2 * Y) / (x + X)**2

	//TODO: can this panic? Does it matter?
	return sdk.NewUintFromBigInt(&fee)
}

func CalcSwapResult(toRowan bool,
	X, x, Y sdk.Uint,
	pmtpCurrentRunningRate sdk.Dec) sdk.Uint {

	if IsAnyZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint()
	}

	tmp := calcSwap(x.BigInt(), X.BigInt(), Y.BigInt())
	y := sdk.NewDecFromBigInt(&tmp)
	pmtpFac := CalcPmtpFactor(pmtpCurrentRunningRate)

	var res sdk.Dec
	if toRowan {
		res = y.Quo(pmtpFac) // res = y / pmtpFac
	} else {
		res = y.Mul(pmtpFac) // res = y * pmtpFac
	}

	return sdk.NewUintFromBigInt(res.RoundInt().BigInt())
}

func calcSwap(x, X, Y *big.Int) big.Int {
	var s, d, d2, d3, y big.Int

	s.Add(X, x)    // s = X + x
	d.Mul(&s, &s)  // d = (X + x)**2
	d2.Mul(X, Y)   // d2 = X * Y
	d3.Mul(x, &d2) // d3 = x * X * Y
	y.Quo(&d3, &d) // y = d3 / d = (x * X * Y) / (X + x)**2

	return y
}

func CalcPmtpFactor(r sdk.Dec) sdk.Dec {
	return sdk.NewDec(1).Add(r)
}

func CalcSwapPriceResult(toRowan bool,
	normalizationFactor sdk.Dec,
	adjustExternalToken bool,
	X, x, Y sdk.Uint,
	pmtpCurrentRunningRate sdk.Dec) sdk.Dec {
	if !ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroDec()
	}

	nf := sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt())
	if adjustExternalToken {
		if toRowan {
			X = X.Mul(nf)
			x = x.Mul(nf)
		} else {
			Y = Y.Mul(nf)
		}
	} else {
		if toRowan {
			Y = Y.Mul(nf)
		} else {
			X = X.Mul(nf)
			x = x.Mul(nf)
		}
	}

	minLen := GetMinLen([]sdk.Uint{X, x, Y})
	Xd := ReducePrecision(sdk.NewDecFromBigInt(X.BigInt()), minLen)
	xd := ReducePrecision(sdk.NewDecFromBigInt(x.BigInt()), minLen)
	Yd := ReducePrecision(sdk.NewDecFromBigInt(Y.BigInt()), minLen)

	s := xd.Add(Xd)
	d := s.Mul(s)
	y := xd.Mul(Xd).Mul(Yd).Quo(d)
	y = IncreasePrecision(y, minLen)
	// we're looking for price in absolute units here
	if toRowan {
		y = y.Quo(normalizationFactor)
	}
	y = CalcSwapPmtp(toRowan, y, pmtpCurrentRunningRate)
	return y
}

func calcPriceImpact(X, x sdk.Uint) sdk.Uint {
	if x.IsZero() {
		return sdk.ZeroUint()
	}
	d := x.Add(X)

	return x.Quo(d)
}

func CalculateAllAssetsForLP(pool types.Pool, lp types.LiquidityProvider) (sdk.Uint, sdk.Uint, sdk.Uint, sdk.Uint) {
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	return CalculateWithdrawal(
		poolUnits,
		nativeAssetBalance.String(),
		externalAssetBalance.String(),
		lp.LiquidityProviderUnits.String(),
		sdk.NewInt(types.MaxWbasis).String(),
		sdk.ZeroInt(),
	)
}

func decToRat(d *sdk.Dec) big.Rat {
	var rat big.Rat

	rat.SetInt(d.BigInt())
	decimals := int64(math.Pow10(sdk.Precision)) // 10**18
	denom := big.NewRat(decimals, 1)
	rat.Quo(&rat, denom)

	return rat
}
