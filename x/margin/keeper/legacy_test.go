package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/margin/keeper"
	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/stretchr/testify/require"
)

func TestKeeper_NewLegacyQuerier(t *testing.T) {
	_, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	got := keeper.NewLegacyQuerier(marginKeeper)

	require.NotNil(t, got)
}

func TestKeeper_NewLegacyHandler(t *testing.T) {
	_, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	got := keeper.NewLegacyHandler(marginKeeper)

	require.NotNil(t, got)
}
