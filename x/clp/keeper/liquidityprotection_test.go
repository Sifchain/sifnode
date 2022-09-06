package keeper_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetNativePrice(t *testing.T) {
	testcases := []struct {
		name                     string
		pricingAsset             string
		createPool               bool
		poolNativeAssetBalance   sdk.Uint
		poolExternalAssetBalance sdk.Uint
		pmtpCurrentRunningRate   sdk.Dec
		expectedPrice            sdk.Dec
		expectedError            error
	}{
		{
			name:          "success",
			pricingAsset:  types.NativeSymbol,
			expectedPrice: sdk.NewDec(1),
		},
		{
			name:          "fail pool does not exist",
			pricingAsset:  "usdc",
			expectedError: types.ErrMaxRowanLiquidityThresholdAssetPoolDoesNotExist,
		},
		{
			name:                     "success non rowan pricing asset",
			pricingAsset:             "usdc",
			createPool:               true,
			poolNativeAssetBalance:   sdk.NewUint(100000),
			poolExternalAssetBalance: sdk.NewUint(1000),
			pmtpCurrentRunningRate:   sdk.OneDec(),
			expectedPrice:            sdk.MustNewDecFromStr("0.02"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {

				if tc.createPool {
					pools := []*types.Pool{
						{
							ExternalAsset:        &types.Asset{Symbol: tc.pricingAsset},
							NativeAssetBalance:   tc.poolNativeAssetBalance,
							ExternalAssetBalance: tc.poolExternalAssetBalance,
						},
					}
					clpGs := types.DefaultGenesisState()

					clpGs.Params = types.Params{
						MinCreatePoolThreshold: 1,
					}
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ := app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			liquidityProtectionParams := app.ClpKeeper.GetLiquidityProtectionParams(ctx)
			liquidityProtectionParams.MaxRowanLiquidityThresholdAsset = tc.pricingAsset
			app.ClpKeeper.SetLiquidityProtectionParams(ctx, liquidityProtectionParams)
			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, tc.pmtpCurrentRunningRate)

			price, err := app.ClpKeeper.GetNativePrice(ctx)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedPrice.String(), price.String()) // compare strings so that the expected amounts can be read from the failure message
		})
	}
}

func TestKeeper_IsBlockedByLiquidityProtection(t *testing.T) {
	testcases := []struct {
		name                           string
		currentRowanLiquidityThreshold sdk.Uint
		nativeAmount                   sdk.Uint
		nativePrice                    sdk.Dec
		expectedIsBlocked              bool
	}{
		{
			name:                           "not blocked",
			currentRowanLiquidityThreshold: sdk.NewUint(180),
			nativeAmount:                   sdk.NewUint(900),
			nativePrice:                    sdk.MustNewDecFromStr("0.2"),
			expectedIsBlocked:              false,
		},
		{
			name:                           "blocked",
			currentRowanLiquidityThreshold: sdk.NewUint(179),
			nativeAmount:                   sdk.NewUint(900),
			nativePrice:                    sdk.MustNewDecFromStr("0.2"),
			expectedIsBlocked:              true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			app, ctx := test.CreateTestApp(false)

			liquidityProtectionRateParams := app.ClpKeeper.GetLiquidityProtectionRateParams(ctx)
			liquidityProtectionRateParams.CurrentRowanLiquidityThreshold = tc.currentRowanLiquidityThreshold
			app.ClpKeeper.SetLiquidityProtectionRateParams(ctx, liquidityProtectionRateParams)

			isBlocked := app.ClpKeeper.IsBlockedByLiquidityProtection(ctx, tc.nativeAmount, tc.nativePrice)

			require.Equal(t, tc.expectedIsBlocked, isBlocked)
		})
	}
}

func TestKeeper_MustUpdateLiquidityProtectionThreshold(t *testing.T) {
	testcases := []struct {
		name                           string
		maxRowanLiquidityThreshold     sdk.Uint
		currentRowanLiquidityThreshold sdk.Uint
		isActive                       bool
		nativeAmount                   sdk.Uint
		nativePrice                    sdk.Dec
		sellNative                     bool
		expectedUpdatedThreshold       sdk.Uint
		expectedPanicError             string
	}{
		{
			name:                           "sell native",
			maxRowanLiquidityThreshold:     sdk.NewUint(100000000),
			currentRowanLiquidityThreshold: sdk.NewUint(180),
			isActive:                       true,
			nativeAmount:                   sdk.NewUint(900),
			nativePrice:                    sdk.MustNewDecFromStr("0.2"),
			sellNative:                     true,
			expectedUpdatedThreshold:       sdk.ZeroUint(),
		},
		{
			name:                           "buy native",
			maxRowanLiquidityThreshold:     sdk.NewUint(100000000),
			currentRowanLiquidityThreshold: sdk.NewUint(180),
			isActive:                       true,
			nativeAmount:                   sdk.NewUint(900),
			nativePrice:                    sdk.MustNewDecFromStr("0.2"),
			sellNative:                     false,
			expectedUpdatedThreshold:       sdk.NewUint(360),
		},
		{
			name:                           "liquidity protection disabled",
			maxRowanLiquidityThreshold:     sdk.NewUint(100000000),
			currentRowanLiquidityThreshold: sdk.NewUint(180),
			isActive:                       false,
			expectedUpdatedThreshold:       sdk.NewUint(180),
		},
		{
			name:                           "panics if sell native value greater than current threshold",
			currentRowanLiquidityThreshold: sdk.NewUint(180),
			isActive:                       true,
			nativeAmount:                   sdk.NewUint(900),
			nativePrice:                    sdk.MustNewDecFromStr("1"),
			sellNative:                     true,
			expectedPanicError:             "expect sell native value to be less than currentRowanLiquidityThreshold",
		},
		{
			name:                           "does not exceed max",
			maxRowanLiquidityThreshold:     sdk.NewUint(150),
			currentRowanLiquidityThreshold: sdk.NewUint(100),
			isActive:                       true,
			nativeAmount:                   sdk.NewUint(80),
			nativePrice:                    sdk.MustNewDecFromStr("1"),
			sellNative:                     false,
			expectedUpdatedThreshold:       sdk.NewUint(150),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			app, ctx := test.CreateTestApp(false)

			liquidityProtectionParams := app.ClpKeeper.GetLiquidityProtectionParams(ctx)
			liquidityProtectionParams.IsActive = tc.isActive
			liquidityProtectionParams.MaxRowanLiquidityThreshold = tc.maxRowanLiquidityThreshold
			app.ClpKeeper.SetLiquidityProtectionParams(ctx, liquidityProtectionParams)

			liquidityProtectionRateParams := app.ClpKeeper.GetLiquidityProtectionRateParams(ctx)
			liquidityProtectionRateParams.CurrentRowanLiquidityThreshold = tc.currentRowanLiquidityThreshold
			app.ClpKeeper.SetLiquidityProtectionRateParams(ctx, liquidityProtectionRateParams)

			if tc.expectedPanicError != "" {
				require.PanicsWithError(t, tc.expectedPanicError, func() {
					app.ClpKeeper.MustUpdateLiquidityProtectionThreshold(ctx, tc.sellNative, tc.nativeAmount, tc.nativePrice)
				})
				return
			}

			app.ClpKeeper.MustUpdateLiquidityProtectionThreshold(ctx, tc.sellNative, tc.nativeAmount, tc.nativePrice)

			liquidityProtectionRateParams = app.ClpKeeper.GetLiquidityProtectionRateParams(ctx)

			require.Equal(t, tc.expectedUpdatedThreshold.String(), liquidityProtectionRateParams.CurrentRowanLiquidityThreshold.String())
		})
	}
}
