package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	keepertest "github.com/Sifchain/sifnode/testutil/keeper"
	"github.com/Sifchain/sifnode/testutil/nullify"
	"github.com/Sifchain/sifnode/x/clp/keeper"
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
