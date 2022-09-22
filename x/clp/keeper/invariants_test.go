package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestBalanceModuleAccountCheck(t *testing.T) {
	tt := []struct {
		name           string
		Pools          []*types.Pool
		moduleBalances sdk.Coins
		stop           bool
	}{
		{
			name: "ok without margin",
			Pools: []*types.Pool{
				{
					ExternalAsset:        &types.Asset{Symbol: "ceth"},
					NativeAssetBalance:   sdk.NewUint(1000),
					ExternalAssetBalance: sdk.NewUint(2000),
					ExternalCustody:      sdk.ZeroUint(),
					NativeCustody:        sdk.ZeroUint(),
				},
			},
			moduleBalances: sdk.NewCoins(
				sdk.NewCoin("rowan", sdk.NewInt(1000)),
				sdk.NewCoin("ceth", sdk.NewInt(2000)),
			),
			stop: false,
		},
		{
			name: "ok with margin",
			Pools: []*types.Pool{
				{
					ExternalAsset:        &types.Asset{Symbol: "ceth"},
					NativeAssetBalance:   sdk.NewUint(1000),
					ExternalAssetBalance: sdk.NewUint(2000),
					ExternalCustody:      sdk.NewUint(1000),
					NativeCustody:        sdk.NewUint(500),
				},
			},
			moduleBalances: sdk.NewCoins(
				sdk.NewCoin("rowan", sdk.NewInt(1500)),
				sdk.NewCoin("ceth", sdk.NewInt(3000)),
			),
			stop: false,
		},
		{
			name: "native balance not ok without margin",
			Pools: []*types.Pool{
				{
					ExternalAsset:        &types.Asset{Symbol: "ceth"},
					NativeAssetBalance:   sdk.NewUint(1000),
					ExternalAssetBalance: sdk.NewUint(2000),
					ExternalCustody:      sdk.ZeroUint(),
					NativeCustody:        sdk.ZeroUint(),
				},
			},
			moduleBalances: sdk.NewCoins(
				sdk.NewCoin("rowan", sdk.NewInt(9000)),
				sdk.NewCoin("ceth", sdk.NewInt(2000)),
			),
			stop: true,
		},
		{
			name: "external balance not ok without margin",
			Pools: []*types.Pool{
				{
					ExternalAsset:        &types.Asset{Symbol: "ceth"},
					NativeAssetBalance:   sdk.NewUint(1000),
					ExternalAssetBalance: sdk.NewUint(2000),
					ExternalCustody:      sdk.ZeroUint(),
					NativeCustody:        sdk.ZeroUint(),
				},
			},
			moduleBalances: sdk.NewCoins(
				sdk.NewCoin("rowan", sdk.NewInt(1000)),
				sdk.NewCoin("ceth", sdk.NewInt(500)),
			),
			stop: true,
		},
		{
			name: "native balance not ok with margin",
			Pools: []*types.Pool{
				{
					ExternalAsset:        &types.Asset{Symbol: "ceth"},
					NativeAssetBalance:   sdk.NewUint(1000),
					ExternalAssetBalance: sdk.NewUint(2000),
					ExternalCustody:      sdk.NewUint(1000),
					NativeCustody:        sdk.NewUint(500),
				},
			},
			moduleBalances: sdk.NewCoins(
				sdk.NewCoin("rowan", sdk.NewInt(1000)),
				sdk.NewCoin("ceth", sdk.NewInt(3000)),
			),
			stop: true,
		},
		{
			name: "external balance not ok with margin",
			Pools: []*types.Pool{
				{
					ExternalAsset:        &types.Asset{Symbol: "ceth"},
					NativeAssetBalance:   sdk.NewUint(1000),
					ExternalAssetBalance: sdk.NewUint(2000),
					ExternalCustody:      sdk.NewUint(1000),
					NativeCustody:        sdk.NewUint(500),
				},
			},
			moduleBalances: sdk.NewCoins(
				sdk.NewCoin("rowan", sdk.NewInt(1500)),
				sdk.NewCoin("ceth", sdk.NewInt(5000)),
			),
			stop: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app, ctx := test.CreateTestApp(false)

			for _, pool := range tc.Pools {
				err := app.ClpKeeper.SetPool(ctx, pool)
				require.NoError(t, err)
			}

			err := app.BankKeeper.MintCoins(ctx, types.ModuleName, tc.moduleBalances)
			require.NoError(t, err)

			_, stop := app.ClpKeeper.BalanceModuleAccountCheck()(ctx)
			require.Equal(t, tc.stop, stop)
		})

	}
}

func TestUnitsCheck(t *testing.T) {
	tt := []struct {
		name  string
		Pools []*types.Pool
		lps   []*types.LiquidityProvider
		stop  bool
	}{
		{
			name: "ok",
			Pools: []*types.Pool{
				{
					ExternalAsset: &types.Asset{Symbol: "ceth"},
					PoolUnits:     sdk.NewUint(2000),
				},
			},
			lps: []*types.LiquidityProvider{
				{
					Asset:                    &types.Asset{Symbol: "ceth"},
					LiquidityProviderUnits:   sdk.NewUint(1000),
					LiquidityProviderAddress: "sif123",
				},
				{
					Asset:                    &types.Asset{Symbol: "ceth"},
					LiquidityProviderUnits:   sdk.NewUint(1000),
					LiquidityProviderAddress: "sif456",
				},
			},
			stop: false,
		},
		{
			name: "not ok",
			Pools: []*types.Pool{
				{
					ExternalAsset: &types.Asset{Symbol: "ceth"},
					PoolUnits:     sdk.NewUint(3000),
				},
			},
			lps: []*types.LiquidityProvider{
				{
					Asset:                    &types.Asset{Symbol: "ceth"},
					LiquidityProviderUnits:   sdk.NewUint(1000),
					LiquidityProviderAddress: "sif123",
				},
				{
					Asset:                    &types.Asset{Symbol: "ceth"},
					LiquidityProviderUnits:   sdk.NewUint(1000),
					LiquidityProviderAddress: "sif456",
				},
			},
			stop: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app, ctx := test.CreateTestApp(false)

			for _, pool := range tc.Pools {
				err := app.ClpKeeper.SetPool(ctx, pool)
				require.NoError(t, err)
			}

			for _, lp := range tc.lps {
				app.ClpKeeper.SetLiquidityProvider(ctx, lp)
			}

			_, stop := app.ClpKeeper.UnitsCheck()(ctx)
			require.Equal(t, tc.stop, stop)
		})

	}
}
