package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
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
		if !stakeUnits.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("Pool_Units | Expected : %s | Got : %s \n", test.Expected, stakeUnits.String())
			errCount++
		}
	}
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
		Yy, _ := calcSwapResult("cusdt",
			true,
			sdk.NewUintFromString(test.X),
			sdk.NewUintFromString(test.Xx),
			sdk.NewUintFromString(test.Y))
		if !Yy.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("SingleSwap-Result | Expected : %s | Got : %s \n", test.Expected, Yy.String())
			errCount++
		}
	}
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
		Yy, _ := calcLiquidityFee("ceth",
			true,
			sdk.NewUintFromString(test.X),
			sdk.NewUintFromString(test.Xx),
			sdk.NewUintFromString(test.Y))
		if !Yy.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("SingleSwap-Liquidityfees | Expected : %s | Got : %s \n", test.Expected, Yy.String())
			errCount++
		}
	}
}

func TestCalculateDoubleSwapResult(t *testing.T) {
	type TestCase struct {
		Ax       string `json:"ax"`
		AX       string `json:"aX"`
		AY       string `json:"aY"`
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
		Ay, _ := calcSwapResult("cusdt",
			true,
			sdk.NewUintFromString(test.AX),
			sdk.NewUintFromString(test.Ax),
			sdk.NewUintFromString(test.AY))

		By, _ := calcSwapResult("cusdt",
			true,
			sdk.NewUintFromString(test.BX),
			Ay,
			sdk.NewUintFromString(test.BY))

		if !By.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("Doubleswap_Result | Expected : %s | Got : %s \n", test.Expected, By.String())
			errCount++
		}
	}
}
