package keeper_test

import (
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tenderminttypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

var multiplierDec1 = sdk.OneDec()
var multiplierDec2 = sdk.MustNewDecFromStr("2.0")
var multiplierDec3 = sdk.MustNewDecFromStr("3.0")
var allocations = sdk.NewUintFromString("300")

//Function to Create 3 Pools with different Pool Depths and Different Multipliers
func createPool(keeper clpkeeper.Keeper, ctx sdk.Context, t *testing.T) {
	err := keeper.SetPool(ctx, &types.Pool{
		ExternalAsset:                 &types.Asset{Symbol: "atom"},
		NativeAssetBalance:            sdk.NewUint(1100),
		ExternalAssetBalance:          sdk.NewUint(1100),
		PoolUnits:                     sdk.NewUint(1100),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
	})
	require.NoError(t, err)
	err = keeper.SetPool(ctx, &types.Pool{
		ExternalAsset:                 &types.Asset{Symbol: "cusdc"},
		NativeAssetBalance:            sdk.NewUint(1000),
		ExternalAssetBalance:          sdk.NewUint(1000),
		PoolUnits:                     sdk.NewUint(1000),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
	})
	require.NoError(t, err)
	err = keeper.SetPool(ctx, &types.Pool{
		ExternalAsset:                 &types.Asset{Symbol: "ceth"},
		NativeAssetBalance:            sdk.NewUint(1000),
		ExternalAssetBalance:          sdk.NewUint(1000),
		PoolUnits:                     sdk.NewUint(1000),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
	})
	require.NoError(t, err)
}

//Function to define a reward period with start and end blocks, total allocation and setting pools with their respective multipliers
func generateRewardDistribution(keeper clpkeeper.Keeper, ctx sdk.Context) {
	params := keeper.GetRewardsParams(ctx)
	params.RewardPeriods = []*types.RewardPeriod{
		{RewardPeriodId: "Test 1", RewardPeriodStartBlock: 1, RewardPeriodEndBlock: 10,
			RewardPeriodAllocation: &allocations, RewardPeriodDefaultMultiplier: &multiplierDec1, RewardPeriodPoolMultipliers: []*types.PoolMultiplier{
				{PoolMultiplierAsset: "atom", Multiplier: &multiplierDec3},
				{PoolMultiplierAsset: "cusdc", Multiplier: &multiplierDec2},
				{PoolMultiplierAsset: "ceth", Multiplier: &multiplierDec1},
			},
		}}

	keeper.SetRewardParams(ctx, params)
}

//Function to validate pool rewards distribution at the end of the reward period
func validator(keeper clpkeeper.Keeper, ctx sdk.Context, t *testing.T) {
	pool, err := keeper.GetPool(ctx, "atom")
	require.NoError(t, err)
	require.EqualValues(t, "1253", pool.NativeAssetBalance.String())

	pool, err = keeper.GetPool(ctx, "cusdc")
	require.NoError(t, err)
	require.EqualValues(t, "1090", pool.NativeAssetBalance.String())

	pool, err = keeper.GetPool(ctx, "ceth")
	require.NoError(t, err)
	require.EqualValues(t, "1040", pool.NativeAssetBalance.String())
}

//TESTCASE#1:Test Function that calls rewards distribution function and validate the distribution against expected rewards values for a valid reward period
func TestDistributeDepthRewards(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	generateRewardDistribution(app.ClpKeeper, ctx)
	createPool(app.ClpKeeper, ctx, t)
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	validator(app.ClpKeeper, ctx, t)
}

//TESTCASE#2:Test Function that calls rewards distribution function for an invalid rewards period and check if it panics
func TestRewardDistributionInvalidPeriod(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	generateRewardDistribution(app.ClpKeeper, ctx)
	createPool(app.ClpKeeper, ctx, t)
	for block := 11; block <= 15; block++ {
		require.Panics(t, func() {
			app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
			app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
			app.Commit()
		})
	}
}

//TESTCASE#3:Test Function that calls rewards distribution function and validate the distribution against expected rewards values for a mix of valid reward period and invalid rewards period
func TestRewardDistributionInvalidExtendedPeriod(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	generateRewardDistribution(app.ClpKeeper, ctx)
	createPool(app.ClpKeeper, ctx, t)
	for block := 1; block <= 15; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	validator(app.ClpKeeper, ctx, t)
}

//TESTCASE#4:Test Function that calls rewards distribution function and validate the distribution against expected rewards values for multiple Rewards Periods
func TestRewardDistributionMultiplePeriods(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	params := app.ClpKeeper.GetRewardsParams(ctx)
	params.RewardPeriods = []*types.RewardPeriod{
		{RewardPeriodId: "Test 1", RewardPeriodStartBlock: 1, RewardPeriodEndBlock: 5,
			RewardPeriodAllocation: &allocations, RewardPeriodDefaultMultiplier: &multiplierDec1, RewardPeriodPoolMultipliers: []*types.PoolMultiplier{
				{PoolMultiplierAsset: "atom", Multiplier: &multiplierDec3},
				{PoolMultiplierAsset: "cusdc", Multiplier: &multiplierDec2},
				{PoolMultiplierAsset: "ceth", Multiplier: &multiplierDec1},
			},
		},
		{RewardPeriodId: "Test 2", RewardPeriodStartBlock: 6, RewardPeriodEndBlock: 10,
			RewardPeriodAllocation: &allocations, RewardPeriodDefaultMultiplier: &multiplierDec1, RewardPeriodPoolMultipliers: []*types.PoolMultiplier{
				{PoolMultiplierAsset: "atom", Multiplier: &multiplierDec3},
				{PoolMultiplierAsset: "cusdc", Multiplier: &multiplierDec2},
				{PoolMultiplierAsset: "ceth", Multiplier: &multiplierDec1},
			},
		},
	}
	app.ClpKeeper.SetRewardParams(ctx, params)
	createPool(app.ClpKeeper, ctx, t)
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	pool, err := app.ClpKeeper.GetPool(ctx, "atom")
	require.NoError(t, err)
	require.EqualValues(t, "1416", pool.NativeAssetBalance.String())

	pool, err = app.ClpKeeper.GetPool(ctx, "ceth")
	require.NoError(t, err)
	require.EqualValues(t, "1085", pool.NativeAssetBalance.String())

	pool, err = app.ClpKeeper.GetPool(ctx, "cusdc")
	require.NoError(t, err)
	require.EqualValues(t, "1181", pool.NativeAssetBalance.String())
}

//Test Function that calls unlock liquidity function and checks various scenarios
func TestUnlockedLiquidity(t *testing.T) {
	testcases := []struct {
		name                 string
		height               int64
		use                  sdk.Uint
		any                  bool
		availableUnits       sdk.Uint
		remainingUnlockUnits sdk.Uint
		unlocks              []*types.LiquidityUnlock
		expected             error
	}{
		//TESTCASE#5: Test Unlockliquidity - When no unlocks requested
		{
			name:                 "No unlocks",
			height:               1,
			use:                  sdk.NewUint(1000000),
			availableUnits:       sdk.NewUint(1000000),
			remainingUnlockUnits: sdk.NewUint(0),
			expected:             types.ErrBalanceNotAvailable,
		},
		//TESTCASE#6: Test Unlockliquidity - When locking period is not yet valid
		{
			name:                 "Unlock not ready",
			height:               5,
			use:                  sdk.NewUint(1000000),
			availableUnits:       sdk.NewUint(1000000),
			remainingUnlockUnits: sdk.NewUint(1000000),
			expected:             types.ErrBalanceNotAvailable,
			unlocks: []*types.LiquidityUnlock{
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1000000),
				},
			},
		},
		//TESTCASE#7: Test Unlockliquidity - When locking period is not yet valid but flagged to true to cancel locking for allowing unlock
		{
			name:                 "Unlock not ready but flag true",
			height:               5,
			any:                  true,
			use:                  sdk.NewUint(1000000),
			availableUnits:       sdk.NewUint(1000000),
			remainingUnlockUnits: sdk.NewUint(1000000),
			expected:             nil,
			unlocks: []*types.LiquidityUnlock{
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1000000),
				},
			},
		},

		//TESTCASE#8: Test Unlockliquidity - When locking period is valid but there is no liquidity or balance
		{
			name:                 "Insufficient balance/liquidity",
			height:               11,
			use:                  sdk.NewUint(1000000),
			availableUnits:       sdk.NewUint(10000),
			remainingUnlockUnits: sdk.NewUint(1000000),
			expected:             nil,
			unlocks: []*types.LiquidityUnlock{
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1000000),
				},
			},
		},
		//TESTCASE#9: Test Unlockliquidity - When locking period is valid and there is sufficient liquidity
		{
			name:                 "Test unlock liquidity with sufficient balance to unlock",
			height:               11,
			use:                  sdk.NewUint(2000000),
			availableUnits:       sdk.NewUint(3000000),
			remainingUnlockUnits: sdk.NewUint(0),
			expected:             types.ErrBalanceNotAvailable,
			unlocks: []*types.LiquidityUnlock{
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(800000),
				},
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1200000),
				},
			},
		},
		//TESTCASE#10: Test Unlockliquidity - When locking period is valid and some of funds requested for unlock are available but not all
		//or partial liquidity exists
		{
			name:                 "Test unlock liquidity with insufficient/partial balance available to unlock",
			height:               11,
			use:                  sdk.NewUint(3000000),
			availableUnits:       sdk.NewUint(2000000),
			remainingUnlockUnits: sdk.NewUint(1000000),
			expected:             nil,
			unlocks: []*types.LiquidityUnlock{
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(1000000),
				},
				{
					RequestHeight: 1,
					Units:         sdk.NewUint(2000000),
				},
			},
		},
	}
	for _, testcase := range testcases {
		tc := testcase
		t.Run(tc.name, func(t *testing.T) {
			app, ctx := test.CreateTestApp(false)
			ctx = ctx.WithBlockHeight(tc.height)
			params := app.ClpKeeper.GetRewardsParams(ctx)
			params.LiquidityRemovalLockPeriod = 10
			params.LiquidityRemovalCancelPeriod = 5
			app.ClpKeeper.SetRewardParams(ctx, params)
			lp := types.LiquidityProvider{
				Asset:                    &types.Asset{Symbol: "atom"},
				LiquidityProviderUnits:   sdk.NewUint(1000),
				LiquidityProviderAddress: "sif123",
				Unlocks:                  tc.unlocks,
			}
			app.ClpKeeper.SetLiquidityProvider(ctx, &lp)
			err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, tc.availableUnits, tc.any)
			require.ErrorIs(t, err, tc.expected)
		})
	}
}
