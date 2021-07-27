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
	whitelisttypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func getDenomWhiteListEntries() whitelisttypes.Registry {
	return whitelisttypes.DefaultWhitelist()
}

func createTestAppForTestTables() (sdk.Context, *sifapp.SifchainApp) {
	wl := getDenomWhiteListEntries()
	ctx, app := clptest.CreateTestAppClp(false)
	for _, entry := range wl.Entries {
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
		for _, entry := range wl.Entries {
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
		for _, entry := range wl.Entries {
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
		for _, entry := range wl.Entries {
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
		for _, entry := range wl.Entries {
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
