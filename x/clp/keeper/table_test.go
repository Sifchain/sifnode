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

// THe normalization has been shifted to a different module , so these tests cannot remain pure anymore , unless we add normalization factor to the json file directly
// Using a stub now for providing hardcoded values of normalization-factor and adjustment flag
// The stub implementation can be replaced with DenomWhitelist types

func GetNormalizationMap() map[string]int64 {
	m := make(map[string]int64)
	m["cel"] = 4
	m["ausdc"] = 6
	m["usdt"] = 6
	m["usdc"] = 6
	m["cro"] = 8
	m["cdai"] = 8
	m["wbtc"] = 8
	m["ceth"] = 8
	m["renbtc"] = 8
	m["cusdc"] = 8
	m["husd"] = 8
	m["ampl"] = 9
	return m
}
func GetNormalizationFactorStub(denom string) (sdk.Dec, bool) {
	normalizationFactor := sdk.NewDec(1)
	nf, ok := GetNormalizationMap()[denom[1:]]
	adjustExternalToken := false
	if ok {
		adjustExternalToken = true
		diffFactor := 18 - nf
		if diffFactor < 0 {
			diffFactor = nf - 18
			adjustExternalToken = false
		}
		normalizationFactor = sdk.NewDec(10).Power(uint64(diffFactor))
	}
	return normalizationFactor, adjustExternalToken
}

func TestCalculatePoolUnits(t *testing.T) {
	type TestCase struct {
		Symbol           string `json:"symbol"`
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
	file, err := ioutil.ReadFile("../../../test/test-tables/pool_units_after_upgrade.json")
	assert.NoError(t, err)

	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)

	testcases := test.TestType
	errCount := 0
	for _, test := range testcases {
		nf, ad := GetNormalizationFactorStub(test.Symbol)
		_, stakeUnits, _ := CalculatePoolUnits(
			sdk.NewUintFromString(test.PoolUnitsBalance),
			sdk.NewUintFromString(test.NativeBalance),
			sdk.NewUintFromString(test.ExternalBalance),
			sdk.NewUintFromString(test.NativeAdded),
			sdk.NewUintFromString(test.ExternalAdded),
			nf, ad,
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
		nf, ad := GetNormalizationFactorStub("cusdt")
		Yy, _ := calcSwapResult(
			true,
			nf,
			ad,
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
		nf, ad := GetNormalizationFactorStub("ceth")
		Yy, _ := calcLiquidityFee(
			true,
			nf, ad,
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
		nf, ad := GetNormalizationFactorStub("cusdt")
		Ay, _ := calcSwapResult(
			true,
			nf, ad,
			sdk.NewUintFromString(test.AX),
			sdk.NewUintFromString(test.Ax),
			sdk.NewUintFromString(test.AY))
		nf, ad = GetNormalizationFactorStub("cusdt")
		By, _ := calcSwapResult(
			true,
			nf, ad,
			sdk.NewUintFromString(test.BX),
			Ay,
			sdk.NewUintFromString(test.BY))

		if !By.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("Doubleswap_Result | Expected : %s | Got : %s \n", test.Expected, By.String())
			errCount++
		}
	}
}

func TestCalculatePoolUnitsAfterUpgrade(t *testing.T) {
	type TestCase struct {
		Symbol           string `json:"symbol"`
		NativeAdded      string `json:"r"`
		ExternalAdded    string `json:"a"`
		NativeBalance    string `json:"R"`
		ExternalBalance  string `json:"A"`
		PoolUnitsBalance string `json:"P"`
		Expected         string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"PoolUnitsAfterUpgrade"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/pool_units_after_upgrade.json")
	assert.NoError(t, err)

	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)

	testcases := test.TestType
	errCount := 0
	for _, test := range testcases {
		nf, ad := GetNormalizationFactorStub(test.Symbol)
		_, stakeUnits, _ := CalculatePoolUnits(
			sdk.NewUintFromString(test.PoolUnitsBalance),
			sdk.NewUintFromString(test.NativeBalance),
			sdk.NewUintFromString(test.ExternalBalance),
			sdk.NewUintFromString(test.NativeAdded),
			sdk.NewUintFromString(test.ExternalAdded),
			nf,
			ad,
		)
		if !stakeUnits.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("Pool_Units_After_Upgrade | Expected : %s | Got : %s \n", test.Expected, stakeUnits.String())
			errCount++
		}
	}
}
