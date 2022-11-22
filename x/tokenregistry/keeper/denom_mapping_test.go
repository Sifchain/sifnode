package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/stretchr/testify/require"
)

func TestGetPeggy2Denom(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)
	peggy1Denom := "peggy1"
	_, found := app.TokenRegistryKeeper.GetPeggy2Denom(ctx, peggy1Denom)
	// not found before set
	require.Equal(t, found, false)

	peggy2Denom := "peggy2"
	app.TokenRegistryKeeper.SetPeggy2Denom(ctx, peggy1Denom, peggy2Denom)
	peggy2DenomInKeeper, found := app.TokenRegistryKeeper.GetPeggy2Denom(ctx, peggy1Denom)
	// found after set
	require.Equal(t, found, true)
	require.Equal(t, peggy2Denom, peggy2DenomInKeeper)
}
