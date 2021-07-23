package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
	maxInputLen := 1
	for _, val := range inputs {
		currentLen := len(val.String())
		if currentLen < minLen {
			minLen = currentLen
		}
		if currentLen > maxInputLen {
			maxInputLen = currentLen
		}
	}
	if maxInputLen >= 27 {
		lenDiff := maxInputLen - 27
		if lenDiff > minLen {
			return int64(minLen - 1)
		} else {
			return int64(lenDiff)
		}
	}
	return int64(1)
}
