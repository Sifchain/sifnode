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

func TestKeeper_CollectProviderDistribution(t *testing.T) {
	blockRate := sdk.MustNewDecFromStr("0.003141590000000000")
	poolDepthRowan := sdk.NewDec(200_000)
	totalProviderDistributioned := sdk.NewUint(628) // blockRate * poolDepthRowan
	poolUnitss := make([]uint64, 5)
	poolUnitss[0] = 10
	poolUnitss[1] = 0
	poolUnitss[2] = 3
	poolUnitss[3] = 5
	poolUnitss[4] = 12

	totalPoolUnits := uint64(0)
	for i := 0; i < len(poolUnitss); i++ {
		totalPoolUnits += poolUnitss[i]
	}

	lps := test.GenerateRandomLPWithUnits(poolUnitss)
	cbm := make(keeper.ProviderDistributionMap)

	keeper.CollectProviderDistribution(poolDepthRowan, blockRate, sdk.NewUint(totalPoolUnits), lps, cbm)

	firstProviderDistributionAmount := sdk.NewUint(209)
	require.Equal(t, firstProviderDistributionAmount, cbm[lps[0].LiquidityProviderAddress])

	secondProviderDistributionAmount := sdk.ZeroUint()
	require.Equal(t, secondProviderDistributionAmount, cbm[lps[1].LiquidityProviderAddress])

	thirdProviderDistributionAmount := sdk.NewUint(63)
	require.Equal(t, thirdProviderDistributionAmount, cbm[lps[2].LiquidityProviderAddress])

	fourthProviderDistributionAmount := sdk.NewUint(105)
	require.Equal(t, fourthProviderDistributionAmount, cbm[lps[3].LiquidityProviderAddress])

	fifthProviderDistributionAmount := sdk.NewUint(251)
	require.Equal(t, fifthProviderDistributionAmount, cbm[lps[4].LiquidityProviderAddress])

	sum := firstProviderDistributionAmount.Add(secondProviderDistributionAmount).Add(thirdProviderDistributionAmount).Add(fourthProviderDistributionAmount).Add(fifthProviderDistributionAmount)
	require.Equal(t, totalProviderDistributioned, sum)
}

func TestKeeper_CollectProviderDistributions(t *testing.T) {
	blockRate := sdk.MustNewDecFromStr("0.003141590000000000")
	nPools := 5
	nLPs := 3
	ctx, app := test.CreateTestAppClp(false)
	pools := test.GeneratePoolsSetLPs(app.ClpKeeper, ctx, nPools, nLPs)
	cbm := app.ClpKeeper.CollectProviderDistributions(ctx, pools, blockRate)

	// TODO: something better
	require.Equal(t, nLPs, len(cbm))
}
