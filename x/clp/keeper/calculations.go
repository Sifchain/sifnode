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

// Reverse Formula: x = ( -2*X*S + X*Y - X*sqrt( Y*(Y - 4*S) ) ) / 2*S

// Objective :  I want to send cToken and receive Rowan , I already have the amount of rowan I require (S) . The below formula should
// Tell me the amount of cToken I need to send to receive that much rowan

// if x = reverse swap result ( amount of cToken to send to get S amount of rowan)

// pool = cToken:rowan pool
// X = pool.NativeBalance   (amount of rowan in pool )
// Y = pool.ExternalBalance (amount of cToken in pool )
// S = Amount of Rowan I want to receive

func ReverseSwap(X sdk.Uint, Y sdk.Uint, S sdk.Uint) (sdk.Uint, error) {
	if S.Equal(sdk.ZeroUint()) || X.Equal(sdk.ZeroUint()) || S.Mul(sdk.NewUint(4)).GTE(Y) {
		return sdk.ZeroUint(), types.ErrNotEnoughAssetTokens
	}
	denominator := S.Add(S)                                //2*S
	innerMostTerm := Y.Sub(S.Mul(sdk.NewUint(4))).BigInt() // ( Y*(Y - 4*S)
	sqRootInnermost := innerMostTerm.Sqrt(innerMostTerm)   // sqrt( Y*(Y - 4*S)
	term3 := X.Mul(sdk.NewUintFromBigInt(sqRootInnermost)) // X*sqrt( Y*(Y - 4*S)
	term2 := X.Mul(Y)                                      //X*Y
	term1 := X.Mul(S).Mul(sdk.NewUint(2))                  //2*X*S
	numerator := term2.Sub(term1).Sub(term3)               //  X*Y - (2*X*S)  - (X*sqrt( Y*(Y - 4*S))
	return numerator.Quo(denominator), nil
}

func SwapOne(from types.Asset, sentAmount sdk.Uint, to types.Asset, pool types.Pool) (sdk.Uint, sdk.Uint, sdk.Uint, types.Pool, error) {
	var X sdk.Uint
	var Y sdk.Uint
	toRowan := true
	if to == types.GetSettlementAsset() {
		Y = pool.NativeAssetBalance
		X = pool.ExternalAssetBalance
	} else {
		X = pool.NativeAssetBalance
		Y = pool.ExternalAssetBalance
		toRowan = false
	}
	x := sentAmount
	liquidityFee, err := calcLiquidityFee(pool.ExternalAsset.Symbol, toRowan, X, x, Y)
	if err != nil {
		return sdk.Uint{}, sdk.Uint{}, sdk.Uint{}, types.Pool{}, err
	}
	priceImpact, err := calcPriceImpact(X, x)
	if err != nil {
		return sdk.Uint{}, sdk.Uint{}, sdk.Uint{}, types.Pool{}, err
	}
	swapResult, err := calcSwapResult(pool.ExternalAsset.Symbol, toRowan, X, x, Y)
	if err != nil {
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

func GetSwapFee(sentAmount sdk.Uint, to types.Asset, pool types.Pool) sdk.Uint {
	var X sdk.Uint
	var Y sdk.Uint
	toRowan := true
	if to == types.GetSettlementAsset() {
		Y = pool.NativeAssetBalance
		X = pool.ExternalAssetBalance
	} else {
		X = pool.NativeAssetBalance
		Y = pool.ExternalAssetBalance
		toRowan = false
	}
	x := sentAmount
	swapResult, err := calcSwapResult(pool.ExternalAsset.Symbol, toRowan, X, x, Y)
	if err != nil {
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

func CalculatePoolUnits(symbol string, oldPoolUnits, nativeAssetBalance, externalAssetBalance,
	nativeAssetAmount, externalAssetAmount sdk.Uint) (sdk.Uint, sdk.Uint, error) {
	normalizationFactor := sdk.NewDec(1)
	nf, ok := types.GetNormalizationMap()[symbol[1:]]
	adjustExternalToken := true
	if ok {
		diffFactor := 18 - nf
		if diffFactor < 0 {
			diffFactor = nf - 18
			adjustExternalToken = false
		}
		normalizationFactor = sdk.NewDec(10).Power(uint64(diffFactor))
	}
	if adjustExternalToken {
		externalAssetAmount = externalAssetAmount.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt())) // Convert token which are not E18 to E18 format
		externalAssetBalance = externalAssetBalance.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
	} else {
		nativeAssetAmount = nativeAssetAmount.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		nativeAssetBalance = nativeAssetBalance.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
	}

	inputs := []sdk.Uint{oldPoolUnits, nativeAssetBalance, externalAssetBalance,
		nativeAssetAmount, externalAssetAmount}

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
	numerator := P.Mul(a.Mul(R).Add(A.Mul(r)))
	denominator := sdk.NewDec(2).Mul(A).Mul(R)
	stakeUnits := numerator.Quo(denominator).Mul(slipAdjustment)
	newPoolUnit := P.Add(stakeUnits)
	newPoolUnit = IncreasePrecision(newPoolUnit, minLen)
	stakeUnits = IncreasePrecision(stakeUnits, minLen)

	return sdk.NewUintFromBigInt(newPoolUnit.RoundInt().BigInt()), sdk.NewUintFromBigInt(stakeUnits.RoundInt().BigInt()), nil
}

func calcLiquidityFee(symbol string, toRowan bool, X, x, Y sdk.Uint) (sdk.Uint, error) {
	if X.IsZero() && x.IsZero() {
		return sdk.ZeroUint(), nil
	}
	if !ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint(), nil
	}
	normalizationFactor := sdk.NewDec(1)
	nf, ok := types.GetNormalizationMap()[symbol[1:]]
	adjustExternalToken := true
	if ok {
		diffFactor := 18 - nf
		if diffFactor < 0 {
			diffFactor = nf - 18
			adjustExternalToken = false
		}
		normalizationFactor = sdk.NewDec(10).Power(uint64(diffFactor))
	}

	if adjustExternalToken {
		if toRowan {
			X = X.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
			x = x.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		} else {
			Y = Y.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		}
	} else {
		if toRowan {
			X = X.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
			x = x.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		} else {
			Y = Y.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		}
	}

	// Assuming the max supply for any token in the world to be 1 trillion
	minLen := int64(6)

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

func calcSwapResult(symbol string, toRowan bool, X, x, Y sdk.Uint) (sdk.Uint, error) {
	if !ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint(), nil
	}
	normalizationFactor := sdk.NewDec(1)
	nf, ok := types.GetNormalizationMap()[symbol[1:]]
	adjustExternalToken := true
	if ok {
		diffFactor := 18 - nf
		if diffFactor < 0 {
			diffFactor = nf - 18
			adjustExternalToken = false
		}
		normalizationFactor = sdk.NewDec(10).Power(uint64(diffFactor))
	}

	if adjustExternalToken {
		if toRowan {
			X = X.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
			x = x.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		} else {
			Y = Y.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		}
	} else {
		if toRowan {
			X = X.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
			x = x.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
		} else {
			Y = Y.Mul(sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt()))
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

	return sdk.NewUintFromBigInt(y.RoundInt().BigInt()), nil
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
	return CalculateWithdrawal(poolUnits, nativeAssetBalance.String(), externalAssetBalance.String(),
		lp.LiquidityProviderUnits.String(), sdk.NewInt(types.MaxWbasis).String(), sdk.ZeroInt())
}
