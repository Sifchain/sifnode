//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"math/big"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPInterestLiabilities(mtp *types.MTP, interestRate sdk.Dec, epochPosition, epochLength int64) sdk.Uint {
	var interestRational, liabilitiesRational, rate, epochPositionRational, epochLengthRational big.Rat

	rate.SetFloat64(interestRate.MustFloat64())

	liabilitiesRational.SetInt(mtp.Liabilities.BigInt().Add(mtp.Liabilities.BigInt(), mtp.InterestUnpaid.BigInt()))
	interestRational.Mul(&rate, &liabilitiesRational)

	if epochPosition > 0 { // prorate interest if within epoch
		epochPositionRational.SetInt64(epochPosition)
		epochLengthRational.SetInt64(epochLength)
		epochPositionRational.Quo(&epochPositionRational, &epochLengthRational)
		interestRational.Mul(&interestRational, &epochPositionRational)
	}

	interestNew := interestRational.Num().Quo(interestRational.Num(), interestRational.Denom())

	return sdk.NewUintFromBigInt(interestNew.Add(interestNew, mtp.InterestUnpaid.BigInt()))
}

func CalcInterestRate(interestRateMax,
	interestRateMin,
	interestRateIncrease,
	interestRateDecrease,
	healthGainFactor,
	sQ sdk.Dec,
	pool clptypes.Pool) (sdk.Dec, error) {

	prevInterestRate := pool.InterestRate

	externalAssetBalance := sdk.NewDecFromBigInt(pool.ExternalAssetBalance.BigInt())
	ExternalLiabilities := sdk.NewDecFromBigInt(pool.ExternalLiabilities.BigInt())
	NativeAssetBalance := sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt())
	NativeLiabilities := sdk.NewDecFromBigInt(pool.NativeLiabilities.BigInt())

	mul1 := externalAssetBalance.Add(ExternalLiabilities).Quo(externalAssetBalance)
	mul2 := NativeAssetBalance.Add(NativeLiabilities).Quo(NativeAssetBalance)

	targetInterestRate := healthGainFactor.Mul(mul1).Mul(mul2).Add(sQ)

	interestRateChange := targetInterestRate.Sub(prevInterestRate)
	interestRate := prevInterestRate
	if interestRateChange.LTE(interestRateDecrease.Mul(sdk.NewDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
		interestRate = targetInterestRate
	} else if interestRateChange.GT(interestRateIncrease) {
		interestRate = prevInterestRate.Add(interestRateIncrease)
	} else if interestRateChange.LT(interestRateDecrease.Mul(sdk.NewDec(-1))) {
		interestRate = prevInterestRate.Sub(interestRateDecrease)
	}

	newInterestRate := interestRate

	if interestRate.GT(interestRateMin) && interestRate.LT(interestRateMax) {
		newInterestRate = interestRate
	} else if interestRate.LTE(interestRateMin) {
		newInterestRate = interestRateMin
	} else if interestRate.GTE(interestRateMax) {
		newInterestRate = interestRateMax
	}

	return newInterestRate, nil
}
