package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

//func TestSwapOne(t *testing.T) {
//	type TestCase struct {
//		Ax       string `json:"ax"`
//		AX       string `json:"aX"`
//		AY       string `json:"aY"`
//		BX       string `json:"bX"`
//		BY       string `json:"bY"`
//		Expected string `json:"expected"`
//	}
//	type Test struct {
//		TestType []TestCase `json:"Swap"`
//	}
//	file, err := ioutil.ReadFile("../../../test/test-tables/sample_swaps.json")
//	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
//	assert.NoError(t, err)
//	var test Test
//	err = json.Unmarshal(file, &test)
//	assert.NoError(t, err)
//	testcases := test.TestType
//	for _, test := range testcases {
//		SwapOne()
//	}
//}
//
//func TestSlip(t *testing.T) {
//	type TestCase struct {
//		X        string `json:"x"`
//		Y        string `json:"X"`
//		Expected string `json:"expected"`
//	}
//	type Test struct {
//		TestType []TestCase `json:"Swap"`
//	}
//	file, err := ioutil.ReadFile("../../../test/test-tables/sample_slip.json")
//	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
//	assert.NoError(t, err)
//	var test Test
//	err = json.Unmarshal(file, &test)
//	assert.NoError(t, err)
//	testcases := test.TestType
//	for _, test := range testcases {
//		CalculatePoolUnits()
//	}
//}
//
//func TestTradeSlip(t *testing.T) {
//	type TestCase struct {
//		X        string `json:"x"`
//		Y        string `json:"X"`
//		Expected string `json:"expected"`
//	}
//	type Test struct {
//		TestType []TestCase `json:"Swap"`
//	}
//	file, err := ioutil.ReadFile("../../../test/test-tables/sample_trade_slip.json")
//	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
//	assert.NoError(t, err)
//	var test Test
//	err = json.Unmarshal(file, &test)
//	assert.NoError(t, err)
//	testcases := test.TestType
//	for _, test := range testcases {
//		CalculatePoolUnits()
//	}
//}
//
//func TestLiquidityFee(t *testing.T) {
//	type TestCase struct {
//		X        string `json:"x"`
//		Y        string `json:"X"`
//		Expected string `json:"expected"`
//	}
//	type Test struct {
//		TestType []TestCase `json:"Swap"`
//	}
//	file, err := ioutil.ReadFile("../../../test/test-tables/sample_liquidity_fee.json")
//	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
//	assert.NoError(t, err)
//	var test Test
//	err = json.Unmarshal(file, &test)
//	assert.NoError(t, err)
//	testcases := test.TestType
//	for _, test := range testcases {
//		CalculatePoolUnits()
//	}
//}

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
