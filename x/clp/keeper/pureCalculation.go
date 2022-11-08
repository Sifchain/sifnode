package keeper

import (
	"fmt"
	"math"
	"math/big"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	maxDecBitLen = 315 // sdk.Dec doesn't export this value but see here: https://github.com/cosmos/cosmos-sdk/blob/main/types/decimal.go#L34
)

func DecToRat(d *sdk.Dec) big.Rat {
	var rat big.Rat

	rat.SetInt(d.BigInt())
	decimals := int64(math.Pow10(sdk.Precision)) // 10**18
	denom := big.NewRat(decimals, 1)
	rat.Quo(&rat, denom)

	return rat
}

func RatToDec(r *big.Rat) (sdk.Dec, error) {
	num := r.Num()
	denom := r.Denom() // big.Rat guarantees that denom is always > 0

	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil) // 10**18

	var d big.Int
	d.Mul(num, multiplier)
	d.Quo(&d, denom)

	// There's a bug in the SDK which allows sdk.NewDecFromBigIntWithPrec to create an sdk.Dec with > maxDecBitLen bits
	// This leads to an error when attempting to unmarshal such sdk.Decs
	if d.BitLen() > maxDecBitLen {
		return sdk.ZeroDec(), fmt.Errorf("decimal out of range; bitLen: got %d, max %d", d.BitLen(), maxDecBitLen)
	}

	return sdk.NewDecFromBigIntWithPrec(&d, sdk.Precision), nil
}

func RatIntQuo(r *big.Rat) *big.Int {
	var i big.Int

	num := r.Num()
	denom := r.Denom()
	return i.Quo(num, denom)
}

func ApproxRatSquareRoot(x *big.Rat) *big.Int {
	var i big.Int
	return i.Sqrt(RatIntQuo(x))
}

func IsAnyZero(inputs []sdk.Uint) bool {
	for _, val := range inputs {
		if val.IsZero() {
			return true
		}
	}
	return false
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
