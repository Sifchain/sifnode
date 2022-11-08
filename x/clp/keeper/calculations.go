package keeper

import (
	"fmt"

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
	normalizationFactor sdk.Dec,
	adjustExternalToken bool,
	pmtpCurrentRunningRate sdk.Dec) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {

	X, x, Y, toRowan := SetInputs(sentAmount, to, pool)
	liquidityFee, err := CalcLiquidityFee(toRowan, normalizationFactor, adjustExternalToken, X, x, Y)
	if err != nil {
		// this branch will never be reached as err will always be nil
		return sdk.Uint{}, sdk.Uint{}, sdk.Uint{}, types.Pool{}, err
	}
	priceImpact, err := calcPriceImpact(X, x)
	if err != nil {
		// this branch will never be reached as err will always be nil
		return sdk.Uint{}, sdk.Uint{}, sdk.Uint{}, types.Pool{}, err
	}
	swapResult, err := CalcSwapResult(toRowan, normalizationFactor, adjustExternalToken, X, x, Y, pmtpCurrentRunningRate)
	if err != nil {
		// this branch will never be reached as err will always be nil
		return sdk.Uint{}, sdk.Uint{}, sdk.Uint{}, types.Pool{}, err
	}
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

func CalcSwapPrice(from types.Asset,
	sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	normalizationFactor sdk.Dec,
	adjustExternalToken bool,
	pmtpCurrentRunningRate sdk.Dec) sdk.Dec {

	X, x, Y, toRowan := SetInputs(sentAmount, to, pool)

	swapResult := CalcSwapPriceResult(toRowan, normalizationFactor, adjustExternalToken, X, x, Y, pmtpCurrentRunningRate)

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

func SetInputs(sentAmount sdk.Uint, to types.Asset, pool types.Pool) (sdk.Uint, sdk.Uint, sdk.Uint, bool) {
	var X sdk.Uint
	var Y sdk.Uint
	var x sdk.Uint
	toRowan := true
	if to == types.GetSettlementAsset() {
		Y = pool.NativeAssetBalance
		X = pool.ExternalAssetBalance
	} else {
		X = pool.NativeAssetBalance
		Y = pool.ExternalAssetBalance
		toRowan = false
	}
	x = sentAmount

	return X, x, Y, toRowan
}

func GetSwapFee(sentAmount sdk.Uint,
	to types.Asset,
	pool types.Pool,
	normalizationFactor sdk.Dec,
	adjustExternalToken bool,
	pmtpCurrentRunningRate sdk.Dec) sdk.Uint {
	X, x, Y, toRowan := SetInputs(sentAmount, to, pool)
	swapResult, err := CalcSwapResult(toRowan, normalizationFactor, adjustExternalToken, X, x, Y, pmtpCurrentRunningRate)
	if err != nil {
		// this branch will never be reached as err will always be nil
		return sdk.Uint{}
	}

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
	externalAssetAmount sdk.Uint, normalizationFactor sdk.Dec, adjustExternalToken bool, symmetryThreshold sdk.Dec) (sdk.Uint, sdk.Uint, error) {
	nf := sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt())

	if adjustExternalToken {
		externalAssetAmount = externalAssetAmount.Mul(nf) // Convert token which are not E18 to E18 format
		externalAssetBalance = externalAssetBalance.Mul(nf)
	} else {
		nativeAssetAmount = nativeAssetAmount.Mul(nf)
		nativeAssetBalance = nativeAssetBalance.Mul(nf)
	}

	inputs := []sdk.Uint{oldPoolUnits, nativeAssetBalance, externalAssetBalance, nativeAssetAmount, externalAssetAmount}

	if nativeAssetAmount.IsZero() && externalAssetAmount.IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAmountTooLow
	}

	minLen := GetMinLen(inputs)

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

	P = ReducePrecision(P, minLen)
	R = ReducePrecision(R, minLen)
	A = ReducePrecision(A, minLen)
	a = ReducePrecision(a, minLen)
	r = ReducePrecision(r, minLen)

	slipAdjDenominator := (r.Add(R)).Mul(a.Add(A))
	var slipAdjustment sdk.Dec
	if R.Mul(a).GT(r.Mul(A)) {
		slipAdjustment = R.Mul(a).Sub(r.Mul(A)).Quo(slipAdjDenominator)
	} else {
		slipAdjustment = r.Mul(A).Sub(R.Mul(a)).Quo(slipAdjDenominator)
	}
	slipAdjustment = sdk.NewDec(1).Sub(slipAdjustment)

	if sdk.OneDec().Sub(slipAdjustment).GT(symmetryThreshold) {
		return sdk.ZeroUint(), sdk.ZeroUint(), types.ErrAsymmetricAdd
	}

	numerator := P.Mul(a.Mul(R).Add(A.Mul(r)))
	denominator := sdk.NewDec(2).Mul(A).Mul(R)
	stakeUnits := numerator.Quo(denominator).Mul(slipAdjustment)
	newPoolUnit := P.Add(stakeUnits)
	newPoolUnit = IncreasePrecision(newPoolUnit, minLen)
	stakeUnits = IncreasePrecision(stakeUnits, minLen)

	return sdk.NewUintFromBigInt(newPoolUnit.RoundInt().BigInt()), sdk.NewUintFromBigInt(stakeUnits.RoundInt().BigInt()), nil
}

func CalcLiquidityFee(toRowan bool, normalizationFactor sdk.Dec, adjustExternalToken bool, X, x, Y sdk.Uint) (sdk.Uint, error) {
	if X.IsZero() && x.IsZero() {
		return sdk.ZeroUint(), nil
	}
	if !ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint(), nil
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

	n := xd.Mul(xd).Mul(Yd)
	s := xd.Add(Xd)
	d := s.Mul(s)
	y := n.Quo(d)

	y = IncreasePrecision(y, minLen)
	if !toRowan {
		y = y.Quo(normalizationFactor)
	}

	return sdk.NewUintFromBigInt(y.RoundInt().BigInt()), nil
}

func CalcSwapResult(toRowan bool,
	normalizationFactor sdk.Dec,
	adjustExternalToken bool,
	X, x, Y sdk.Uint,
	pmtpCurrentRunningRate sdk.Dec) (sdk.Uint, error) {
	if !ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint(), nil
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
	if !toRowan {
		y = y.Quo(normalizationFactor)
	}
	y = CalcSwapPmtp(toRowan, y, pmtpCurrentRunningRate)
	return sdk.NewUintFromBigInt(y.RoundInt().BigInt()), nil
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

func calcPriceImpact(X, x sdk.Uint) (sdk.Uint, error) {
	if x.IsZero() {
		return sdk.ZeroUint(), nil
	}
	d := x.Add(X)
	return x.Quo(d), nil
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
