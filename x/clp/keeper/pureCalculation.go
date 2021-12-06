package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ValidateZero(inputs []sdk.Uint) bool {
	return ValidateZero(inputs)
}

func ValidateZero(inputs []sdk.Uint) bool {
	for _, val := range inputs {
		if val.IsZero() {
			return false
		}
	}
	return true
}

func (k Keeper) ReducePrecision(dec sdk.Dec, po int64) sdk.Dec {
	return ReducePrecision(dec, po)
}

func ReducePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.Quo(p)
}

func (k Keeper) IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec {
	return IncreasePrecision(dec, po)
}

func IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.Mul(p)
}

func (k Keeper) GetMinLen(inputs []sdk.Uint) int64 {
	return GetMinLen(inputs)
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
	if minLen <= 6 {
		return int64(6)
	}
	if maxInputLen >= 27 {
		lenDiff := maxInputLen - 27
		if lenDiff < minLen {
			if lenDiff <= 6 {
				return int64(6)
			}
			return int64(lenDiff)
		}
	}
	return int64(minLen - 1)
}
