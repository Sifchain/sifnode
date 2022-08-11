package keeper

import (
	"fmt"
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
)

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

func CalculateExternalSwapAmountAsymmetricFloat(Y, X, y, x, f, r float64) float64 {

	return math.Abs((math.Sqrt(Y*(-1*(x+X))*(-1*f*f*x*Y-f*f*X*Y-2*f*r*x*Y+4*f*r*X*y+2*f*r*X*Y+4*f*X*y+4*f*X*Y-r*r*x*Y-r*r*X*Y-4*r*X*y-4*r*X*Y-4*X*y-4*X*Y)) + f*x*Y + f*X*Y + r*x*Y - 2*r*X*y - r*X*Y - 2*X*y - 2*X*Y) / (2 * (r + 1) * (y + Y)))
}

func CalculateSwapAmountAsymmetricBrokenDown(Y, X, y, x, f, r float64) float64 {
	a_ := x + X
	b_ := -1 * a_
	c_ := Y * b_

	d_ := f * f * x * Y
	e_ := f * f * X * Y
	f_ := 2 * f * r * x * Y
	g_ := 4 * f * r * X * y
	h_ := 2 * f * r * X * Y
	i_ := 4 * f * X * y
	j_ := 4 * f * X * Y
	k_ := r * r * x * Y
	l_ := r * r * X * Y
	m_ := 4 * r * X * y
	n_ := 4 * r * X * Y
	o_ := 4 * X * y
	p_ := 4 * X * Y
	q_ := f * x * Y
	r_ := f * X * Y
	s_ := r * x * Y
	t_ := 2 * r * X * y
	u_ := r * X * Y
	v_ := 2 * X * y
	w_ := 2 * X * Y

	r1 := (r + 1)
	x_ := (y + Y)
	y_ := -d_ - e_ - f_ + g_ + h_ + i_ + j_ - k_ - l_ - m_ - n_ - o_ - p_
	z_ := c_ * y_
	aa_ := math.Sqrt(z_)
	ab_ := (aa_ + q_ + r_ + s_ - t_ - u_ - v_ - w_)
	ac_ := (2 * r1 * x_)
	ad_ := ab_ / ac_

	return math.Abs(ad_)
}

// Calculates how much external asset to swap for an asymmetric add
func CalculateExternalSwapAmountAsymmetric(Y, X, y, x, f, r *big.Rat) big.Rat {

	var a_, b_, c_, d_, e_, f_, g_, h_, i_, j_, k_, l_, m_, n_, o_, p_, q_, r_, s_, t_, u_, v_, w_, x_, y_, z_, aa_, ab_, ac_, ad_, minusOne, one, two, four, r1 big.Rat
	minusOne.SetInt64(-1)
	one.SetInt64(1)
	two.SetInt64(2)
	four.SetInt64(4)
	r1.Add(r, &one)

	a_.Add(x, X)           // a_ = x + X
	b_.Mul(&a_, &minusOne) // b_ = -1 * (x + X)
	c_.Mul(Y, &b_)         // c_ = Y * -1 * (x + X)

	d_.Mul(f, f).Mul(&d_, x).Mul(&d_, Y)                 // d_ = f * f * x * Y
	e_.Mul(f, f).Mul(&e_, X).Mul(&e_, Y)                 // e_ := f * f * X * Y
	f_.Mul(&two, f).Mul(&f_, r).Mul(&f_, x).Mul(&f_, Y)  // f_ := 2 * f * r * x * Y
	g_.Mul(&four, f).Mul(&g_, r).Mul(&g_, X).Mul(&g_, y) // g_ := 4 * f * r * X * y
	h_.Mul(&two, f).Mul(&h_, r).Mul(&h_, X).Mul(&h_, Y)  // h_ := 2 * f * r * X * Y
	i_.Mul(&four, f).Mul(&i_, X).Mul(&i_, y)             // i_ := 4 * f * X * y
	j_.Mul(&four, f).Mul(&j_, X).Mul(&j_, Y)             // j_ := 4 * f * X * Y
	k_.Mul(r, r).Mul(&k_, x).Mul(&k_, Y)                 // k_ := r * r * x * Y
	l_.Mul(r, r).Mul(&l_, X).Mul(&l_, Y)                 // l_ := r * r * X * Y
	m_.Mul(&four, r).Mul(&m_, X).Mul(&m_, y)             // m_ := 4 * r * X * y
	n_.Mul(&four, r).Mul(&n_, X).Mul(&n_, Y)             // n_ := 4 * r * X * Y
	o_.Mul(&four, X).Mul(&o_, y)                         // o_ := 4 * X * y
	p_.Mul(&four, X).Mul(&p_, Y)                         // p_ := 4 * X * Y
	q_.Mul(f, x).Mul(&q_, Y)                             // q_ := f * x * Y
	r_.Mul(f, X).Mul(&r_, Y)                             // r_ := f * X * Y
	s_.Mul(r, x).Mul(&s_, Y)                             // s_ := r * x * Y
	t_.Mul(&two, r).Mul(&t_, X).Mul(&t_, y)              // t_ := 2 * r * X * y
	u_.Mul(r, X).Mul(&u_, Y)                             // u_ := r * X * Y
	v_.Mul(&two, X).Mul(&v_, y)                          // v_ := 2 * X * y
	w_.Mul(&two, X).Mul(&w_, Y)                          // w_ := 2 * X * Y

	x_.Add(y, Y) // x_ := (y + Y)

	y_.Add(&g_, &h_).Add(&y_, &i_).Add(&y_, &j_).Sub(&y_, &d_).Sub(&y_, &e_).Sub(&y_, &f_).Sub(&y_, &k_).Sub(&y_, &l_).Sub(&y_, &m_).Sub(&y_, &n_).Sub(&y_, &o_).Sub(&y_, &p_) // y_ :=  g_ + h_ + i_ + j_ - d_ - e_ - f_ - k_ - l_ - m_ - n_ - o_ - p_ // y_ := -d_ - e_ - f_ + g_ + h_ + i_ + j_ - k_ - l_ - m_ - n_ - o_ - p_

	z_.Mul(&c_, &y_) // z_ := c_ * y_

	aa_.SetInt(ApproxRatSquareRoot(&z_)) // aa_ := math.Sqrt(z_)

	ab_.Add(&aa_, &q_).Add(&ab_, &r_).Add(&ab_, &s_).Sub(&ab_, &t_).Sub(&ab_, &u_).Sub(&ab_, &v_).Sub(&ab_, &w_) // ab_ := (aa_ + q_ + r_ + s_ - t_ - u_ - v_ - w_)

	ac_.Mul(&two, &r1).Mul(&ac_, &x_) // ac_ := (2 * r1 * x_)

	ad_.Quo(&ab_, &ac_) // ad_ := ab_ / ac_

	return *ad_.Abs(&ad_)
}
