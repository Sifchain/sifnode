package keeper_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	sifapp "github.com/Sifchain/sifnode/app"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptest "github.com/Sifchain/sifnode/x/clp/test"
	whitelisttypes "github.com/Sifchain/sifnode/x/whitelist/types"
)

func getDenomWhiteListEntries() whitelisttypes.DenomWhitelist {
	return whitelisttypes.DenomWhitelist{
		DenomWhitelistEntries: []*whitelisttypes.DenomWhitelistEntry{
			{Denom: "ccel", Decimals: 4},
			{Denom: "causc", Decimals: 6},
			{Denom: "cusdt", Decimals: 6},
			{Denom: "cusdc", Decimals: 6},
			{Denom: "ccro", Decimals: 8},
			{Denom: "ccdai", Decimals: 8},
			{Denom: "cwbtc", Decimals: 8},
			{Denom: "cceth", Decimals: 8},
			{Denom: "crenbtc", Decimals: 8},
			{Denom: "ccusdc", Decimals: 8},
			{Denom: "chusd", Decimals: 8},
			{Denom: "campl", Decimals: 9},
			{Denom: "ceth", Decimals: 18},
			{Denom: "cdai", Decimals: 18},
			{Denom: "cyfi", Decimals: 18},
			{Denom: "czrx", Decimals: 18},
			{Denom: "cwscrt", Decimals: 18},
			{Denom: "cwfil", Decimals: 18},
			{Denom: "cwbtc", Decimals: 18},
			{Denom: "cuni", Decimals: 18},
			{Denom: "cuma", Decimals: 18},
			{Denom: "ctusd", Decimals: 18},
			{Denom: "csxp", Decimals: 18},
			{Denom: "csushi", Decimals: 18},
			{Denom: "csusd", Decimals: 18},
			{Denom: "csrm", Decimals: 18},
			{Denom: "csnx", Decimals: 18},
			{Denom: "csand", Decimals: 18},
			{Denom: "crune", Decimals: 18},
			{Denom: "creef", Decimals: 18},
			{Denom: "cogn", Decimals: 18},
			{Denom: "cocean", Decimals: 18},
			{Denom: "cmana", Decimals: 18},
			{Denom: "clrc", Decimals: 18},
			{Denom: "clon", Decimals: 18},
			{Denom: "clink", Decimals: 18},
			{Denom: "ciotx", Decimals: 18},
			{Denom: "cgrt", Decimals: 18},
			{Denom: "cftm", Decimals: 18},
			{Denom: "cesd", Decimals: 18},
			{Denom: "cenj", Decimals: 18},
			{Denom: "ccream", Decimals: 18},
			{Denom: "ccomp", Decimals: 18},
			{Denom: "ccocos", Decimals: 18},
			{Denom: "cbond", Decimals: 18},
			{Denom: "cbnt", Decimals: 18},
			{Denom: "cbat", Decimals: 18},
			{Denom: "cband", Decimals: 18},
			{Denom: "cbal", Decimals: 18},
			{Denom: "cant", Decimals: 18},
			{Denom: "caave", Decimals: 18},
			{Denom: "c1inch", Decimals: 18},
		},
	}
}

func createTestAppForTestTables() (sdk.Context, *sifapp.SifchainApp) {
	wl := getDenomWhiteListEntries()
	ctx, app := clptest.CreateTestAppClp(false)
	for _, entry := range wl.DenomWhitelistEntries {
		app.WhitelistKeeper.SetDenom(ctx, entry.Denom, entry.Decimals)
	}
	return ctx, app
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

	ctx, app := createTestAppForTestTables()

	testcases := test.TestType
	for _, test := range testcases {
		wl := getDenomWhiteListEntries()
		for _, entry := range wl.DenomWhitelistEntries {
			nf, ad := app.ClpKeeper.GetNormalizationFactor(ctx, entry.Denom)
			_, stakeUnits, _ := clpkeeper.CalculatePoolUnits(
				sdk.NewUintFromString(test.PoolUnitsBalance),
				sdk.NewUintFromString(test.NativeBalance),
				sdk.NewUintFromString(test.ExternalBalance),
				sdk.NewUintFromString(test.NativeAdded),
				sdk.NewUintFromString(test.ExternalAdded),
				nf,
				ad,
			)
			assert.True(t, stakeUnits.Equal(sdk.NewUintFromString(test.Expected)))
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

	ctx, app := createTestAppForTestTables()

	testcases := test.TestType
	for _, test := range testcases {
		wl := getDenomWhiteListEntries()
		for _, entry := range wl.DenomWhitelistEntries {
			nf, ad := app.ClpKeeper.GetNormalizationFactor(ctx, entry.Denom)
			Yy, _ := clpkeeper.CalcSwapResult(
				true,
				nf,
				ad,
				sdk.NewUintFromString(test.X),
				sdk.NewUintFromString(test.Xx),
				sdk.NewUintFromString(test.Y),
			)
			assert.True(t, Yy.Equal(sdk.NewUintFromString(test.Expected)))
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

	ctx, app := createTestAppForTestTables()

	testcases := test.TestType
	for _, test := range testcases {
		wl := getDenomWhiteListEntries()
		for _, entry := range wl.DenomWhitelistEntries {
			nf, ad := app.ClpKeeper.GetNormalizationFactor(ctx, entry.Denom)
			Yy, _ := clpkeeper.CalcLiquidityFee(
				true,
				nf,
				ad,
				sdk.NewUintFromString(test.X),
				sdk.NewUintFromString(test.Xx),
				sdk.NewUintFromString(test.Y),
			)
			assert.True(t, Yy.Equal(sdk.NewUintFromString(test.Expected)))
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

	ctx, app := createTestAppForTestTables()

	testcases := test.TestType
	for _, test := range testcases {
		wl := getDenomWhiteListEntries()
		for _, entry := range wl.DenomWhitelistEntries {
			nf, ad := app.ClpKeeper.GetNormalizationFactor(ctx, entry.Denom)
			Ay, _ := clpkeeper.CalcSwapResult(
				true,
				nf,
				ad,
				sdk.NewUintFromString(test.AX),
				sdk.NewUintFromString(test.Ax),
				sdk.NewUintFromString(test.AY),
			)
			By, _ := clpkeeper.CalcSwapResult(
				false,
				nf,
				ad,
				sdk.NewUintFromString(test.BX),
				Ay,
				sdk.NewUintFromString(test.BY),
			)
			assert.True(t, By.Equal(sdk.NewUintFromString(test.Expected)))
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

	ctx, app := createTestAppForTestTables()

	testcases := test.TestType
	for _, test := range testcases {
		nf, ad := app.ClpKeeper.GetNormalizationFactor(ctx, test.Symbol)
		_, stakeUnits, _ := clpkeeper.CalculatePoolUnits(
			sdk.NewUintFromString(test.PoolUnitsBalance),
			sdk.NewUintFromString(test.NativeBalance),
			sdk.NewUintFromString(test.ExternalBalance),
			sdk.NewUintFromString(test.NativeAdded),
			sdk.NewUintFromString(test.ExternalAdded),
			nf,
			ad,
		)
		assert.True(t, stakeUnits.Equal(sdk.NewUintFromString(test.Expected)))
	}
}
