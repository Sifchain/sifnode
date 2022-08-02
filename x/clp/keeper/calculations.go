package keeper

import (
	"fmt"
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
	externalAssetAmount sdk.Uint, externalDecimals uint8, symmetryThreshold, ratioThreshold sdk.Dec) (sdk.Uint, sdk.Uint, error) {

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

	slipAdjustmentValues := calculateSlipAdjustment(nativeAssetBalance.BigInt(), externalAssetBalance.BigInt(),
		nativeAssetAmount.BigInt(), externalAssetAmount.BigInt())

	one := big.NewRat(1, 1)
	symmetryThresholdRat := DecToRat(&symmetryThreshold)

	var diff big.Rat
	diff.Sub(one, slipAdjustmentValues.slipAdjustment)
	if diff.Cmp(&symmetryThresholdRat) == 1 { // this is: if diff > symmetryThresholdRat
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAsymmetricAdd
	}

	ratioThresholdRat := DecToRat(&ratioThreshold)
	normalisingFactor := CalcDenomChangeMultiplier(externalDecimals, types.NativeAssetDecimals)
	ratioThresholdRat.Mul(&ratioThresholdRat, &normalisingFactor)
	ratioDiff, err := CalculateRatioDiff(externalAssetBalance.BigInt(), nativeAssetBalance.BigInt(), externalAssetAmount.BigInt(), nativeAssetAmount.BigInt())
	if err != nil {
		return sdk.ZeroUint(), sdk.ZeroUint(), err
	}
	if ratioDiff.Cmp(&ratioThresholdRat) == 1 { //if ratioDiff > ratioThreshold
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAsymmetricRatioAdd
	}

	stakeUnits := calculateStakeUnits(oldPoolUnits.BigInt(), nativeAssetBalance.BigInt(),
		externalAssetBalance.BigInt(), nativeAssetAmount.BigInt(), slipAdjustmentValues)

	var newPoolUnit big.Int
	newPoolUnit.Add(oldPoolUnits.BigInt(), stakeUnits)

	return sdk.NewUintFromBigInt(&newPoolUnit), sdk.NewUintFromBigInt(stakeUnits), nil
}

// | A/R - a/r |
func CalculateRatioDiff(A, R, a, r *big.Int) (big.Rat, error) {
	if R.Cmp(big.NewInt(0)) == 0 || r.Cmp(big.NewInt(0)) == 0 { // check for zeros
		return *big.NewRat(0, 1), types.ErrAsymmetricRatioAdd
	}
	var AdivR, adivr, diff big.Rat

	AdivR.SetFrac(A, R)
	adivr.SetFrac(a, r)
	diff.Sub(&AdivR, &adivr)
	diff.Abs(&diff)

	return diff, nil
}

// units = ((P (a R + A r))/(2 A R))*slidAdjustment
func calculateStakeUnits(P, R, A, r *big.Int, slipAdjustmentValues *slipAdjustmentValues) *big.Int {
	var add, numerator big.Int
	add.Add(slipAdjustmentValues.RTimesa, slipAdjustmentValues.rTimesA)
	numerator.Mul(P, &add)

	var denominator big.Int
	denominator.Mul(big.NewInt(2), A)
	denominator.Mul(&denominator, R)

	var n, d, stakeUnits big.Rat
	n.SetInt(&numerator)
	d.SetInt(&denominator)
	stakeUnits.Quo(&n, &d)
	stakeUnits.Mul(&stakeUnits, slipAdjustmentValues.slipAdjustment)

	return RatIntQuo(&stakeUnits)
}

// slipAdjustment = (1 - ABS((R a - r A)/((r + R) (a + A))))
type slipAdjustmentValues struct {
	slipAdjustment *big.Rat
	RTimesa        *big.Int
	rTimesA        *big.Int
}

func calculateSlipAdjustment(R, A, r, a *big.Int) *slipAdjustmentValues {
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

	return &slipAdjustmentValues{slipAdjustment: &slipAdjustment, RTimesa: &RTimesa, rTimesA: &rTimesA}
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

	return sdk.NewUintFromBigInt(&fee)
}

func CalcSwapResult(toRowan bool,
	X, x, Y sdk.Uint,
	pmtpCurrentRunningRate sdk.Dec) sdk.Uint {

	if IsAnyZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint()
	}

	y := calcSwap(x.BigInt(), X.BigInt(), Y.BigInt())
	pmtpFac := calcPmtpFactor(pmtpCurrentRunningRate)

	var res big.Rat
	if toRowan {
		res.Quo(&y, &pmtpFac) // res = y / pmtpFac
	} else {
		res.Mul(&y, &pmtpFac) // res = y * pmtpFac
	}

	num := RatIntQuo(&res)
	return sdk.NewUintFromBigInt(num)
}

func calcSwap(x, X, Y *big.Int) big.Rat {
	var s, d, d2, d3 big.Int
	var numerator, denominator, y big.Rat

	s.Add(X, x)    // s = X + x
	d.Mul(&s, &s)  // d = (X + x)**2
	d2.Mul(X, Y)   // d2 = X * Y
	d3.Mul(x, &d2) // d3 = x * X * Y

	denominator.SetInt(&d)
	numerator.SetInt(&d3)
	y.Quo(&numerator, &denominator) // y = d3 / d = (x * X * Y) / (X + x)**2

	return y
}

func calcPmtpFactor(r sdk.Dec) big.Rat {
	rRat := DecToRat(&r)
	one := big.NewRat(1, 1)

	one.Add(one, &rRat)

	return *one
}

func CalcSpotPriceNative(pool *types.Pool, decimalsExternal uint8, pmtpCurrentRunningRate sdk.Dec) (sdk.Dec, error) {
	return CalcSpotPriceX(pool.NativeAssetBalance, pool.ExternalAssetBalance, types.NativeAssetDecimals, decimalsExternal, pmtpCurrentRunningRate, true)
}

func CalcSpotPriceExternal(pool *types.Pool, decimalsExternal uint8, pmtpCurrentRunningRate sdk.Dec) (sdk.Dec, error) {
	return CalcSpotPriceX(pool.ExternalAssetBalance, pool.NativeAssetBalance, decimalsExternal, types.NativeAssetDecimals, pmtpCurrentRunningRate, false)
}

// Calculates the spot price of asset X in the preferred denominations accounting for PMTP.
// Since this method applies PMTP adjustment, one of X, Y must be the native asset.
func CalcSpotPriceX(X, Y sdk.Uint, decimalsX, decimalsY uint8, pmtpCurrentRunningRate sdk.Dec, isXNative bool) (sdk.Dec, error) {
	if X.Equal(sdk.ZeroUint()) {
		return sdk.ZeroDec(), types.ErrInValidAmount
	}

	var price big.Rat
	price.SetFrac(Y.BigInt(), X.BigInt())

	pmtpFac := calcPmtpFactor(pmtpCurrentRunningRate)
	var pmtpPrice big.Rat
	if isXNative {
		pmtpPrice.Mul(&price, &pmtpFac) // pmtpPrice = price * pmtpFac
	} else {
		pmtpPrice.Quo(&price, &pmtpFac) // pmtpPrice = price / pmtpFac
	}

	dcm := CalcDenomChangeMultiplier(decimalsX, decimalsY)
	pmtpPrice.Mul(&pmtpPrice, &dcm)

	res := RatToDec(&pmtpPrice)
	return res, nil
}
func CalcRowanValue(pool *types.Pool, pmtpCurrentRunningRate sdk.Dec, rowanAmount sdk.Uint) (sdk.Uint, error) {
	spotPrice, err := CalcRowanSpotPrice(pool, pmtpCurrentRunningRate)
	if err != nil {
		return sdk.ZeroUint(), err
	}
	value := spotPrice.Mul(sdk.NewDecFromBigInt(rowanAmount.BigInt()))
	return sdk.NewUintFromBigInt(value.RoundInt().BigInt()), nil
}

// Calculates spot price of Rowan accounting for PMTP
func CalcRowanSpotPrice(pool *types.Pool, pmtpCurrentRunningRate sdk.Dec) (sdk.Dec, error) {
	rowanBalance := sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt())
	if rowanBalance.Equal(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrInValidAmount
	}
	externalAssetBalance := sdk.NewDecFromBigInt(pool.ExternalAssetBalance.BigInt())
	unadjusted := externalAssetBalance.Quo(rowanBalance)
	return unadjusted.Mul(pmtpCurrentRunningRate.Add(sdk.OneDec())), nil
}

// Denom change multiplier = 10**decimalsX / 10**decimalsY
func CalcDenomChangeMultiplier(decimalsX, decimalsY uint8) big.Rat {
	diff := Abs(int16(decimalsX) - int16(decimalsY))
	dec := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(diff)), nil) // 10**|decimalsX - decimalsY|

	var res big.Rat
	if decimalsX > decimalsY {
		return *res.SetInt(dec)
	}
	return *res.SetFrac(big.NewInt(1), dec)
}

func calcPriceImpact(X, x sdk.Uint) sdk.Uint {
	if x.IsZero() {
		return sdk.ZeroUint()
	}

	Xb := X.BigInt()
	xb := x.BigInt()

	var d, q big.Int
	d.Add(xb, Xb)
	q.Quo(xb, &d) // q = x / (x + X)

	return sdk.NewUintFromBigInt(&q)
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
