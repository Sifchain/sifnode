package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestTradeSlip(t *testing.T) {
	type TestCase struct {
		X        string `json:"x"`
		Y        string `json:"X"`
		Expected string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"TradeSlip"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/sample_trade_slip.json")
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	assert.NoError(t, err)
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	for _, test := range testcases {
		e := calcTradeSlip(GetInt(t, test.Y), GetInt(t, test.X))
		assert.Equal(t, e, GetInt(t, test.Expected))
	}
}

func TestLiquidityFee(t *testing.T) {
	type TestCase struct {
		Sx       string `json:"x"`
		X        string `json:"X"`
		Y        string `json:"Y"`
		Expected string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"LiquidityFee"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/sample_liquidity_fee.json")
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	assert.NoError(t, err)
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	totalCount := 0
	failedCount := 0
	for _, test := range testcases {
		totalCount++
		res := calcLiquidityFee(GetInt(t, test.X), GetInt(t, test.Sx), GetInt(t, test.Y))
		assert.Equal(t, res, GetInt(t, test.Expected))
		if !res.Equal(GetInt(t, test.Expected)) {
			failedCount++
		}
	}
	fmt.Println("TotalCount :", totalCount)
	fmt.Println("FailedCount :", failedCount)
}

func TestCalculatePoolUnits(t *testing.T) {
	type TestCase struct {
		NativeAdded      string `json:"r"`
		ExternalAdded    string `json:"a"`
		NativeBalance    string `json:"R"`
		ExternalBalance  string `json:"A"`
		PoolUnitsBalance string `json:"P"`
		Expected         string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"PoolUnits"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/sample_pool_units.json")
	assert.NoError(t, err)
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	for _, test := range testcases {
		_, newPoolUnit, err := CalculatePoolUnits(
			sdk.NewUintFromString(test.PoolUnitsBalance),
			sdk.NewUintFromString(test.NativeBalance),
			sdk.NewUintFromString(test.ExternalBalance),
			sdk.NewUintFromString(test.NativeAdded),
			sdk.NewUintFromString(test.ExternalAdded),
		)
		assert.NoError(t, err)
		assert.Equal(t, newPoolUnit, sdk.NewUintFromString(test.Expected))
		if !newPoolUnit.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("%s", test)
		}
	}
}

func GetInt(t *testing.T, x string) sdk.Uint {
	flt, _, err := big.ParseFloat(x, 10, 0, big.ToNearestEven)
	if err != nil {
		fmt.Println("Error parsing float :", x)
	}
	var i = new(big.Int)
	i, _ = flt.Int(i)
	X := sdk.NewDecFromBigInt(i)
	assert.NoError(t, err)
	if err != nil {
		return sdk.ZeroUint()
	}
	Xr := X.Mul(sdk.NewDec(10).Power(18)).RoundInt()
	return sdk.NewUintFromBigInt(Xr.BigInt())
}
