package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetSwapFeeRate(t *testing.T) {

	testcases := []struct {
		name                string
		asset               types.Asset
		swapFeeParams       types.SwapFeeParams
		marginEnabled       bool
		expectedSwapFeeRate sdk.Dec
	}{
		{
			name:                "empty token params",
			asset:               types.NewAsset("ceth"),
			swapFeeParams:       types.SwapFeeParams{DefaultSwapFeeRate: sdk.NewDecWithPrec(3, 3)},
			marginEnabled:       false,
			expectedSwapFeeRate: sdk.NewDecWithPrec(3, 3),
		},
		{
			name:  "match",
			asset: types.NewAsset("ceth"),
			swapFeeParams: types.SwapFeeParams{
				DefaultSwapFeeRate: sdk.NewDecWithPrec(3, 3),
				TokenParams: []*types.SwapFeeTokenParams{
					{
						Asset:       "ceth",
						SwapFeeRate: sdk.NewDecWithPrec(1, 3),
					},
					{
						Asset:       "cusdc",
						SwapFeeRate: sdk.NewDecWithPrec(2, 3),
					},
				},
			},
			marginEnabled:       false,
			expectedSwapFeeRate: sdk.NewDecWithPrec(1, 3),
		},
		{
			name:  "no match",
			asset: types.NewAsset("rowan"),
			swapFeeParams: types.SwapFeeParams{
				DefaultSwapFeeRate: sdk.NewDecWithPrec(3, 3),
				TokenParams: []*types.SwapFeeTokenParams{
					{
						Asset:       "ceth",
						SwapFeeRate: sdk.NewDecWithPrec(1, 3),
					},
					{
						Asset:       "cusdc",
						SwapFeeRate: sdk.NewDecWithPrec(2, 3),
					},
				},
			},
			marginEnabled:       false,
			expectedSwapFeeRate: sdk.NewDecWithPrec(3, 3),
		},
		{
			name:  "match but fallback to default rate as margin enabled",
			asset: types.NewAsset("ceth"),
			swapFeeParams: types.SwapFeeParams{
				DefaultSwapFeeRate: sdk.NewDecWithPrec(3, 3),
				TokenParams: []*types.SwapFeeTokenParams{
					{
						Asset:       "ceth",
						SwapFeeRate: sdk.NewDecWithPrec(1, 3),
					},
					{
						Asset:       "cusdc",
						SwapFeeRate: sdk.NewDecWithPrec(2, 3),
					},
				},
			},
			marginEnabled:       true,
			expectedSwapFeeRate: sdk.NewDecWithPrec(3, 3),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			ctx, app := test.CreateTestAppClp(false)

			app.ClpKeeper.SetSwapFeeParams(ctx, &tc.swapFeeParams)

			swapFeeRate := app.ClpKeeper.GetSwapFeeRate(ctx, tc.asset, tc.marginEnabled)

			require.Equal(t, tc.expectedSwapFeeRate.String(), swapFeeRate.String())
		})
	}
}
