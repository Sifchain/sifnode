package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetAssetSwapFeeParams(t *testing.T) {

	testcases := []struct {
		name                string
		asset               types.Asset
		swapFeeParams       types.SwapFeeParams
		expectedSwapFeeRate sdk.Dec
		expectedMinSwapFee  sdk.Uint
	}{
		{
			name:                "empty token params",
			asset:               types.NewAsset("ceth"),
			swapFeeParams:       types.SwapFeeParams{DefaultSwapFeeRate: sdk.NewDecWithPrec(3, 3)},
			expectedSwapFeeRate: sdk.NewDecWithPrec(3, 3),
			expectedMinSwapFee:  sdk.ZeroUint(),
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
						MinSwapFee:  sdk.NewUint(100),
					},
					{
						Asset:       "cusdc",
						SwapFeeRate: sdk.NewDecWithPrec(2, 3),
						MinSwapFee:  sdk.NewUint(300),
					},
				},
			},
			expectedSwapFeeRate: sdk.NewDecWithPrec(1, 3),
			expectedMinSwapFee:  sdk.NewUint(100),
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
						MinSwapFee:  sdk.NewUint(100),
					},
					{
						Asset:       "cusdc",
						SwapFeeRate: sdk.NewDecWithPrec(2, 3),
						MinSwapFee:  sdk.NewUint(300),
					},
				},
			},
			expectedSwapFeeRate: sdk.NewDecWithPrec(3, 3),
			expectedMinSwapFee:  sdk.ZeroUint(),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			swapFeeRate, minSwapFee := keeper.GetAssetSwapFeeParams(tc.asset, &tc.swapFeeParams)

			require.Equal(t, tc.expectedSwapFeeRate.String(), swapFeeRate.String())
			require.Equal(t, tc.expectedMinSwapFee.String(), minSwapFee.String())
		})
	}
}
