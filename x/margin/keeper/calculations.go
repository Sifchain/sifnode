package keeper

import (
	"math/big"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPInterestLiabilities(mtp *types.MTP, interestRate sdk.Dec, epochPosition, epochLength int64) sdk.Uint {
	var interestRational, liabilitiesRational, rate, epochPositionRational, epochLengthRational big.Rat

	rate.SetFloat64(interestRate.MustFloat64())

	liabilitiesRational.SetInt(mtp.Liabilities.BigInt().Add(mtp.Liabilities.BigInt(), mtp.InterestUnpaidCollateral.BigInt()))
	interestRational.Mul(&rate, &liabilitiesRational)

	if epochPosition > 0 { // prorate interest if within epoch
		epochPositionRational.SetInt64(epochPosition)
		epochLengthRational.SetInt64(epochLength)
		epochPositionRational.Quo(&epochPositionRational, &epochLengthRational)
		interestRational.Mul(&interestRational, &epochPositionRational)
	}

	interestNew := interestRational.Num().Quo(interestRational.Num(), interestRational.Denom())

	interestNewUint := sdk.NewUintFromBigInt(interestNew.Add(interestNew, mtp.InterestUnpaidCollateral.BigInt()))
	// round up to lowest digit if interest too low and rate not 0
	if interestNewUint.IsZero() && !interestRate.IsZero() {
		interestNewUint = sdk.OneUint()
	}

	return interestNewUint
}
