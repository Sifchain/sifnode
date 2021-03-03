package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common/math"
)

//------------------------------------------------------------------------------------------------------------------
// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
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
// slipAdjustment = (1 - ABS((R a - r A)/((2 r + R) (a + A))))
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
		externalAssetAmount = externalAssetAmount.Mul(sdk.NewUintFromBigInt(normalizationFactor.BigInt())) // Convert token which are not E18 to E18 format
		externalAssetBalance = externalAssetBalance.Mul(sdk.NewUintFromBigInt(normalizationFactor.BigInt()))
	} else {
		nativeAssetAmount = nativeAssetAmount.Mul(sdk.NewUintFromBigInt(normalizationFactor.BigInt()))
		nativeAssetBalance = nativeAssetBalance.Mul(sdk.NewUintFromBigInt(normalizationFactor.BigInt()))
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

	slipAdjDenominator := (r.MulInt64(2).Add(R)).Mul(a.Add(A))
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

func GetMinLen(inputs []sdk.Uint) int64 {
	minLen := math.MaxInt64
	for _, val := range inputs {
		currentLen := len(val.String())
		if currentLen < minLen {
			minLen = currentLen
		}
	}
	if minLen <= 1 {
		return int64(1)
	}
	return int64(minLen - 1)
}
