//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

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
	pmtpCurrentRunningRate sdk.Dec, marginEnabled bool) sdk.Uint {

	X, Y, toRowan := pool.ExtractValues(to)

	if marginEnabled {
		X, Y = pool.ExtractDebt(X, Y, toRowan)
	}

	return CalcSwapResult(toRowan, X, sentAmount, Y, pmtpCurrentRunningRate)
}
