package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

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
	file, err := ioutil.ReadFile("../../../test/test-tables/pool_units.json")
	assert.NoError(t, err)
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	errCount := 0
	for _, test := range testcases {
		_, stakeUnits, _ := CalculatePoolUnits(
			"cusdt",
			sdk.NewUintFromString(test.PoolUnitsBalance),
			sdk.NewUintFromString(test.NativeBalance),
			sdk.NewUintFromString(test.ExternalBalance),
			sdk.NewUintFromString(test.NativeAdded),
			sdk.NewUintFromString(test.ExternalAdded),
		)
		if test.Expected != "0" && !stakeUnits.Equal(sdk.NewUintFromString(test.Expected)) {
			errCount++
			fmt.Printf("Got %s , Expected %s \n", stakeUnits, test.Expected)
		}

	}
	fmt.Printf("Total/Failed: %d/%d", len(testcases), errCount)
}

func TestCalculateSwapResult(t *testing.T) {
	type TestCase struct {
		Xx       string `json:"x"`
		X        string `json:"X"`
		Y        string `json:"Y"`
		Expected string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"SingleSwapResult"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/singleswap_result.json")
	assert.NoError(t, err)
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	errCount := 0
	for _, test := range testcases {
		res, _ := calcSwapResult("cusdt",
			true, //100000000000000000000
			sdk.NewUintFromString(test.X),
			sdk.NewUintFromString(test.Xx),
			sdk.NewUintFromString(test.Y))
		if test.Expected != "0" && !res.Equal(sdk.NewUintFromString(test.Expected)) {
			errCount++
			fmt.Printf("Got %s , Expected %s \n", res, strings.Split(test.Expected, ".")[0])
		}
	}
	fmt.Printf("Total/Failed: %d/%d", len(testcases), errCount)
}

func TestCalculateSwapLiquidityFee(t *testing.T) {
	type TestCase struct {
		Xx       string `json:"x"`
		X        string `json:"X"`
		Y        string `json:"Y"`
		Expected string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"SingleSwapLiquidityFee"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/singleswap_liquidityfees.json")
	assert.NoError(t, err)
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	errCount := 0
	for _, test := range testcases {
		res, _ := calcLiquidityFee(
			sdk.NewUintFromString(test.X),
			sdk.NewUintFromString(test.Xx),
			sdk.NewUintFromString(test.Y))
		if test.Expected != "0" && !res.Equal(sdk.NewUintFromString(test.Expected)) {
			errCount++
			fmt.Printf("Got %s , Expected %s \n", res, strings.Split(test.Expected, ".")[0])
		}
	}
	fmt.Printf("Total/Failed: %d/%d", len(testcases), errCount)
}

func TestCalculateDoubleSwapResult(t *testing.T) {
	type TestCase struct {
		Xx       string `json:"ax"`
		X        string `json:"aX"`
		Y        string `json:"aY"`
		BX       string `json:"bX"`
		BY       string `json:"bY"`
		Expected string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"DoubleSwap"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/doubleswap_result.json")
	assert.NoError(t, err)
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	errCount := 0
	for _, test := range testcases {
		res, _ := calcSwapResult("cusdt",
			true, //100000000000000000000
			sdk.NewUintFromString(test.X),
			sdk.NewUintFromString(test.Xx),
			sdk.NewUintFromString(test.Y))

		res2, _ := calcSwapResult("cusdt",
			true, //100000000000000000000
			sdk.NewUintFromString(test.X),
			sdk.NewUintFromString(test.Xx),
			sdk.NewUintFromString(test.Y))

		if test.Expected != "0" && !res2.Equal(sdk.NewUintFromString(test.Expected)) {
			errCount++
			fmt.Printf("Got %s , Expected %s \n", res, strings.Split(test.Expected, ".")[0])
		}
	}
	fmt.Printf("Total/Failed: %d/%d", len(testcases), errCount)
}
