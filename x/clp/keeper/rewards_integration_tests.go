package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tenderminttypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

type unlockliquidity struct {
	name     string
	height   int64
	use      sdk.Uint
	unlocks  []*types.LiquidityUnlock
	expected error
}

//TestCase#1: Testing reward allocation for even number of pools with equal depth
func TestEndBlockWithEvenPoolsEqualDepth(t *testing.T) {
	fmt.Println("Testing TestEndBlockWithEvenPoolsEqualDepth")
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	var err error
	for i := 1; i <= 2; i++ {
		err = app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: "pool" + strconv.Itoa(i)},
			NativeAssetBalance:   sdk.NewUint(10000000),
			ExternalAssetBalance: sdk.NewUint(10000000),
			PoolUnits:            sdk.NewUint(10000000),
		})
		require.NoError(t, err)
	}
	allocation := sdk.NewUint(600000)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
	}
	app.ClpKeeper.SetParams(ctx, params)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.True(t, periodOneSupply.IsGTE(startingSupply))
	// check pool has expected increase
	for i := 1; i <= 2; i++ {
		pool, err := app.ClpKeeper.GetPool(ctx, "pool"+strconv.Itoa(i))
		require.NoError(t, err)
		require.Equal(t, "10300000", pool.NativeAssetBalance.String())
	}
}

//TestCase#2: Testing reward allocation for odd number of pools with equal depth
func TestEndBlockWithOddPoolsEqualDepth(t *testing.T) {
	fmt.Println("Testing TestEndBlockWithOddPoolsEqualDepth")
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	var err error
	for i := 1; i <= 3; i++ {
		err = app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: "pool" + strconv.Itoa(i)},
			NativeAssetBalance:   sdk.NewUint(10000000),
			ExternalAssetBalance: sdk.NewUint(10000000),
			PoolUnits:            sdk.NewUint(10000000),
		})
		require.NoError(t, err)
	}
	allocation := sdk.NewUint(600000)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
	}
	app.ClpKeeper.SetParams(ctx, params)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.True(t, periodOneSupply.IsGTE(startingSupply))
	// check pool has expected increase
	for i := 1; i <= 3; i++ {
		pool, err := app.ClpKeeper.GetPool(ctx, "pool"+strconv.Itoa(i))
		require.NoError(t, err)
		require.Equal(t, "10199990", pool.NativeAssetBalance.String())
	}

}

//TestCase#3: Testing reward allocation for even number of pools with different depth
func TestEndBlockWithEvenPoolsDifferentDepth(t *testing.T) {
	fmt.Println("Testing TestEndBlockWithEvenPoolsDifferentDepth")
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	var err error
	for i := 1; i <= 2; i++ {
		err = app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: "pool" + strconv.Itoa(i)},
			NativeAssetBalance:   sdk.NewUint(uint64(i * 1000000)),
			ExternalAssetBalance: sdk.NewUint(uint64(i * 1000000)),
			PoolUnits:            sdk.NewUint(uint64(i * 1000000)),
		})
		require.NoError(t, err)
	}
	allocation := sdk.NewUint(600000)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
	}
	app.ClpKeeper.SetParams(ctx, params)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.True(t, periodOneSupply.IsGTE(startingSupply))

	pool1, _ := app.ClpKeeper.GetPool(ctx, "pool1")
	require.Equal(t, "1199990", pool1.NativeAssetBalance.String())

	pool2, _ := app.ClpKeeper.GetPool(ctx, "pool2")
	require.Equal(t, "2400000", pool2.NativeAssetBalance.String())

}

//TestCase#4: Testing reward allocation for odd number of pools with diferent depth
func TestEndBlockWithOddPoolsDifferentDepth(t *testing.T) {
	fmt.Println("Testing TestEndBlockWithOddNumberOfPools")
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	var err error
	for i := 1; i <= 3; i++ {
		err = app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: "pool" + strconv.Itoa(i)},
			NativeAssetBalance:   sdk.NewUint(uint64(i * 1000000)),
			ExternalAssetBalance: sdk.NewUint(uint64(i * 1000000)),
			PoolUnits:            sdk.NewUint(uint64(i * 1000000)),
		})
		require.NoError(t, err)
	}
	allocation := sdk.NewUint(600000)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
	}
	app.ClpKeeper.SetParams(ctx, params)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.True(t, periodOneSupply.IsGTE(startingSupply))

	pool1, _ := app.ClpKeeper.GetPool(ctx, "pool1")
	require.Equal(t, "1100000", pool1.NativeAssetBalance.String())

	pool2, _ := app.ClpKeeper.GetPool(ctx, "pool2")
	require.Equal(t, "2199990", pool2.NativeAssetBalance.String())

	pool3, _ := app.ClpKeeper.GetPool(ctx, "pool3")
	require.Equal(t, "3300000", pool3.NativeAssetBalance.String())
}

//TestCase#5: Testing when reward allocation is 0
//Test is Failing with this error: --- FAIL: TestEndBlockWithAllocationZero (0.00s)
//panic: value from ParamSetPair is invalid: reward period allocation must be positive: &{824651015872} [recovered]
//panic: value from ParamSetPair is invalid: reward period allocation must be positive: &{824651015872}
func TestEndBlockWithAllocationZero(t *testing.T) {
	fmt.Println("Testing TestEndBlockWithAllocationZero")
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	var err error
	for i := 1; i <= 2; i++ {
		err = app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: "pool" + strconv.Itoa(i)},
			NativeAssetBalance:   sdk.NewUint(10000000),
			ExternalAssetBalance: sdk.NewUint(10000000),
			PoolUnits:            sdk.NewUint(10000000),
		})
		require.NoError(t, err)
	}
	allocation := sdk.NewUint(0)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
	}
	app.ClpKeeper.SetParams(ctx, params)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.Equal(t, periodOneSupply, startingSupply)

	for i := 1; i <= 2; i++ {
		pool, err := app.ClpKeeper.GetPool(ctx, "pool"+strconv.Itoa(i))
		require.NoError(t, err)
		require.Equal(t, "10000000", pool.NativeAssetBalance.String())
	}

}

//TestCase#6: Reward distribution outside valid reward period
func TestEndBlockRewardDistrForInvalidPeriod(t *testing.T) {
	fmt.Println("Testing TestEndBlockRewardDistrForInvalidPeriod")
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	var err error
	for i := 1; i <= 2; i++ {
		err = app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: "pool" + strconv.Itoa(i)},
			NativeAssetBalance:   sdk.NewUint(10000000),
			ExternalAssetBalance: sdk.NewUint(10000000),
			PoolUnits:            sdk.NewUint(10000000),
		})
		require.NoError(t, err)
	}
	allocation := sdk.NewUint(600000)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 11, EndBlock: 13, Allocation: &allocation},
	}
	app.ClpKeeper.SetParams(ctx, params)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 10; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.Equal(t, periodOneSupply, startingSupply)

	for i := 1; i <= 2; i++ {
		pool, err := app.ClpKeeper.GetPool(ctx, "pool"+strconv.Itoa(i))
		require.NoError(t, err)
		require.Equal(t, "10000000", pool.NativeAssetBalance.String())
	}
}

//TestCase#7: Test Unlockliquidity - When no unlocks requested
func TestUnlockedLiquidityNoUnlock(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	params.LiquidityRemovalLockPeriod = 10
	tt := unlockliquidity{
		name:     "No unlocks",
		height:   1,
		use:      sdk.NewUint(1000000),
		expected: types.ErrBalanceNotAvailable,
	}
	ctx = ctx.WithBlockHeight(tt.height)
	app.ClpKeeper.SetParams(ctx, params)
	lp := types.LiquidityProvider{
		Asset:                    &types.Asset{Symbol: "pool1"},
		LiquidityProviderAddress: "sif123",
		Unlocks:                  tt.unlocks,
	}
	app.ClpKeeper.SetLiquidityProvider(ctx, &lp)
	err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, sdk.NewUint(1000000))
	require.Equal(t, err, tt.expected)
}

//TestCase#8: Test Unlockliquidity - When locking period is not yet valid
func TestUnlockedLiquidityInvalidPeriod(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	params.LiquidityRemovalLockPeriod = 10
	tt := unlockliquidity{
		name:     "Unlock not ready",
		height:   5,
		use:      sdk.NewUint(1000000),
		expected: types.ErrBalanceNotAvailable,
		unlocks: []*types.LiquidityUnlock{
			{
				RequestHeight: 1,
				Units:         sdk.NewUint(1000000),
			},
		},
	}
	ctx = ctx.WithBlockHeight(tt.height)
	app.ClpKeeper.SetParams(ctx, params)
	lp := types.LiquidityProvider{
		Asset:                    &types.Asset{Symbol: "pool1"},
		LiquidityProviderAddress: "sif123",
		Unlocks:                  tt.unlocks,
	}
	app.ClpKeeper.SetLiquidityProvider(ctx, &lp)
	err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, sdk.NewUint(1000000))
	require.Equal(t, err, tt.expected)
}

//TestCase#9: Test Unlockliquidity - When locking period is valid but no liquidity or balance
func TestUnlockedLiquidityNoBalances(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	params.LiquidityRemovalLockPeriod = 10
	tt := unlockliquidity{
		name:     "Insufficient balance/liquidity",
		height:   11,
		use:      sdk.NewUint(1000000),
		expected: types.ErrBalanceNotAvailable,
		unlocks: []*types.LiquidityUnlock{
			{
				RequestHeight: 1,
				Units:         sdk.NewUint(1000000),
			},
		},
	}
	ctx = ctx.WithBlockHeight(tt.height)
	app.ClpKeeper.SetParams(ctx, params)
	lp := types.LiquidityProvider{
		Asset:                    &types.Asset{Symbol: "pool1"},
		LiquidityProviderAddress: "sif123",
		Unlocks:                  tt.unlocks,
	}
	app.ClpKeeper.SetLiquidityProvider(ctx, &lp)
	err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, sdk.NewUint(10000))
	require.Equal(t, err, tt.expected)
}

//TestCase#10: Test Unlockliquidity - When locking period is valid and sufficient liquidity
func TestUnlockedLiquidity(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	params.LiquidityRemovalLockPeriod = 10
	tt := unlockliquidity{
		name:     "",
		height:   11,
		use:      sdk.NewUint(2000000),
		expected: nil,
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
	}
	ctx = ctx.WithBlockHeight(tt.height)
	app.ClpKeeper.SetParams(ctx, params)
	lp := types.LiquidityProvider{
		Asset:                    &types.Asset{Symbol: "pool1"},
		LiquidityProviderAddress: "sif123",
		Unlocks:                  tt.unlocks,
	}
	app.ClpKeeper.SetLiquidityProvider(ctx, &lp)
	err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, sdk.NewUint(3000000))
	require.Equal(t, err, tt.expected)
}

//TestCase#11: Test Unlockliquidity - When locking period is valid and some of funds requested for unlock are available but not all
//or partial liquidity exists
func TestUnlockedLiquidityWithPartialBalances(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	params.LiquidityRemovalLockPeriod = 10
	tt := unlockliquidity{
		name:     "",
		height:   11,
		use:      sdk.NewUint(3000000),
		expected: nil,
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
	}
	ctx = ctx.WithBlockHeight(tt.height)
	app.ClpKeeper.SetParams(ctx, params)
	lp := types.LiquidityProvider{
		Asset:                    &types.Asset{Symbol: "pool1"},
		LiquidityProviderAddress: "sif123",
		Unlocks:                  tt.unlocks,
	}
	app.ClpKeeper.SetLiquidityProvider(ctx, &lp)
	err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, sdk.NewUint(2000000))
	//	liquidity,_ := app.keeper.GetLiquidityProvider(ctx,lp.Asset.Symbol,LiquidityProviderAddress);
	require.Equal(t, err, tt.expected)
	//	s.Require().Equal(err, ()
}

//TestCase#12: Reward allocation test when you have two reward period sets one after other
func TestEndBlockWithTwoRewardPeriods(t *testing.T) {
	fmt.Println("Testing TestEndBlockWithTwoRewardPeriods")
	app, ctx := test.CreateTestApp(false)
	params := app.ClpKeeper.GetParams(ctx)
	var err error
	for i := 1; i <= 2; i++ {
		err = app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: "pool" + strconv.Itoa(i)},
			NativeAssetBalance:   sdk.NewUint(uint64(i * 10000000)),
			ExternalAssetBalance: sdk.NewUint(uint64(i * 10000000)),
			PoolUnits:            sdk.NewUint(uint64(i * 10000000)),
		})
		require.NoError(t, err)
	}
	allocation := sdk.NewUint(600000)
	allocation1 := sdk.NewUint(1200000)
	params.RewardPeriods = []*types.RewardPeriod{
		{Id: "Test 1", StartBlock: 1, EndBlock: 10, Allocation: &allocation},
		{Id: "Test 2", StartBlock: 11, EndBlock: 15, Allocation: &allocation1},
	}
	app.ClpKeeper.SetParams(ctx, params)
	startingSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	for block := 1; block <= 15; block++ {
		app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
		app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
		app.Commit()
	}
	/*	for block := 11; block <= 15; block++ {
			app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
			app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
			app.Commit()
		}
	*/
	// check total supply change is as expected
	periodOneSupply := app.BankKeeper.GetSupply(ctx, "rowan")
	require.True(t, periodOneSupply.IsGTE(startingSupply))
	// check pool has expected increase
	/*pool1, _ := app.ClpKeeper.GetPool(ctx, "pool1")
	require.Equal(t, "10300000", pool1.NativeAssetBalance.String())

	pool2, _ := app.ClpKeeper.GetPool(ctx, "pool2")
	require.Equal(t, "10300000", pool2.NativeAssetBalance.String())*/

}
