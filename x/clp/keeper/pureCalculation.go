package keeper

import (
	"math"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func IsAnyZero(inputs []sdk.Uint) bool {
	for _, val := range inputs {
		if val.IsZero() {
			return true
		}
	}
	return false
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

func Int64ToUint8Safe(x int64) (uint8, error) {
	trial := uint8(x)
	if int64(trial) != x {
		return 0, types.ErrTypeCast
	}
	return trial, nil
}

func Abs(a int16) uint16 {
	if a < 0 {
		return uint16(-a)
	}
	return uint16(a)
}
