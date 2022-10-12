package keeper

import (
	"fmt"
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
func CalculateWithdrawal(poolUnits sdk.Uint, nativeAssetDepth string,
	externalAssetDepth string, lpUnits string, wBasisPoints string, asymmetry sdk.Int) (sdk.Uint, sdk.Uint, sdk.Uint, sdk.Uint) {
	poolUnitsF := sdk.NewDecFromBigInt(poolUnits.BigInt())

	nativeAssetDepthF, err := sdk.NewDecFromStr(nativeAssetDepth)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetDepth, err))
	}
	externalAssetDepthF, err := sdk.NewDecFromStr(externalAssetDepth)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetDepth, err))
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
	withdrawExternalAssetAmount := externalAssetDepthF.Quo(poolUnitsF.Quo(unitsToClaim))
	withdrawNativeAssetAmount := nativeAssetDepthF.Quo(poolUnitsF.Quo(unitsToClaim))

	swapAmount := sdk.NewDec(0)
	//if asymmetry is positive we need to swap from native to external
	if asymmetry.IsPositive() {
		unitsToSwap := unitsToClaim.Quo(sdk.NewDec(10000).Quo(asymmetryF.Abs()))
		swapAmount = nativeAssetDepthF.Quo(poolUnitsF.Quo(unitsToSwap))
	}
	//if asymmetry is negative we need to swap from external to native
	if asymmetry.IsNegative() {
		unitsToSwap := unitsToClaim.Quo(sdk.NewDec(10000).Quo(asymmetryF.Abs()))
		swapAmount = externalAssetDepthF.Quo(poolUnitsF.Quo(unitsToSwap))
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
func CalculateWithdrawalFromUnits(poolUnits sdk.Uint, nativeAssetDepth string,
	externalAssetDepth string, lpUnits string, withdrawUnits sdk.Uint) (sdk.Uint, sdk.Uint, sdk.Uint) {
	poolUnitsF := sdk.NewDecFromBigInt(poolUnits.BigInt())

	nativeAssetDepthF, err := sdk.NewDecFromStr(nativeAssetDepth)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetDepth, err))
	}
	externalAssetDepthF, err := sdk.NewDecFromStr(externalAssetDepth)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetDepth, err))
	}
	lpUnitsF, err := sdk.NewDecFromStr(lpUnits)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", lpUnits, err))
	}
	withdrawUnitsF, err := sdk.NewDecFromStr(withdrawUnits.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", withdrawUnits, err))
	}

	withdrawExternalAssetAmount := externalAssetDepthF.Quo(poolUnitsF.Quo(withdrawUnitsF))
	withdrawNativeAssetAmount := nativeAssetDepthF.Quo(poolUnitsF.Quo(withdrawUnitsF))

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

func CalculatePoolUnits(oldPoolUnits, nativeAssetDepth, externalAssetDepth, nativeAssetAmount,
	externalAssetAmount sdk.Uint, externalDecimals uint8, symmetryThreshold, ratioThreshold sdk.Dec) (sdk.Uint, sdk.Uint, error) {

	if nativeAssetAmount.IsZero() && externalAssetAmount.IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAmountTooLow
	}

	if nativeAssetDepth.Add(nativeAssetAmount).IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), errors.Wrap(errors.ErrInsufficientFunds, nativeAssetAmount.String())
	}
	if externalAssetDepth.Add(externalAssetAmount).IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), errors.Wrap(errors.ErrInsufficientFunds, externalAssetAmount.String())
	}
	if nativeAssetDepth.IsZero() || externalAssetDepth.IsZero() {
		return nativeAssetAmount, nativeAssetAmount, nil
	}

	slipAdjustmentValues := calculateSlipAdjustment(nativeAssetDepth.BigInt(), externalAssetDepth.BigInt(),
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
	ratioDiff, err := CalculateRatioDiff(externalAssetDepth.BigInt(), nativeAssetDepth.BigInt(), externalAssetAmount.BigInt(), nativeAssetAmount.BigInt())
	if err != nil {
		return sdk.ZeroUint(), sdk.ZeroUint(), err
	}
	if ratioDiff.Cmp(&ratioThresholdRat) == 1 { //if ratioDiff > ratioThreshold
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAsymmetricRatioAdd
	}

	stakeUnits := calculateStakeUnits(oldPoolUnits.BigInt(), nativeAssetDepth.BigInt(),
		externalAssetDepth.BigInt(), nativeAssetAmount.BigInt(), slipAdjustmentValues)

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

func CalcSwapResult(toRowan bool,
	X, x, Y sdk.Uint,
	pmtpCurrentRunningRate, swapFeeRate sdk.Dec, minSwapFee sdk.Uint) (sdk.Uint, sdk.Uint) {

	// if either side of the pool is empty or swap amount iz zero then return zero
	if IsAnyZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint(), sdk.ZeroUint()
	}

	rawXYK := calcRawXYK(x.BigInt(), X.BigInt(), Y.BigInt())

	pmtpFac := calcPmtpFactor(pmtpCurrentRunningRate)
	var adjustedR big.Rat
	if toRowan {
		adjustedR.Quo(&rawXYK, &pmtpFac) // adjusted = rawXYK / pmtpFac
	} else {
		adjustedR.Mul(&rawXYK, &pmtpFac) // adjusted = rawXYK * pmtpFac
	}

	swapFeeRateR := DecToRat(&swapFeeRate)
	var percentFeeR big.Rat
	percentFeeR.Mul(&adjustedR, &swapFeeRateR)

	percentFee := sdk.NewUintFromBigInt(RatIntQuo(&percentFeeR))
	adjusted := sdk.NewUintFromBigInt(RatIntQuo(&adjustedR))

	fee := sdk.MinUint(sdk.MaxUint(percentFee, minSwapFee), adjusted)
	y := adjusted.Sub(fee)

	return y, fee
}

func calcRawXYK(x, X, Y *big.Int) big.Rat {
	var numerator, denominator, xR, XR, YR, y big.Rat

	xR.SetInt(x)
	XR.SetInt(X)
	YR.SetInt(Y)
	numerator.Mul(&xR, &YR)   // x * Y
	denominator.Add(&XR, &xR) // X + x

	y.Quo(&numerator, &denominator) // y = (x * Y) / (X + x)

	return y
}

func calcPmtpFactor(r sdk.Dec) big.Rat {
	rRat := DecToRat(&r)
	one := big.NewRat(1, 1)

	one.Add(one, &rRat)

	return *one
}

func CalcSpotPriceNative(pool *types.Pool, decimalsExternal uint8, pmtpCurrentRunningRate sdk.Dec) (sdk.Dec, error) {
	X, Y := pool.ExtractDebt(pool.NativeAssetBalance, pool.ExternalAssetBalance, false)

	return CalcSpotPriceX(X, Y, types.NativeAssetDecimals, decimalsExternal, pmtpCurrentRunningRate, true)
}

func CalcSpotPriceExternal(pool *types.Pool, decimalsExternal uint8, pmtpCurrentRunningRate sdk.Dec) (sdk.Dec, error) {
	X, Y := pool.ExtractDebt(pool.ExternalAssetBalance, pool.NativeAssetBalance, true)

	return CalcSpotPriceX(X, Y, decimalsExternal, types.NativeAssetDecimals, pmtpCurrentRunningRate, false)
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

	return RatToDec(&pmtpPrice)
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
	rowanBal, externalAssetBal := pool.ExtractDebt(pool.NativeAssetBalance, pool.ExternalAssetBalance, false)

	rowanBalance := sdk.NewDecFromBigInt(rowanBal.BigInt())
	if rowanBalance.Equal(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrInValidAmount
	}
	externalAssetBalance := sdk.NewDecFromBigInt(externalAssetBal.BigInt())
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
	nativeAssetDepth, externalAssetDepth := pool.ExtractDebt(pool.NativeAssetBalance, pool.ExternalAssetBalance, false)
	return CalculateWithdrawal(
		poolUnits,
		nativeAssetDepth.String(),
		externalAssetDepth.String(),
		lp.LiquidityProviderUnits.String(),
		sdk.NewInt(types.MaxWbasis).String(),
		sdk.ZeroInt(),
	)
}

func ConvUnitsToWBasisPoints(total, units sdk.Uint) sdk.Int {
	totalDec, err := sdk.NewDecFromStr(total.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", total, err))
	}
	unitsDec, err := sdk.NewDecFromStr(units.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", total, err))
	}
	wbasis := sdk.NewDec(10000).Quo(totalDec.Quo(unitsDec))
	return wbasis.TruncateInt()
}

func ConvWBasisPointsToUnits(total sdk.Uint, wbasis sdk.Int) sdk.Uint {
	wbasisUint := sdk.NewUintFromString(wbasis.String())
	return total.Quo(sdk.NewUint(10000).Quo(wbasisUint))
}

func CalculateWithdrawalRowanValue(
	sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	pmtpCurrentRunningRate sdk.Dec, swapFeeParams types.SwapFeeParams) sdk.Uint {

	minSwapFee := GetMinSwapFee(to, swapFeeParams.TokenParams)

	X, Y, toRowan := pool.ExtractValues(to)

	X, Y = pool.ExtractDebt(X, Y, toRowan)

	value, _ := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate, swapFeeParams.SwapFeeRate, minSwapFee)

	return value
}

func SwapOne(from types.Asset,
	sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	pmtpCurrentRunningRate sdk.Dec, swapFeeParams types.SwapFeeParams) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {

	minSwapFee := GetMinSwapFee(to, swapFeeParams.TokenParams)

	X, Y, toRowan := pool.ExtractValues(to)

	var Xincl, Yincl sdk.Uint

	Xincl, Yincl = pool.ExtractDebt(X, Y, toRowan)

	priceImpact := calcPriceImpact(Xincl, sentAmount)
	swapResult, liquidityFee := CalcSwapResult(toRowan, Xincl, sentAmount, Yincl, pmtpCurrentRunningRate, swapFeeParams.SwapFeeRate, minSwapFee)

	// NOTE: impossible... pre-pmtp at least
	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), types.Pool{}, types.ErrNotEnoughAssetTokens
	}

	pool.UpdateBalances(toRowan, X, sentAmount, Y, swapResult)

	return swapResult, liquidityFee, priceImpact, pool, nil
}

func GetSwapFee(sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	pmtpCurrentRunningRate sdk.Dec, swapFeeParams types.SwapFeeParams) sdk.Uint {
	minSwapFee := GetMinSwapFee(to, swapFeeParams.TokenParams)

	X, Y, toRowan := pool.ExtractValues(to)

	X, Y = pool.ExtractDebt(X, Y, toRowan)

	swapResult, _ := CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate, swapFeeParams.SwapFeeRate, minSwapFee)

	if swapResult.GTE(Y) {
		return sdk.ZeroUint()
	}
	return swapResult
}
