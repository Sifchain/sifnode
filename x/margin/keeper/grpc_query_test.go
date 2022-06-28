//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/margin/keeper"
	"github.com/Sifchain/sifnode/x/margin/test"
)

func TestKeeper_NewQueryServer(t *testing.T) {
	_, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	_ = keeper.NewQueryServer(marginKeeper)

	// FIXME: query server not implemented yet
	// require.Nil(t, got)
}
