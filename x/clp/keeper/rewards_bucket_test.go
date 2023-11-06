package keeper_test

import (
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
	keeper, ctx := keepertest.ClpKeeper(t)
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
	keeper, ctx := keepertest.ClpKeeper(t)
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
	keeper, ctx := keepertest.ClpKeeper(t)
	items := createNRewardsBucket(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllRewardsBucket(ctx)),
	)
}
