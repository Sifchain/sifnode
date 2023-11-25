package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	keepertest "github.com/Sifchain/sifnode/testutil/keeper"
	"github.com/Sifchain/sifnode/testutil/nullify"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNRewardsBucket(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.RewardsBucket {
	items := make([]types.RewardsBucket, n)
	for i := range items {
		items[i].Denom = strconv.Itoa(i)
		items[i].Amount = sdk.NewInt(int64(i))

		keeper.SetRewardsBucket(ctx, items[i])
	}
	return items
}

func TestRewardsBucketGet(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	items := createNRewardsBucket(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetRewardsBucket(ctx,
			item.Denom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item), //nolint
			nullify.Fill(&rst),
		)
	}
}
func TestRewardsBucketRemove(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	items := createNRewardsBucket(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveRewardsBucket(ctx,
			item.Denom,
		)
		_, found := keeper.GetRewardsBucket(ctx,
			item.Denom,
		)
		require.False(t, found)
	}
}

func TestRewardsBucketGetAll(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	items := createNRewardsBucket(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllRewardsBucket(ctx)),
	)
}

func TestAddToRewardsBucket(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	// Create a RewardsBucket with a denom and amount
	denom := "atom"
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  denom,
		Amount: initialAmount,
	}
	keeper.SetRewardsBucket(ctx, rewardsBucket)

	// Add amount to the RewardsBucket
	addAmount := sdk.NewInt(50)
	err := keeper.AddToRewardsBucket(ctx, denom, addAmount)
	require.NoError(t, err)

	// Check if the amount has been added
	storedRewardsBucket, found := keeper.GetRewardsBucket(ctx, denom)
	require.True(t, found)
	require.Equal(t, initialAmount.Add(addAmount), storedRewardsBucket.Amount)
}

func TestSubtractFromRewardsBucket(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	// Create a RewardsBucket with a denom and amount
	denom := "atom"
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  denom,
		Amount: initialAmount,
	}
	keeper.SetRewardsBucket(ctx, rewardsBucket)

	// Subtract amount from the RewardsBucket
	subtractAmount := sdk.NewInt(50)
	err := keeper.SubtractFromRewardsBucket(ctx, denom, subtractAmount)
	require.NoError(t, err)

	// Check if the amount has been subtracted
	storedRewardsBucket, found := keeper.GetRewardsBucket(ctx, denom)
	require.True(t, found)
	require.Equal(t, initialAmount.Sub(subtractAmount), storedRewardsBucket.Amount)
}

func TestAddMultipleCoinsToRewardsBuckets(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)
	// Create multiple RewardsBuckets
	createNRewardsBucket(keeper, ctx, 5)

	// Define coins to add
	coinsToAdd := sdk.NewCoins(
		sdk.NewCoin("aaa", sdk.NewInt(10)),
		sdk.NewCoin("bbb", sdk.NewInt(20)),
		sdk.NewCoin("ccc", sdk.NewInt(30)),
	)

	// Add multiple coins to the respective RewardsBuckets
	addedCoins, err := keeper.AddMultipleCoinsToRewardsBuckets(ctx, coinsToAdd)
	require.NoError(t, err)

	// Check if the amounts have been added correctly
	for _, coin := range addedCoins {
		storedRewardsBucket, found := keeper.GetRewardsBucket(ctx, coin.Denom)
		require.True(t, found)
		// The expected amount is the initial amount (which is equal to the index) plus the added amount
		expectedAmount := coinsToAdd.AmountOf(coin.Denom)
		require.Equal(t, expectedAmount, storedRewardsBucket.Amount)
	}
}

func TestAddToRewardsBucket_Errors(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)

	// Test for empty denom error
	err := keeper.AddToRewardsBucket(ctx, "", sdk.NewInt(10))
	require.ErrorIs(t, err, types.ErrDenomCantBeEmpty)

	// Test for negative amount error
	err = keeper.AddToRewardsBucket(ctx, "atom", sdk.NewInt(-1))
	require.ErrorIs(t, err, types.ErrAmountCantBeNegative)
}

func TestSubtractFromRewardsBucket_Errors(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)

	// Test for empty denom error
	err := keeper.SubtractFromRewardsBucket(ctx, "", sdk.NewInt(10))
	require.ErrorIs(t, err, types.ErrDenomCantBeEmpty)

	// Test for negative amount error
	err = keeper.SubtractFromRewardsBucket(ctx, "atom", sdk.NewInt(-1))
	require.ErrorIs(t, err, types.ErrAmountCantBeNegative)

	// Test for rewards bucket not found error
	err = keeper.SubtractFromRewardsBucket(ctx, "atom", sdk.NewInt(1))
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Errorf(types.ErrRewardsBucketNotFound.Error(), "atom").Error())

	// Test for not enough balance error
	// First, create a RewardsBucket with a small amount
	keeper.SetRewardsBucket(ctx, types.RewardsBucket{
		Denom:  "atom",
		Amount: sdk.NewInt(1),
	})
	// Try to subtract more than the available amount
	err = keeper.SubtractFromRewardsBucket(ctx, "atom", sdk.NewInt(10))
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Errorf(types.ErrNotEnoughBalanceInRewardsBucket.Error(), "atom").Error())
}

// ShouldDistributeRewards returns true if the epoch identifier is not in the list of already distributed epochs
func TestShouldDistributeRewards(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)

	params := keeper.GetRewardsParams(ctx)

	// Check if the rewards epoch should be distributed
	require.True(t, keeper.ShouldDistributeRewards(ctx, params.RewardsEpochIdentifier))

	// Check if the rewards epoch should be distributed with a wrong epoch identifier
	require.False(t, keeper.ShouldDistributeRewards(ctx, "wrong_epoch_identifier"))
}

// ShouldDistributeRewardsToLPWallet returns true if the rewards distribute to LP wallet parameter is true
func TestShouldDistributeRewardsToLPWallet(t *testing.T) {
	keeper, ctx, _ := keepertest.ClpKeeper(t)

	// Check if the rewards should be distributed to LP wallet is set to default false value
	require.False(t, keeper.ShouldDistributeRewardsToLPWallet(ctx))

	// set distribute rewards to lp addresses
	rewardsParams := types.GetDefaultRewardParams()
	rewardsParams.RewardsDistribute = true
	keeper.SetRewardParams(ctx, rewardsParams)

	// Check if the rewards should be distributed to LP wallet is set to true
	require.True(t, keeper.ShouldDistributeRewardsToLPWallet(ctx))
}

// DistributeLiquidityProviderRewards distributes rewards to a liquidity provider
func TestDistributeLiquidityProviderRewards(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper

	asset := types.NewAsset("cusdc")

	// Create a liquidity provider
	lpAddress, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
	require.NoError(t, err)
	lp := types.NewLiquidityProvider(&asset, sdk.NewUint(1), lpAddress, ctx.BlockHeight())
	clpKeeper.SetLiquidityProvider(ctx, &lp)

	// Create a RewardsBucket with a denom and amount
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  asset.Symbol,
		Amount: initialAmount,
	}
	clpKeeper.SetRewardsBucket(ctx, rewardsBucket)

	// mint coins to the module account
	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(asset.Symbol, initialAmount),
	))
	require.NoError(t, err)

	// check module balance
	moduleBalance := app.BankKeeper.GetBalance(ctx, types.GetCLPModuleAddress(), asset.Symbol)
	require.Equal(t, initialAmount, moduleBalance.Amount)

	// check account balance
	preBalance := app.BankKeeper.GetBalance(ctx, lpAddress, asset.Symbol)
	require.Equal(t, sdk.ZeroInt(), preBalance.Amount)

	// Distribute rewards to the liquidity provider
	err = clpKeeper.DistributeLiquidityProviderRewards(ctx, &lp, asset.Symbol, initialAmount)
	require.NoError(t, err)

	// check account balance
	postBalance := app.BankKeeper.GetBalance(ctx, lpAddress, asset.Symbol)
	require.Equal(t, preBalance.Amount.Add(initialAmount), postBalance.Amount)

	// Check distributed rewards got subtracted from the rewards bucket
	storedRewardsBucket, found := clpKeeper.GetRewardsBucket(ctx, asset.Symbol)
	require.True(t, found)
	require.Equal(t, sdk.ZeroInt(), storedRewardsBucket.Amount)
}

// CalculateRewardShareForLiquidityProviders calculates the reward share for each liquidity provider
func TestCalculateRewardShareForLiquidityProviders(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper

	// Create a liquidity provider
	lps := test.GenerateRandomLP(10)
	for _, lp := range lps {
		clpKeeper.SetLiquidityProvider(ctx, lp)
	}

	// Calculate reward share for liquidity providers
	rewardShares := clpKeeper.CalculateRewardShareForLiquidityProviders(ctx, lps)

	// Check if the rewards amount is correct
	require.Equal(t, []sdk.Dec{
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
		sdk.MustNewDecFromStr("0.1"),
	}, rewardShares)
}

// CalculateRewardAmountForLiquidityProviders calculates the reward amount for each liquidity provider
func TestCalculateRewardAmountForLiquidityProviders(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper

	// Create a RewardsBucket with a denom and amount
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  "atom",
		Amount: initialAmount,
	}
	clpKeeper.SetRewardsBucket(ctx, rewardsBucket)

	// Create a liquidity provider
	lps := test.GenerateRandomLP(10)
	for _, lp := range lps {
		clpKeeper.SetLiquidityProvider(ctx, lp)
	}

	// Calculate reward share for liquidity providers
	rewardShares := clpKeeper.CalculateRewardShareForLiquidityProviders(ctx, lps)

	// Calculate reward amount for liquidity providers
	rewardAmounts := clpKeeper.CalculateRewardAmountForLiquidityProviders(ctx, rewardShares, rewardsBucket.Amount)

	// Check if the rewards amount is correct
	require.Equal(t, []sdk.Int{
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
		sdk.NewInt(10),
	}, rewardAmounts)
}
