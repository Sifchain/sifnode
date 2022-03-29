package clp_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/stretchr/testify/require"
)

func TestBeginBlocker(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clp.BeginBlocker(ctx, app.ClpKeeper)
}

func TestBeginBlocker_ZeroCounters(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)

	// params := types.DefaultParams()
	// params.PmtpPeriodStartBlock = int64(0)

	// app.ClpKeeper.SetParams(ctx, params)

	require.Equal(t, ctx.BlockHeight(), int64(0))
	require.Equal(t, app.ClpKeeper.GetPmtpStartBlock(ctx), int64(1))

	epochParams := types.PmtpEpoch{
		EpochCounter: 0,
		BlockCounter: 0,
	}
	app.ClpKeeper.SetPmtpEpoch(ctx, epochParams)
	clp.BeginBlocker(ctx, app.ClpKeeper)
}
