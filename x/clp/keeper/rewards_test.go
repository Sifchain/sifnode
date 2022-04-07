package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tenderminttypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestEndBlock(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	// Setup reward period
	params := app.ClpKeeper.GetRewardsParams(ctx)
	allocation := sdk.NewUintFromString("200000000000000000000000000")
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
	}
	app.ClpKeeper.SetRewardParams(ctx, params)
	err := app.ClpKeeper.SetPool(ctx, &types.Pool{
		ExternalAsset:        &types.Asset{Symbol: "atom"},
		NativeAssetBalance:   sdk.NewUint(1000),
		ExternalAssetBalance: sdk.NewUint(1000),
		PoolUnits:            sdk.NewUint(1000),
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &types.Pool{
		ExternalAsset:        &types.Asset{Symbol: "cusdc"},
		NativeAssetBalance:   sdk.NewUint(1000),
		ExternalAssetBalance: sdk.NewUint(1000),
		PoolUnits:            sdk.NewUint(1000),
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &types.Pool{
		ExternalAsset:        &types.Asset{Symbol: "ceth"},
		NativeAssetBalance:   sdk.NewUint(1000),
		ExternalAssetBalance: sdk.NewUint(1000),
		PoolUnits:            sdk.NewUint(1000),
	})
	require.NoError(t, err)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.False(t, startingSupply.Equal(periodOneSupply), "starting: %s period: %s", startingSupply.String(), periodOneSupply.String())
	require.True(t, periodOneSupply.IsGTE(startingSupply))
	// check pool has expected increase
	pool, err := app.ClpKeeper.GetPool(ctx, "atom")
	require.NoError(t, err)
	require.Equal(t, "66666666666666666600001000", pool.NativeAssetBalance.String())
	expected := sdk.NewUintFromString("66666666666666666666667666")
	accuracy := sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()).Quo(sdk.NewDecFromBigInt(expected.BigInt()))
	require.True(t, accuracy.GT(sdk.MustNewDecFromStr("0.99")))
	// TODO continue through another portion of the period and ensure supply is increased.
	// continue through a non reward period
	for block := 11; block <= 20; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply is unchanged
	supplyCheck := app.BankKeeper.GetSupply(ctx, "rowan")
	//log.Printf("starting supply: %s final supply: %s after period one: %s", startingSupply.String(), supplyCheck.String(), periodOneSupply.String())
	require.True(t, supplyCheck.Equal(periodOneSupply))
}

func TestUseUnlockedLiquidity(t *testing.T) {
	tt := []struct {
		name     string
		height   int64
		use      sdk.Uint
		unlocks  []*types.LiquidityUnlock
		expected error
	}{
		{
			name:     "No unlocks",
			height:   1,
			use:      sdk.NewUint(1000),
			expected: types.ErrBalanceNotAvailable,
		}, {
			name:     "Unlock not ready",
			height:   5,
			use:      sdk.NewUint(1000),
			expected: types.ErrBalanceNotAvailable,
			unlocks: []*types.LiquidityUnlock{
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1000),
				},
			},
		},
		{
			name:     "Available via split",
			height:   50,
			use:      sdk.NewUint(2000),
			expected: nil,
			unlocks: []*types.LiquidityUnlock{
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1000),
				},
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1000),
				},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			app, ctx := test.CreateTestApp(false)
			ctx = ctx.WithBlockHeight(tc.height)
			params := app.ClpKeeper.GetRewardsParams(ctx)
			params.LiquidityRemovalLockPeriod = 10
			params.LiquidityRemovalCancelPeriod = 5
			app.ClpKeeper.SetRewardParams(ctx, params)
			lp := types.LiquidityProvider{
				Asset:                    &types.Asset{Symbol: "atom"},
				LiquidityProviderUnits:   sdk.NewUint(100),
				LiquidityProviderAddress: "sif123",
				Unlocks:                  tc.unlocks,
			}
			app.ClpKeeper.SetLiquidityProvider(ctx, &lp)
			err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, sdk.NewUint(1000))
			require.ErrorIs(t, err, tc.expected)
		})
	}
}
