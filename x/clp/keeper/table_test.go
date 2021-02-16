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

func TestCalculatePoolUnits(t *testing.T) {
	//type TestCase struct {
	//	NativeAdded      string `json:"r"`
	//	ExternalAdded    string `json:"a"`
	//	NativeBalance    string `json:"R"`
	//	ExternalBalance  string `json:"A"`
	//	PoolUnitsBalance string `json:"P"`
	//	Expected         string `json:"expected"`
	//}
	//type Test struct {
	//	TestType []TestCase `json:"PoolUnits"`
	//}
	//file, err := ioutil.ReadFile("../../../test/test-tables/pool_units.json")
	//assert.NoError(t, err)
	//file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	//var test Test
	//err = json.Unmarshal(file, &test)
	//assert.NoError(t, err)
	//testcases := test.TestType
	//errCount := 0
	//for _, test := range testcases {
	//	_, stakeUnits, _ := CalculatePoolUnits(
	//		"cusdt",
	//		sdk.NewUintFromString(test.PoolUnitsBalance),
	//		sdk.NewUintFromString(test.NativeBalance),
	//		sdk.NewUintFromString(test.ExternalBalance),
	//		sdk.NewUintFromString(test.NativeAdded),
	//		sdk.NewUintFromString(test.ExternalAdded),
	//	)
	//	//assert.NoError(t, err)
	//	if test.Expected != "0" && !stakeUnits.Equal(sdk.NewUintFromString(test.Expected)) {
	//		errCount++
	//		fmt.Printf("Got %s , Expected %s \n", stakeUnits, test.Expected)
	//		//fmt.Printf("%+v \n", test)
	//
	//	}
	//
	//}
	//fmt.Println("Error Count :", errCount, "Total :", len(testcases))
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
	//for _, test := range testcases {
	//	res, _ := calcSwapResult("cusdt",
	//		true,//100000000000000000000
	//		sdk.NewUintFromString("500000000"),//(strings.Split(test.X,".")[0]),
	//		sdk.NewUintFromString("10000000"),//(strings.Split(test.Xx,".")[0]),
	//		sdk.NewUintFromString("100000000000000000000"))//(strings.Split(test.Y,".")[0]))
	//	//assert.NoError(t, err)
	//	//if test.Expected != "0" && !res.Equal(sdk.NewUintFromString(strings.Split(test.Expected,".")[0])) {
	//	//	errCount++
	//	//	fmt.Printf("Got %s , Expected %s \n", res, strings.Split(test.Expected,".")[0])
	//	//	//fmt.Printf("%+v \n", test)
	//	//
	//	//}
	//	fmt.Println(res)
	//   //1922337562475971000
	//   res2, _ := calcSwapResult("ceth",
	//	false,//100000000000000000000
	//	sdk.NewUintFromString("101922337562475971000"),//(strings.Split(test.X,".")[0]),
	//	res,//(strings.Split(test.Xx,".")[0]),
	//	sdk.NewUintFromString("400000000000000000000"))
	//   fmt.Println(res2)

	res, _ := calcSwapResult("ceth",
		true, //100000000000000000000
		sdk.NewUintFromString("70000000000000000000"),  //(strings.Split(test.X,".")[0]),
		sdk.NewUintFromString("14000000000000000000"),  //(strings.Split(test.Xx,".")[0]),
		sdk.NewUintFromString("100000000000000000000")) //(strings.Split(test.Y,".")[0]))
	fmt.Println(res)
	res2, _ := calcSwapResult("cusdt",
		false, //100000000000000000000
		sdk.NewUintFromString("113888888888888888890"), //(strings.Split(test.X,".")[0]),
		res, //(strings.Split(test.Xx,".")[0]),
		sdk.NewUintFromString("2500000000"))
	fmt.Println(res2)
	//	}
	fmt.Println("Error Count :", errCount, "Total :", len(testcases))
}
