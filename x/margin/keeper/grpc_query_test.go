package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/margin/keeper"
	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/stretchr/testify/require"
)

func TestKeeper_NewQueryServer(t *testing.T) {
	_, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	got := keeper.NewQueryServer(marginKeeper)

	require.NotNil(t, got)
}
