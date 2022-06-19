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
	_ = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.NativeSymbol, sdk.NewInt(628))))
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
	cbm := make(keeper.ProviderDistributionMap)

	keeper.CollectProviderDistribution(poolDepthRowan, blockRate, sdk.NewUint(totalPoolUnits), lps, cbm)

	for i, providerDistribution := range providerDistributions {
		require.Equal(t, providerDistribution, cbm[lps[i].LiquidityProviderAddress])

		// We clear the EventManager before every call as Events accumulate throughout calls
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		addr, _ := sdk.AccAddressFromBech32(lps[i].LiquidityProviderAddress)
		err := app.ClpKeeper.TransferProviderDistribution(ctx, addr, providerDistribution)
		require.Nil(t, err)
		transferEvents := createTransferEvents(providerDistribution, addr)

		// NOTE: we use Subset here as bankKeeper.SendCoinsFromModuleToAccount does fire Events itself which we do not care for at this point
		require.Subset(t, ctx.EventManager().Events(), transferEvents)
	}
}

func createTransferEvents(amount sdk.Uint, receiver sdk.AccAddress) []sdk.Event {
	return []sdk.Event{sdk.NewEvent("lppd_distribution",
		sdk.NewAttribute("lppd_distribution_amount", sdk.NewCoin(types.NativeSymbol, sdk.Int(amount)).String()),
		sdk.NewAttribute("lppd_distribution_receiver", receiver.String()),
		sdk.NewAttribute("height", "0")),
	}
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
