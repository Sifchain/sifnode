package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestKeeper_DecRatIdentity(t *testing.T) {

	// Test for: dec == RatToDec(DecToRat(dec))
	// NOTE: rat == DecToRat(RatToDec(rat)) does not hold for all rat in big.Rat due to loss of precision when converting to sdk.Dec
	// E.g rat 1/3 becomes dec 0.333333... which then becomes 1000.../3000... and not 1/3
	testIdentity := func(t *rapid.T) {
		expected := genDec(t)

		rat := clpkeeper.DecToRat(&expected)
		actual, _ := clpkeeper.RatToDec(&rat)

		if !expected.Equal(actual) {
			t.Fatalf("\nexpected %s\nactual   %s", expected.String(), actual.String())
		}
	}

	rapid.Check(t, testIdentity)
}

func TestKeeper_CanConvertMinDec(t *testing.T) {
	min := sdk.SmallestDec()
	minRat := clpkeeper.DecToRat(&min)

	// 18 zeros
	assert.Equal(t, "1/1000000000000000000", minRat.String())
}

func genDec(t *rapid.T) sdk.Dec {
	const numInt64 = 5 // 4 * 64bit = 256 bits at most; max size sdk.Dec

	// ints := rapid.ArrayOf(numInt64, genInt64ButZero()).Draw(t, "ints").([numInt64]int64)
	ints := rapid.SliceOfN[int64](genInt64ButZero(), numInt64, numInt64).Draw(t, "ints")
	dec := sdk.NewDec(ints[0])

	for i := 1; i < numInt64-1; i++ {
		dec = dec.MulInt64(ints[i])
	}

	dec = dec.QuoInt64(ints[numInt64-1])

	return dec
}

func genInt64ButZero() *rapid.Generator[int64] {
	return rapid.OneOf(rapid.Int64Max(-1), rapid.Int64Min(1))
}

func TestKeeper_RatToDec(t *testing.T) {
	testcases := []struct {
		name     string
		num      *big.Int
		denom    *big.Int
		expected sdk.Dec
		err      error
	}{
		{
			name:     "small values",
			num:      big.NewInt(1),
			denom:    big.NewInt(3),
			expected: sdk.MustNewDecFromStr("0.333333333333333333"),
		},
		{
			name:     "small values",
			num:      big.NewInt(7),
			denom:    big.NewInt(3),
			expected: sdk.MustNewDecFromStr("2.333333333333333333"),
		},
		{
			name:     "negative numerator",
			num:      big.NewInt(-7),
			denom:    big.NewInt(3),
			expected: sdk.MustNewDecFromStr("-2.333333333333333333"),
		},
		{
			name:  "big numbers",
			num:   big.NewInt(1).Exp(big.NewInt(2), big.NewInt(400), nil), // 2**400
			denom: big.NewInt(3),
			err:   errors.New("decimal out of range; bitLen: got 459, max 315"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			var rat big.Rat
			rat.SetFrac(tc.num, tc.denom)
			y, err := clpkeeper.RatToDec(&rat)

			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, y)
		})
	}
}

// If the cosmos SDK reduces the value of maxDecBitLen below 315 then (without
// further changes) RatToDec will be able to return an sdk.Dec with 316 bits.
// Such sdk.Decs would cause panics and errors when unmarshaling. This test
// will fail if the cosmos SDK reduces the value of maxDecBitLen < 315.
func TestKeeper_RatToDec_Overflow(t *testing.T) {
	numerator, _ := big.NewInt(1).SetString("66749594872528440074844428317798503581334516323645399060845050244444366430645", 10) // (2**315 - 1) / 10**18
	denominator := big.NewInt(1)

	var rat big.Rat
	rat.SetFrac(numerator, denominator)

	// The rat passed here uses 315 bits when represented as an sdk.Dec. Even if the SDK
	// reduces the maxDecBitLen < 315 (without fixing sdk.NewDecFromBigIntWithPrec) the RatToDec
	// will still succeed - returning an sdk.Dec which cannot be unmarshaled or that will cause a panic
	// even when adding zero.
	dec, err := clpkeeper.RatToDec(&rat)

	require.NoError(t, err)

	// If the SDK ever reduces the value of maxDecBitLen < 315 then this will panic
	dec.Add(sdk.NewDec(0))
}

func TestKeeper_Int64ToUint8Safe(t *testing.T) {

	testcases := []struct {
		name      string
		x         int64
		expected  uint8
		errString error
	}{
		{
			name:     "success",
			x:        128,
			expected: 128,
		},
		{
			name:     "success 0",
			x:        0,
			expected: 0,
		},
		{
			name:     "success 255",
			x:        255,
			expected: 255,
		},
		{
			name:      "fail - below range",
			x:         -1,
			errString: errors.New("Could not perform type cast"),
		},
		{
			name:      "fail - above range",
			x:         256,
			errString: errors.New("Could not perform type cast"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			y, err := clpkeeper.Int64ToUint8Safe(tc.x)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, y)
		})
	}
}

func TestKeeper_Abs(t *testing.T) {

	testcases := []struct {
		name     string
		x        int16
		expected uint16
	}{
		{
			name:     "no change",
			x:        128,
			expected: 128,
		},
		{
			name:     "0 case",
			x:        0,
			expected: 0,
		},
		{
			name:     "flip sign",
			x:        -100,
			expected: 100,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			y := clpkeeper.Abs(tc.x)

			require.Equal(t, tc.expected, y)
		})
	}
}

func TestKeeper_DecToRat(t *testing.T) {
	testcases := []struct {
		name     string
		dec      sdk.Dec
		expected big.Rat
	}{
		{
			name:     "small values",
			dec:      sdk.MustNewDecFromStr("0.333333333333333333"),
			expected: *big.NewRat(333333333333333333, 1000000000000000000),
		},
		{
			name:     "small values",
			dec:      sdk.MustNewDecFromStr("2.333333333333333333"),
			expected: *big.NewRat(2333333333333333333, 1000000000000000000),
		},
		{
			name:     "negative numerator",
			dec:      sdk.MustNewDecFromStr("-2.333333333333333333"),
			expected: *big.NewRat(-2333333333333333333, 1000000000000000000),
		},
		{
			name:     "big numbers",
			dec:      sdk.NewDecFromBigIntWithPrec(getFirstArg(big.NewInt(1).SetString("860749959362302863218639724001003958109901930943074504276886452180215874005613731543215117760045943811967723990915831125333333333333333333", 10)), 18),
			expected: getFirstArgRat(new(big.Rat).SetString("860749959362302863218639724001003958109901930943074504276886452180215874005613731543215117760045943811967723990915831125333333333333333333/1000000000000000000")),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			y := clpkeeper.DecToRat(&tc.dec)

			require.Equal(t, tc.expected.String(), y.String())
		})
	}
}

func getFirstArgRat(r *big.Rat, ignore bool) big.Rat {
	return *r
}

func TestKeeper_RatIntQuo(t *testing.T) {
	testcases := []struct {
		name     string
		rat      big.Rat
		expected big.Int
	}{
		{
			name:     "small values",
			rat:      *big.NewRat(6, 3),
			expected: *big.NewInt(2),
		},
		{
			name:     "small values",
			rat:      *big.NewRat(7, 3),
			expected: *big.NewInt(2),
		},
		{
			name:     "negative numerator",
			rat:      *big.NewRat(-7, 3),
			expected: *big.NewInt(-2),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			y := clpkeeper.RatIntQuo(&tc.rat)

			require.Equal(t, tc.expected.String(), y.String())
		})
	}
}

func TestKeeper_ApproxRatSquareRoot(t *testing.T) {
	testcases := []struct {
		name     string
		x        big.Rat
		expected big.Int
	}{
		{
			name:     "square number",
			x:        *big.NewRat(50, 2),
			expected: *big.NewInt(5),
		},
		{
			name:     "non square integer number",
			x:        *big.NewRat(100, 2),
			expected: *big.NewInt(7),
		},
		{
			name:     "non square non number",
			x:        *big.NewRat(101, 2),
			expected: *big.NewInt(7),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			y := clpkeeper.ApproxRatSquareRoot(&tc.x)

			require.Equal(t, tc.expected.String(), y.String())
		})
	}
}
