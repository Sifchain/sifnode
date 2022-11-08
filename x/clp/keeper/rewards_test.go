package keeper_test

import (
	"fmt"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tenderminttypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestEndBlock(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	// Setup reward period
	params := app.ClpKeeper.GetRewardsParams(ctx)

	allocation := sdk.NewUintFromString("200000000000000000000000000")
	oneDec := sdk.OneDec()
	params.RewardPeriods = []*types.RewardPeriod{
		{RewardPeriodId: "Test 1", RewardPeriodStartBlock: 1, RewardPeriodEndBlock: 10, RewardPeriodAllocation: &allocation, RewardPeriodDefaultMultiplier: &oneDec},
	}
	app.ClpKeeper.SetRewardParams(ctx, params)
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin("rowan", sdk.NewInt(3000)),
		sdk.NewCoin("atom", sdk.NewInt(1000)),
		sdk.NewCoin("cusdc", sdk.NewInt(1000)),
		sdk.NewCoin("ceth", sdk.NewInt(1000)),
	))
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &types.Pool{
		ExternalAsset:                 &types.Asset{Symbol: "atom"},
		NativeAssetBalance:            sdk.NewUint(1000),
		ExternalAssetBalance:          sdk.NewUint(1000),
		PoolUnits:                     sdk.NewUint(1000),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &types.Pool{
		ExternalAsset:                 &types.Asset{Symbol: "cusdc"},
		NativeAssetBalance:            sdk.NewUint(1000),
		ExternalAssetBalance:          sdk.NewUint(1000),
		PoolUnits:                     sdk.NewUint(1000),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &types.Pool{
		ExternalAsset:                 &types.Asset{Symbol: "ceth"},
		NativeAssetBalance:            sdk.NewUint(1000),
		ExternalAssetBalance:          sdk.NewUint(1000),
		PoolUnits:                     sdk.NewUint(1000),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
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
	// TODO : Modify reward policy so that the numbers asserted match with expected
	//pool, err := app.ClpKeeper.GetPool(ctx, "atom")
	//require.NoError(t, err)
	//require.Equal(t, "66666666666666666600001000", pool.NativeAssetBalance.String())
	//expected := sdk.NewUintFromString("66666666666666666666667666")
	//accuracy := sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()).Quo(sdk.NewDecFromBigInt(expected.BigInt()))
	//require.True(t, accuracy.GT(sdk.MustNewDecFromStr("0.99")))
	//// TODO continue through another portion of the period and ensure supply is increased.
	//// continue through a non reward period
	//for block := 11; block <= 20; block++ {
	//	app.BeginBlock(abci.RequestBeginBlock{Header: tenderminttypes.Header{Height: int64(block)}})
	//	app.EndBlock(abci.RequestEndBlock{Height: int64(block)})
	//	app.Commit()
	//}
	//// check total supply is unchanged
	//supplyCheck := app.BankKeeper.GetSupply(ctx, "rowan")
	////log.Printf("starting supply: %s final supply: %s after period one: %s", startingSupply.String(), supplyCheck.String(), periodOneSupply.String())
	//require.True(t, supplyCheck.Equal(periodOneSupply))
}

func TestUseUnlockedLiquidity(t *testing.T) {
	tt := []struct {
		name     string
		height   int64
		use      sdk.Uint
		any      bool
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
			name:     "Unlock in any state",
			height:   5,
			use:      sdk.NewUint(1000),
			any:      true,
			expected: nil,
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
			err := app.ClpKeeper.UseUnlockedLiquidity(ctx, lp, sdk.NewUint(1000), tc.any)
			require.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestKeeper_RewardsDistribution(t *testing.T) {
	startBalance := sdk.NewCoin(types.NativeSymbol, sdk.NewInt(42))
	ctx, app := test.CreateTestAppClp(false)
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(startBalance))

	allocation := sdk.NewUintFromString("200000000000000000000000000")
	totalCoinsDistribution := sdk.NewCoin(types.NativeSymbol, sdk.NewIntFromBigInt(allocation.BigInt()))
	oneDec := sdk.OneDec()
	period := types.RewardPeriod{
		RewardPeriodId: "Test 1", RewardPeriodStartBlock: 0, RewardPeriodEndBlock: 0, RewardPeriodAllocation: &allocation,
		RewardPeriodDefaultMultiplier: &oneDec, RewardPeriodDistribute: true, RewardPeriodMod: 0,
		RewardPeriodPoolMultipliers: nil}

	pools := test.GeneratePoolsSetLPs(app.ClpKeeper, ctx, 1, 1)
	pool := pools[0]
	lps, _ := app.ClpKeeper.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
	lp := lps[0]
	lpAddr, _ := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)

	lpCoinsBefore := app.BankKeeper.GetBalance(ctx, lpAddr, types.NativeSymbol)
	// Disitribute coins to the LP
	blockDistribution := keeper.CalcBlockDistribution(&period)
	err := app.ClpKeeper.DistributeDepthRewards(ctx, blockDistribution, &period, pools)
	require.Nil(t, err)
	lpCoinsAfter1 := app.BankKeeper.GetBalance(ctx, lpAddr, types.NativeSymbol)
	require.True(t, lpCoinsBefore.IsLT(lpCoinsAfter1))
	require.Equal(t, lpCoinsAfter1, totalCoinsDistribution)
	require.Subset(t, ctx.EventManager().Events(), createRewardsDistributeEvent(totalCoinsDistribution, pool.ExternalAsset))

	distributed1 := pool.RewardPeriodNativeDistributed
	moduleBalance1 := app.ClpKeeper.GetModuleRowan(ctx)
	// We transferred all minted coins
	require.Equal(t, startBalance.String(), moduleBalance1.String())

	// This time, we do not distribute coins to the LP
	period.RewardPeriodDistribute = false
	// reset to easier keep track of it
	pool.RewardPeriodNativeDistributed = sdk.ZeroUint()
	err = app.ClpKeeper.DistributeDepthRewards(ctx, blockDistribution, &period, pools)
	require.Nil(t, err)
	distributed2 := pool.RewardPeriodNativeDistributed
	require.Equal(t, distributed1, distributed2)
	require.Subset(t, ctx.EventManager().Events(), createRewardsAccumEvent(totalCoinsDistribution, pool.ExternalAsset))

	lpCoinsAfter2 := app.BankKeeper.GetBalance(ctx, lpAddr, types.NativeSymbol)
	require.Equal(t, lpCoinsAfter1, lpCoinsAfter2)

	moduleBalance2 := app.ClpKeeper.GetModuleRowan(ctx)
	diffBalance := sdk.NewCoin(types.NativeSymbol, sdk.NewIntFromBigInt(distributed2.BigInt()))
	// we did not distribute the newly minted coins
	require.Equal(t, startBalance.Add(diffBalance).String(), moduleBalance2.String())
}

// nolint
func createRewardsDistributeEvent(totalCoinsDistribution sdk.Coin, asset *types.Asset) []sdk.Event {
	amountsStr := fmt.Sprintf("[{\"pool\":\"%s\",\"amount\":\"200000000000000000000000000\"}]", asset.Symbol)
	return []sdk.Event{sdk.NewEvent("rewards/distribution",
		sdk.NewAttribute("total_amount", totalCoinsDistribution.Amount.String()),
		sdk.NewAttribute("amounts", amountsStr)),
	}
}

func createRewardsAccumEvent(totalCoinsDistribution sdk.Coin, asset *types.Asset) []sdk.Event {
	amountsStr := fmt.Sprintf("[{\"pool\":\"%s\",\"amount\":\"200000000000000000000000000\"}]", asset.Symbol)
	return []sdk.Event{sdk.NewEvent("rewards/accumulation",
		sdk.NewAttribute("total_amount", totalCoinsDistribution.Amount.String()),
		sdk.NewAttribute("amounts", amountsStr)),
	}
}

func TestKeeper_RewardsDistributionFailure(t *testing.T) {
	sifapp.SetConfig(false) // needed for GenerateAddress to generate a proper address
	// We cheat here to get the first LP's address in order to add it to the blacklist
	lpAddress := test.GenerateAddress2(fmt.Sprintf("%d%d%d%d", 0, 0, 0, 0))
	blacklist := []sdk.AccAddress{lpAddress}
	startBalance := sdk.NewCoin(types.NativeSymbol, sdk.NewInt(42))

	ctx, app := test.CreateTestAppClpWithBlacklist(false, blacklist)
	require.True(t, app.BankKeeper.BlockedAddr(lpAddress))
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(startBalance))

	allocation := sdk.NewUintFromString("200000000000000000000000000")
	oneDec := sdk.OneDec()
	period := types.RewardPeriod{
		RewardPeriodId: "Test 1", RewardPeriodStartBlock: 0, RewardPeriodEndBlock: 0, RewardPeriodAllocation: &allocation,
		RewardPeriodDefaultMultiplier: &oneDec, RewardPeriodDistribute: true, RewardPeriodMod: 0,
		RewardPeriodPoolMultipliers: nil}

	pools := test.GeneratePoolsSetLPs(app.ClpKeeper, ctx, 1, 1)
	pool := pools[0]
	lps, _ := app.ClpKeeper.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
	lp := lps[0]
	lpAddr, _ := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	require.Equal(t, lpAddress, lpAddr)

	lpCoinsBefore := app.BankKeeper.GetBalance(ctx, lpAddress, types.NativeSymbol)
	// Distribute coins to the LP
	blockDistribution := keeper.CalcBlockDistribution(&period)
	err := app.ClpKeeper.DistributeDepthRewards(ctx, blockDistribution, &period, pools)
	require.Nil(t, err)

	// Nope, distribution failed
	lpCoinsAfter := app.BankKeeper.GetBalance(ctx, lpAddress, types.NativeSymbol)
	require.Equal(t, lpCoinsBefore, lpCoinsAfter)
	require.Equal(t, pool.RewardPeriodNativeDistributed.String(), sdk.ZeroInt().String())

	// Check tokens got burnt
	moduleBalance := app.ClpKeeper.GetModuleRowan(ctx)
	require.Equal(t, startBalance, moduleBalance)

	failedEvent := createFailedEvent(lpAddress)
	require.Subset(t, ctx.EventManager().Events(), failedEvent)
}

func createFailedEvent(receiver sdk.AccAddress) []sdk.Event {
	return []sdk.Event{sdk.NewEvent("rewards/distribution_error",
		sdk.NewAttribute("liquidity_provider", receiver.String()),
		sdk.NewAttribute("error", fmt.Sprint(receiver.String(), " is not allowed to receive funds: unauthorized")),
		sdk.NewAttribute("height", "0")),
	}
}
