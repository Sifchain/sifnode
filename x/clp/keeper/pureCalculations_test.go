package keeper_test

import (
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func TestKeeper_DecRatIdentity(t *testing.T) {

	// Test for: dec == RatToDec(DecToRat(dec))
	// NOTE: rat == DecToRat(RatToDec(rat)) does not hold for all rat in big.Rat due to loss of precision when converting to sdk.Dec
	// E.g rat 1/3 becomes dec 0.333333... which then becomes 1000.../3000... and not 1/3
	testIdentity := func(t *rapid.T) {
		expected := genDec(t)

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
	minRat := clpkeeper.DecToRat(&min)

	// 18 zeros
	assert.Equal(t, "1/1000000000000000000", minRat.String())

	// https://github.com/cosmos/cosmos-sdk/blob/main/types/decimal.go#L34
	max := sdk.MustNewDecFromStr("43556142965880123323311949751266331066368") // 2**315
	maxRat := clpkeeper.DecToRat(&max)

	assert.Equal(t, "43556142965880123323311949751266331066368/1", maxRat.String())
}

func genDec(t *rapid.T) sdk.Dec {
	const numInt64 = 5 // 4 * 64bit = 256 bits at most; max size sdk.Dec

	ints := rapid.ArrayOf(numInt64, genInt64ButZero()).Draw(t, "ints").([numInt64]int64)
	dec := sdk.NewDec(ints[0])

	for i := 1; i < numInt64-1; i++ {
		dec = dec.MulInt64(ints[i])
	}

	dec = dec.QuoInt64(ints[numInt64-1])

	return dec
}
func genInt64ButZero() *rapid.Generator {
	return rapid.OneOf(rapid.Int64Max(-1), rapid.Int64Min(1))
}
