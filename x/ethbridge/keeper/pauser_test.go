package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/stretchr/testify/assert"
)

func TestSetPauser(t *testing.T) {
	ctx, app := test.CreateSimulatorApp(false)

	// test the default value before any setting
	paused := app.EthbridgeKeeper.IsPaused(ctx)
	assert.False(t, paused)

	// pause
	app.EthbridgeKeeper.SetPauser(ctx, &types.Pauser{
		IsPaused: true,
	})

	paused = app.EthbridgeKeeper.IsPaused(ctx)
	assert.True(t, paused)

	// unpause
	app.EthbridgeKeeper.SetPauser(ctx, &types.Pauser{
		IsPaused: false,
	})

	paused = app.EthbridgeKeeper.IsPaused(ctx)
	assert.False(t, paused)
}
