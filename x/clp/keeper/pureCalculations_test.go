package keeper_test

import (
	"fmt"
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func TestKeeper_DecRatIdentity(t *testing.T) {

	// Tests for: dec == RatToDec(DecToRat(dec))
	// NOTE: rat == DecToRat(RatToDec(rat)) does not hold for all rat in big.Rat due to loss of precision when converting to sdk.Dec
	// E.g rat 1/3 becomes dec 0.333333 thich then becomes 1000.../3000... and not 1/3
	testIdentity := func(t *rapid.T) {
		expected := genDec(t)
		fmt.Println("expected ", expected.String())

		rat := clpkeeper.DecToRat(&expected)
		actual := clpkeeper.RatToDec(&rat)

		if !expected.Equal(actual) {
			t.Fatalf("\nexpected %s\nactual   %s", expected.String(), actual.String())
		}
	}

	rapid.Check(t, testIdentity)
}

func TestKeeper_CanConvertDecBoundaries(t *testing.T) {
	min := sdk.SmallestDec()
	min_rat := clpkeeper.DecToRat(&min)

	// 18 zeros
	assert.Equal(t, "1/1000000000000000000", min_rat.String())

	// https://github.com/cosmos/cosmos-sdk/blob/main/types/decimal.go#L34
	max := sdk.MustNewDecFromStr("43556142965880123323311949751266331066368") // 2**315
	max_rat := clpkeeper.DecToRat(&max)

	// 18 zeros
	assert.Equal(t, "43556142965880123323311949751266331066368/1", max_rat.String())
}

func genDec(t *rapid.T) sdk.Dec {
	const num_int64 = 4 // 4 * 64bit = 256 bits at most; max size sdk.Dec
	ints := rapid.ArrayOf(num_int64, rapid.Int64()).Draw(t, "ints").([num_int64]int64)
	dec := sdk.NewDec(ints[0])

	for i := 1; i < num_int64; i++ {
		dec = dec.MulInt64(ints[i])
	}

	return dec
}
