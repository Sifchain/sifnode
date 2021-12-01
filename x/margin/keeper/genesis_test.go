package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sifchain/sifnode/x/margin/test"
)

func TestKeeperExportGenesis(t *testing.T) {
	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper
	assert.NotNil(t, marginKeeper)
	state := marginKeeper.ExportGenesis(ctx)
	assert.NotNil(t, state)
}

func TestKeeperInitGenesis(t *testing.T) {
	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper
	assert.NotNil(t, marginKeeper)
	marginKeeper.ExportGenesis(ctx)
	// state2 := marginKeeper.InitGenesis(ctx, *state)
	// assert.NotNil(t, state2)
}