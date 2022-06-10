package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_CalcCashbackAmount(t *testing.T) {
	rowanCashbacked := sdk.NewDec(10)
	totalPoolUnits := sdk.NewUint(999)
	lpPoolUnits := sdk.NewUint(333)
	expectedAmount := sdk.NewUint(3)

	amount := keeper.CalcCashbackAmount(rowanCashbacked, totalPoolUnits, lpPoolUnits)

	require.Equal(t, expectedAmount, amount)
}

func TestKeeper_FindActivePeriod(t *testing.T) {
	firstPeriod := types.CashbackPeriod{CashbackPeriodStartBlock: 4, CashbackPeriodEndBlock: 10, CashbackPeriodBlockRate: sdk.NewDec(1)}
	secondPeriod := types.CashbackPeriod{CashbackPeriodStartBlock: 8, CashbackPeriodEndBlock: 12, CashbackPeriodBlockRate: sdk.NewDec(1)}
	thirdPeriod := types.CashbackPeriod{CashbackPeriodStartBlock: 20, CashbackPeriodEndBlock: 20, CashbackPeriodBlockRate: sdk.NewDec(1)}

	currentHeight := 0
	period := keeper.FindActiveCashbackPeriod(int64(currentHeight), []*types.CashbackPeriod{&firstPeriod, &secondPeriod, &thirdPeriod})
	require.Nil(t, period)

	currentHeight = 4
	period = keeper.FindActiveCashbackPeriod(int64(currentHeight), []*types.CashbackPeriod{&firstPeriod, &secondPeriod, &thirdPeriod})
	require.Equal(t, &firstPeriod, period)

	currentHeight = 10
	period = keeper.FindActiveCashbackPeriod(int64(currentHeight), []*types.CashbackPeriod{&firstPeriod, &secondPeriod, &thirdPeriod})
	require.Equal(t, &firstPeriod, period)

	currentHeight = 11
	period = keeper.FindActiveCashbackPeriod(int64(currentHeight), []*types.CashbackPeriod{&firstPeriod, &secondPeriod, &thirdPeriod})
	require.Equal(t, &secondPeriod, period)

	currentHeight = 20
	period = keeper.FindActiveCashbackPeriod(int64(currentHeight), []*types.CashbackPeriod{&firstPeriod, &secondPeriod, &thirdPeriod})
	require.Equal(t, &thirdPeriod, period)

	currentHeight = 30
	period = keeper.FindActiveCashbackPeriod(int64(currentHeight), []*types.CashbackPeriod{&firstPeriod, &secondPeriod, &thirdPeriod})
	require.Nil(t, period)
}

func TestKeeper_CollectCashback(t *testing.T) {
	blockRate := sdk.MustNewDecFromStr("0.003141590000000000")
	poolDepthRowan := sdk.NewDec(200_000)
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
	cbm := make(keeper.CashbackMap)

	keeper.CollectCashback(poolDepthRowan, blockRate, sdk.NewUint(totalPoolUnits), lps, cbm)

	firstCashbackAmount := sdk.NewUint(209)
	require.Equal(t, firstCashbackAmount, cbm[lps[0].LiquidityProviderAddress])

	secondCashbackAmount := sdk.ZeroUint()
	require.Equal(t, secondCashbackAmount, cbm[lps[1].LiquidityProviderAddress])

	thirdCashbackAmount := sdk.NewUint(63)
	require.Equal(t, thirdCashbackAmount, cbm[lps[2].LiquidityProviderAddress])

	fourthCashbackAmount := sdk.NewUint(105)
	require.Equal(t, fourthCashbackAmount, cbm[lps[3].LiquidityProviderAddress])

	fifthCashbackAmount := sdk.NewUint(251)
	require.Equal(t, fifthCashbackAmount, cbm[lps[4].LiquidityProviderAddress])
}

// multiple pools, non-disjoint lps
func TestKeeper_CollectCashbacks(t *testing.T) {
	blockRate := sdk.MustNewDecFromStr("0.003141590000000000")
	nPools := 5
	nLPs := 3
	ctx, app := test.CreateTestAppClp(false)
	pools := test.GeneratePoolsSetLPs(app.ClpKeeper, ctx, nPools, nLPs)
	cbm := app.ClpKeeper.CollectCashbacks(ctx, pools, blockRate)

	require.Equal(t, nLPs, len(cbm))
}
