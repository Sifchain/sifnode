package keeper_test

import (
	"testing"
	"time"

	keepertest "github.com/Sifchain/sifnode/testutil/keeper"
	"github.com/Sifchain/sifnode/testutil/nullify"
	"github.com/Sifchain/sifnode/x/epochs/keeper"
	"github.com/Sifchain/sifnode/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createEpochInfos(keeper *keeper.Keeper, ctx sdk.Context) []types.EpochInfo {

	items := make([]types.EpochInfo, 3)
	items[0] = types.EpochInfo{
		Identifier:            "daily",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}
	items[1] = types.EpochInfo{
		Identifier:            "weekly",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24 * 7,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}
	items[2] = types.EpochInfo{
		Identifier:            "monthly",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24 * 30,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}
	for i := range items {
		keeper.SetEpochInfo(ctx, items[i])
	}
	return items
}

func TestEpochInfoGet(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	items := createEpochInfos(keeper, ctx)
	for _, item := range items {
		rst, found := keeper.GetEpochInfo(ctx,
			item.Identifier,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item), //nolint:gosec
			nullify.Fill(&rst),
		)
	}
}
func TestEpochInfoRemove(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	items := createEpochInfos(keeper, ctx)
	for _, item := range items {
		keeper.DeleteEpochInfo(ctx,
			item.Identifier,
		)
		_, found := keeper.GetEpochInfo(ctx,
			item.Identifier,
		)
		require.False(t, found)
	}
}

func TestEntryGetAll(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	items := createEpochInfos(keeper, ctx)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.AllEpochInfos(ctx)),
	)
}

func TestGetSetDeleteEpochInfo(t *testing.T) {
	// Create a new instance of the EpochsKeeper and context
	k, ctx := keepertest.EpochsKeeper(t)

	// Create a new EpochInfo struct
	epoch := types.EpochInfo{
		Identifier:              "epoch1",
		StartTime:               time.Now().UTC(),
		Duration:                time.Hour * 24 * 7, // 1 week
		CurrentEpoch:            1,
		CurrentEpochStartTime:   time.Now().UTC(),
		EpochCountingStarted:    true,
		CurrentEpochStartHeight: 1,
	}

	// Set the epoch info
	k.SetEpochInfo(ctx, epoch)

	// Retrieve the epoch info
	storedEpoch, found := k.GetEpochInfo(ctx, epoch.Identifier)

	// Ensure the epoch info was found
	require.True(t, found)

	// Ensure the epoch info is correct
	require.Equal(t, epoch, storedEpoch)

	// Delete the epoch info
	k.DeleteEpochInfo(ctx, epoch.Identifier)

	// Ensure the epoch info was deleted
	_, found = k.GetEpochInfo(ctx, epoch.Identifier)
	require.False(t, found)
}

func TestIterateEpochInfo(t *testing.T) {
	// Create a new instance of the EpochsKeeper and context
	k, ctx := keepertest.EpochsKeeper(t)

	// Create some sample epoch infos
	epoch1 := types.EpochInfo{
		Identifier:              "epoch1",
		StartTime:               time.Now().UTC(),
		Duration:                time.Hour * 24 * 7, // 1 week
		CurrentEpoch:            1,
		CurrentEpochStartTime:   time.Now().UTC(),
		EpochCountingStarted:    true,
		CurrentEpochStartHeight: 1,
	}

	epoch2 := types.EpochInfo{
		Identifier:              "epoch2",
		StartTime:               time.Now().UTC(),
		Duration:                time.Hour * 24 * 7, // 1 week
		CurrentEpoch:            2,
		CurrentEpochStartTime:   time.Now().UTC(),
		EpochCountingStarted:    true,
		CurrentEpochStartHeight: 2,
	}

	epoch3 := types.EpochInfo{
		Identifier:              "epoch3",
		StartTime:               time.Now().UTC(),
		Duration:                time.Hour * 24 * 7, // 1 week
		CurrentEpoch:            3,
		CurrentEpochStartTime:   time.Now().UTC(),
		EpochCountingStarted:    true,
		CurrentEpochStartHeight: 3,
	}

	// Set the epoch infos
	k.SetEpochInfo(ctx, epoch1)
	k.SetEpochInfo(ctx, epoch2)
	k.SetEpochInfo(ctx, epoch3)

	// Iterate over the epoch infos and ensure they are correct
	expectedEpochs := []types.EpochInfo{epoch1, epoch2, epoch3}
	var i int64 = 0
	k.IterateEpochInfo(ctx, func(index int64, epoch types.EpochInfo) (stop bool) {
		require.Equal(t, expectedEpochs[index], epoch)
		require.Equal(t, i, index)
		i++
		return false
	})
}
