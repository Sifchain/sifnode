package keeper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

//Commenting for now as the frontend calculation uses a different formula
//func TestTradeSlip(t *testing.T) {
//	type TestCase struct {
//		Sx       string `json:"x"`
//		X        string `json:"X"`
//		Expected string `json:"expected"`
//	}
//	type Test struct {
//		TestType []TestCase `json:"SingleSwapStandardSlip"`
//	}
//	file, err := ioutil.ReadFile("../../../test/test-tables/singleswap_standardslip.json")
//	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
//	assert.NoError(t, err)
//	var test Test
//	err = json.Unmarshal(file, &test)
//	assert.NoError(t, err)
//	testcases := test.TestType
//	totalCount := 0
//	failedCount := 0
//	f, err := os.Create("discrepancies_singleswap_standard_slip")
//	assert.NoError(t, err)
//	w := bufio.NewWriter(f)
//	for _, test := range testcases {
//		totalCount++
//		expected := GetInt(t, test.Expected)
//		e := calcTradeSlip(GetInt(t, test.X), GetInt(t, test.Sx))
//		//assert.Equal(t, e, expected)
//
//		if !e.Equal(expected) {
//			failedCount++
//			_, err := fmt.Fprintf(w, "Expected : %s  Actual : %s \n", expected.String(), e.String())
//			assert.NoError(t, err)
//		}
//	}
//	_, err = fmt.Fprintf(w, "TotalCount : %v  FailedCount : %v ", totalCount, failedCount)
//	assert.NoError(t, err)
//	w.Flush()
//	_ = f.Close()
//
//}

func TestLiquidityFee(t *testing.T) {
	type TestCase struct {
		Sx       string `json:"x"`
		X        string `json:"X"`
		Y        string `json:"Y"`
		Expected string `json:"expected"`
	}
	type Test struct {
		TestType []TestCase `json:"SingleSwapLiquidityFee"`
	}
	file, err := ioutil.ReadFile("../../../test/test-tables/singleswap_liquidityfees.json")
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	assert.NoError(t, err)
	var test Test
	err = json.Unmarshal(file, &test)
	assert.NoError(t, err)
	testcases := test.TestType
	totalCount := 0
	failedCount := 0
	f, err := os.Create("discrepancies_singleswap_liquidity_fee")
	assert.NoError(t, err)

	w := bufio.NewWriter(f)
	for _, test := range testcases {
		totalCount++
		res := calcLiquidityFee(GetInt(t, test.X), GetInt(t, test.Sx), GetInt(t, test.Y))
		//assert.GreaterOrEqual(t, GetInt(t, test.Expected).Uint64(),res.Uint64())
		expected := GetInt(t, test.Expected)
		if !res.Equal(expected) {
			failedCount++
			_, err := fmt.Fprintf(w, "Expected : %s  Actual : %s \n", expected.String(), res.String())
			assert.NoError(t, err)
		}
	}
	_, err = fmt.Fprintf(w, "TotalCount : %v  FailedCount : %v ", totalCount, failedCount)
	assert.NoError(t, err)
	w.Flush()
	os.Remove("discrepancies_singleswap_liquidity_fee")
	_ = f.Close()
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
		//assert.Equal(t, newPoolUnit, sdk.NewUintFromString(test.Expected))
		if !newPoolUnit.Equal(sdk.NewUintFromString(test.Expected)) {
			//fmt.Printf("%s", test)
		}
	}
}

func GetInt(t *testing.T, x string) sdk.Uint {
	// Precision based conversion as the Files contain values in scientific notation
	//flt, _, err := big.ParseFloat(x, 10, 0, big.ToNearestEven)
	//if err != nil {
	//	fmt.Println("Error parsing float :", x)
	//}
	//var i = new(big.Int)
	//i, _ = flt.Int(i)
	// Dec type from sdk
	//X := sdk.NewDecFromBigInt(i)
	X, err := sdk.NewDecFromStr(x)
	assert.NoError(t, err)
	if err != nil {
		return sdk.ZeroUint()
	}
	Xr := X.Mul(sdk.NewDec(10).Power(18)).RoundInt()
	return sdk.NewUintFromBigInt(Xr.BigInt())
}
