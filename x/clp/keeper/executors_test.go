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

//func TestCalculatePoolUnits(t *testing.T) {
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
//  testcases := test.TestType
//  for _,test := range testcases{
//		// r = native asset added;
//		// a = external asset added
//		// R = native Balance (before)
//		// A = external Balance (before)
//		// P = existing Pool Units
//  	CalculatePoolUnits()
//	}
//}

func TestCalculatePoolUnits(t *testing.T) {
	type TestCase struct {
		R        string `json:"r"`
		A        string `json:"a"`
		RR       string `json:"R"`
		AA       string `json:"A"`
		P        string `json:"P"`
		Expected string `json:"expected"`
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
			sdk.NewUintFromString(test.P),
			sdk.NewUintFromString(test.RR),
			sdk.NewUintFromString(test.AA),
			sdk.NewUintFromString(test.R),
			sdk.NewUintFromString(test.A))
		assert.NoError(t, err)
		assert.Equal(t, newPoolUnit, sdk.NewUintFromString(test.Expected))
		if !newPoolUnit.Equal(sdk.NewUintFromString(test.Expected)) {
			fmt.Printf("%s", test)
		}
	}
}
