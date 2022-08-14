package keeper

import (
	"math"
	"math/big"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DecToRat(d *sdk.Dec) big.Rat {
	var rat big.Rat

	rat.SetInt(d.BigInt())
	decimals := int64(math.Pow10(sdk.Precision)) // 10**18
	denom := big.NewRat(decimals, 1)
	rat.Quo(&rat, denom)

	return rat
}

// The sdk.Dec returned by this method can exceed the sdk.Decimal maxDecBitLen
func RatToDec(r *big.Rat) sdk.Dec {
	num := r.Num()
	denom := r.Denom() // big.Rat guarantees that denom is always > 0

	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil) // 10**18

	var d big.Int
	d.Mul(num, multiplier)
	d.Quo(&d, denom)

	return sdk.NewDecFromBigIntWithPrec(&d, sdk.Precision)
}

func RatIntQuo(r *big.Rat) *big.Int {
	var i big.Int

	num := r.Num()
	denom := r.Denom()
	return i.Quo(num, denom)
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
