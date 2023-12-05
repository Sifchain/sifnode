package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/epochs"
	"github.com/Sifchain/sifnode/x/epochs/types"
)

func TestEpochInfoChangesBeginBlockerAndInitGenesis(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)

	var (
		epochInfo types.EpochInfo
		found     bool
	)

	now := time.Now()

	for _, tc := range []struct {
		desc                       string
		expCurrentEpochStartTime   time.Time
		expCurrentEpochStartHeight int64
		expCurrentEpoch            int64
		expInitialEpochStartTime   time.Time
		fn                         func()
	}{
		{
			// Only advance 2 seconds, do not increment epoch
			expCurrentEpochStartHeight: 2,
			expCurrentEpochStartTime:   now,
			expCurrentEpoch:            1,
			expInitialEpochStartTime:   now,
			fn: func() {
				ctx = ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				app.EpochsKeeper.BeginBlocker(ctx)
				epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
				require.True(t, found)
			},
		},
		{
			expCurrentEpochStartHeight: 2,
			expCurrentEpochStartTime:   now,
			expCurrentEpoch:            1,
			expInitialEpochStartTime:   now,
			fn: func() {
				ctx = ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				app.EpochsKeeper.BeginBlocker(ctx)
				epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
				require.True(t, found)
			},
		},
		{
			expCurrentEpochStartHeight: 2,
			expCurrentEpochStartTime:   now,
			expCurrentEpoch:            1,
			expInitialEpochStartTime:   now,
			fn: func() {
				ctx = ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				app.EpochsKeeper.BeginBlocker(ctx)
				ctx = ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 31))
				app.EpochsKeeper.BeginBlocker(ctx)
				epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
				require.True(t, found)
			},
		},
		// Test that incrementing _exactly_ 1 month increments the epoch count.
		{
			expCurrentEpochStartHeight: 3,
			expCurrentEpochStartTime:   now.Add(time.Hour * 24 * 31),
			expCurrentEpoch:            2,
			expInitialEpochStartTime:   now,
			fn: func() {
				ctx = ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				app.EpochsKeeper.BeginBlocker(ctx)
				ctx = ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 32))
				app.EpochsKeeper.BeginBlocker(ctx)
				epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
				require.True(t, found)
			},
		},
		{
			expCurrentEpochStartHeight: 3,
			expCurrentEpochStartTime:   now.Add(time.Hour * 24 * 31),
			expCurrentEpoch:            2,
			expInitialEpochStartTime:   now,
			fn: func() {
				ctx = ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				app.EpochsKeeper.BeginBlocker(ctx)
				ctx = ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 32))
				app.EpochsKeeper.BeginBlocker(ctx)
				ctx.WithBlockHeight(4).WithBlockTime(now.Add(time.Hour * 24 * 33))
				app.EpochsKeeper.BeginBlocker(ctx)
				epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
				require.True(t, found)
			},
		},
		{
			expCurrentEpochStartHeight: 3,
			expCurrentEpochStartTime:   now.Add(time.Hour * 24 * 31),
			expCurrentEpoch:            2,
			expInitialEpochStartTime:   now,
			fn: func() {
				ctx = ctx.WithBlockHeight(2).WithBlockTime(now.Add(time.Second))
				app.EpochsKeeper.BeginBlocker(ctx)
				ctx = ctx.WithBlockHeight(3).WithBlockTime(now.Add(time.Hour * 24 * 32))
				app.EpochsKeeper.BeginBlocker(ctx)
				ctx.WithBlockHeight(4).WithBlockTime(now.Add(time.Hour * 24 * 33))
				app.EpochsKeeper.BeginBlocker(ctx)
				epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
				require.True(t, found)
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			// On init genesis, default epochs information is set
			// To check init genesis again, should make it fresh status

			epochInfos := app.EpochsKeeper.AllEpochInfos(ctx)
			for _, epochInfo := range epochInfos {
				app.EpochsKeeper.DeleteEpochInfo(ctx, epochInfo.Identifier)
			}

			ctx = ctx.WithBlockHeight(1).WithBlockTime(now)

			// check init genesis
			epochs.InitGenesis(ctx, app.EpochsKeeper, types.GenesisState{
				Epochs: []types.EpochInfo{
					{
						Identifier:              "monthly",
						StartTime:               time.Time{},
						Duration:                time.Hour * 24 * 31,
						CurrentEpoch:            0,
						CurrentEpochStartHeight: ctx.BlockHeight(),
						CurrentEpochStartTime:   time.Time{},
						EpochCountingStarted:    false,
					},
				},
			})

			tc.fn()

			require.Equal(t, epochInfo.Identifier, "monthly")
			require.Equal(t, epochInfo.StartTime.UTC().String(), tc.expInitialEpochStartTime.UTC().String())
			require.Equal(t, epochInfo.Duration, time.Hour*24*31)
			require.Equal(t, epochInfo.CurrentEpoch, tc.expCurrentEpoch)
			require.Equal(t, epochInfo.CurrentEpochStartHeight, tc.expCurrentEpochStartHeight)
			require.Equal(t, epochInfo.CurrentEpochStartTime.UTC().String(), tc.expCurrentEpochStartTime.UTC().String())
			require.Equal(t, epochInfo.EpochCountingStarted, true)
		})
	}
}

func TestEpochStartingOneMonthAfterInitGenesis(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)

	// On init genesis, default epochs information is set
	// To check init genesis again, should make it fresh status
	epochInfos := app.EpochsKeeper.AllEpochInfos(ctx)
	for _, epochInfo := range epochInfos {
		app.EpochsKeeper.DeleteEpochInfo(ctx, epochInfo.Identifier)
	}

	now := time.Now()
	week := time.Hour * 24 * 7
	month := time.Hour * 24 * 30
	initialBlockHeight := int64(1)
	ctx = ctx.WithBlockHeight(initialBlockHeight).WithBlockTime(now)

	epochs.InitGenesis(ctx, app.EpochsKeeper, types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               now.Add(month),
				Duration:                time.Hour * 24 * 30,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    false,
			},
		},
	})

	// epoch not started yet
	epochInfo, found := app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
	require.True(t, found)
	require.Equal(t, epochInfo.CurrentEpoch, int64(0))
	require.Equal(t, epochInfo.CurrentEpochStartHeight, initialBlockHeight)
	require.Equal(t, epochInfo.CurrentEpochStartTime, time.Time{})
	require.Equal(t, epochInfo.EpochCountingStarted, false)

	// after 1 week
	ctx = ctx.WithBlockHeight(2).WithBlockTime(now.Add(week))
	app.EpochsKeeper.BeginBlocker(ctx)

	// epoch not started yet
	epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
	require.True(t, found)
	require.Equal(t, epochInfo.CurrentEpoch, int64(0))
	require.Equal(t, epochInfo.CurrentEpochStartHeight, initialBlockHeight)
	require.Equal(t, epochInfo.CurrentEpochStartTime, time.Time{})
	require.Equal(t, epochInfo.EpochCountingStarted, false)

	// after 1 month
	ctx = ctx.WithBlockHeight(3).WithBlockTime(now.Add(month))
	app.EpochsKeeper.BeginBlocker(ctx)

	// epoch started
	epochInfo, found = app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
	require.True(t, found)
	require.Equal(t, epochInfo.CurrentEpoch, int64(1))
	require.Equal(t, epochInfo.CurrentEpochStartHeight, ctx.BlockHeight())
	require.Equal(t, epochInfo.CurrentEpochStartTime.UTC().String(), now.Add(month).UTC().String())
	require.Equal(t, epochInfo.EpochCountingStarted, true)
}
