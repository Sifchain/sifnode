package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetMinSwapFee(t *testing.T) {

	testcases := []struct {
		name               string
		asset              types.Asset
		tokenParams        []*types.SwapFeeTokenParams
		expectedMinSwapFee sdk.Uint
	}{
		{
			name:               "empty token params",
			asset:              types.NewAsset("ceth"),
			expectedMinSwapFee: sdk.ZeroUint(),
		},
		{
			name:               "match",
			asset:              types.NewAsset("ceth"),
			tokenParams:        []*types.SwapFeeTokenParams{{Asset: "ceth", MinSwapFee: sdk.NewUint(100)}, {Asset: "cusdc", MinSwapFee: sdk.NewUint(300)}},
			expectedMinSwapFee: sdk.NewUint(100),
		},
		{
			name:               "no match",
			asset:              types.NewAsset("rowan"),
			tokenParams:        []*types.SwapFeeTokenParams{{Asset: "ceth", MinSwapFee: sdk.NewUint(100)}, {Asset: "cusdc", MinSwapFee: sdk.NewUint(300)}},
			expectedMinSwapFee: sdk.ZeroUint(),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			minSwapFee := keeper.GetMinSwapFee(tc.asset, tc.tokenParams)

			require.Equal(t, tc.expectedMinSwapFee.String(), minSwapFee.String())
		})
	}
}
