package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_CalcProviderDistributionAmount(t *testing.T) {
	rowanProviderDistributioned := sdk.NewDec(10)
	totalPoolUnits := sdk.NewUint(999)
	lpPoolUnits := sdk.NewUint(333)
	expectedAmount := sdk.NewUint(3)

	amount := keeper.CalcProviderDistributionAmount(rowanProviderDistributioned, totalPoolUnits, lpPoolUnits)

	require.Equal(t, expectedAmount, amount)
}

func TestKeeper_FindActivePeriod(t *testing.T) {
	firstPeriod := types.ProviderDistributionPeriod{DistributionPeriodStartBlock: 4, DistributionPeriodEndBlock: 10, DistributionPeriodBlockRate: sdk.NewDec(1)}
	secondPeriod := types.ProviderDistributionPeriod{DistributionPeriodStartBlock: 8, DistributionPeriodEndBlock: 12, DistributionPeriodBlockRate: sdk.NewDec(1)}
	thirdPeriod := types.ProviderDistributionPeriod{DistributionPeriodStartBlock: 20, DistributionPeriodEndBlock: 20, DistributionPeriodBlockRate: sdk.NewDec(1)}

	periods := make([]*types.ProviderDistributionPeriod, 3)
	periods[0] = &firstPeriod
	periods[1] = &secondPeriod
	periods[2] = &thirdPeriod

	currentHeight := int64(0)
	period := keeper.FindProviderDistributionPeriod(currentHeight, periods)
	require.Nil(t, period)

	currentHeight = 4
	period = keeper.FindProviderDistributionPeriod(currentHeight, periods)
	require.Equal(t, &firstPeriod, period)

	currentHeight = 10
	period = keeper.FindProviderDistributionPeriod(currentHeight, periods)
	require.Equal(t, &firstPeriod, period)

	currentHeight = 11
	period = keeper.FindProviderDistributionPeriod(currentHeight, periods)
	require.Equal(t, &secondPeriod, period)

	currentHeight = 20
	period = keeper.FindProviderDistributionPeriod(currentHeight, periods)
	require.Equal(t, &thirdPeriod, period)

	currentHeight = 30
	period = keeper.FindProviderDistributionPeriod(currentHeight, periods)
	require.Nil(t, period)
}

func TestKeeper_CollectProviderDistributionAndEvents(t *testing.T) {
	blockRate := sdk.MustNewDecFromStr("0.003141590000000000")
	poolDepthRowan := sdk.NewDec(200_000)
	totalProviderDistributioned := sdk.NewUint(628) // blockRate * poolDepthRowan
	// only used for events collection
	ctx, app := test.CreateTestAppClp(false)
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.NativeSymbol, sdk.NewInt(2*628)))) //x2 since there's 2 pools
	// clear MintCoins events
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	poolUnitss := []uint64{10, 0, 3, 5, 12}
	providerDistributions := []sdk.Uint{sdk.NewUint(209), sdk.ZeroUint(), sdk.NewUint(63), sdk.NewUint(105), sdk.NewUint(251)}
	totalPoolUnits := uint64(0)
	providerSum := sdk.ZeroUint()

	for i := 0; i < len(poolUnitss); i++ {
		totalPoolUnits += poolUnitss[i]
		providerSum = providerSum.Add(providerDistributions[i])
	}
	require.Equal(t, totalProviderDistributioned, providerSum)

	lps := test.GenerateRandomLPWithUnits(poolUnitss)
	assetStr := "cusdc"
	asset := types.NewAsset(assetStr)
	pool := types.NewPool(&asset, totalProviderDistributioned, sdk.ZeroUint(), sdk.ZeroUint())

	lpsFiltered := keeper.FilterValidLiquidityProviders(ctx, lps)
	lpRowanMap := make(keeper.LpRowanMap, 0)
	lpPoolMap := make(keeper.LpPoolMap, 0)
	poolRowanMap := make(keeper.PoolRowanMap, 1)
	rowanToDistribute := keeper.CollectProviderDistribution(ctx, &pool, poolDepthRowan, blockRate, sdk.NewUint(totalPoolUnits), lpsFiltered, lpRowanMap, lpPoolMap)
	require.Equal(t, totalProviderDistributioned, rowanToDistribute)
	poolRowanMap[&pool] = rowanToDistribute

	for i, providerDistribution := range providerDistributions {
		addr := lps[i].LiquidityProviderAddress
		//addr, _ := sdk.AccAddressFromBech32(lps[i].LiquidityProviderAddress)
		//tuple := findTupleByAddress(addr, tuples)
		require.Equal(t, providerDistribution, lpRowanMap[addr])

		// We clear the EventManager before every call as Events accumulate throughout calls
		//ctx = ctx.WithEventManager(sdk.NewEventManager())

		//transferEvents := createTransferEvents(providerDistribution, addr)
		// NOTE: we use Subset here as bankKeeper.SendCoinsFromModuleToAccount does fire Events itself which we do not care for at this point
		//require.Subset(t, ctx.EventManager().Events(), transferEvents)
	}

	// repeat for a second pool
	assetStr2 := "ceth"
	asset2 := types.NewAsset(assetStr2)
	pool2 := types.NewPool(&asset2, totalProviderDistributioned, sdk.ZeroUint(), sdk.ZeroUint())
	rowanToDistribute2 := keeper.CollectProviderDistribution(ctx, &pool2, poolDepthRowan, blockRate, sdk.NewUint(totalPoolUnits), lpsFiltered, lpRowanMap, lpPoolMap)
	poolRowanMap[&pool2] = rowanToDistribute2

	app.ClpKeeper.TransferProviderDistribution(ctx, poolRowanMap, lpRowanMap, lpPoolMap)

	// pool empty after all LPs got paid
	poolStored, _ := app.ClpKeeper.GetPool(ctx, assetStr)
	require.Equal(t, sdk.ZeroUint().String(), poolStored.NativeAssetBalance.String())

	require.Subset(t, ctx.EventManager().Events(), createDistributeEvent(lps[len(lps)-1].LiquidityProviderAddress))
}

//nolint
func createDistributeEvent(address string) []sdk.Event {
	return []sdk.Event{sdk.NewEvent("lppd/distribution",
		sdk.NewAttribute("recipient", address),
		sdk.NewAttribute("total_amount", "502"),
		sdk.NewAttribute("amounts", "[{\"pool\":\"cusdc\",\"amount\":\"251\"},{\"pool\":\"ceth\",\"amount\":\"251\"}]")),
	}
}

func TestKeeper_CollectProviderDistributions(t *testing.T) {
	blockRate := sdk.MustNewDecFromStr("0.003141590000000000")
	nPools := 100
	nLPs := 800
	ctx, app := test.CreateTestAppClp(false)
	pools := test.GeneratePoolsSetLPs(app.ClpKeeper, ctx, nPools, nLPs)
	poolRowanMap, lpRowanMap, lpPoolMap := app.ClpKeeper.CollectProviderDistributions(ctx, pools, blockRate)

	// TODO: something better
	require.Equal(t, nPools, len(poolRowanMap))
	require.Equal(t, nLPs, len(lpRowanMap))
	require.Equal(t, len(lpPoolMap), len(lpRowanMap))
}

func TestKeeper_IsDistributionBlock(t *testing.T) {
	startHeight := uint64(12)
	blockHeight := int64(12)
	mod := uint64(4)

	require.True(t, keeper.IsDistributionBlockPure(blockHeight, startHeight, mod))

	blockHeight = 13
	require.False(t, keeper.IsDistributionBlockPure(blockHeight, startHeight, mod))

	blockHeight = 14
	require.False(t, keeper.IsDistributionBlockPure(blockHeight, startHeight, mod))

	blockHeight = 15
	require.False(t, keeper.IsDistributionBlockPure(blockHeight, startHeight, mod))

	blockHeight = 16
	require.True(t, keeper.IsDistributionBlockPure(blockHeight, startHeight, mod))
}
