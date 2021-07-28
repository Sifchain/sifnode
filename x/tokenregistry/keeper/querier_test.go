package keeper_test

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/keeper"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestQueryEntries(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)

	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:    "rowan",
		Decimals: 18,
	})

	expectedRegistry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)

	querier := keeper.NewLegacyQuerier(app.TokenRegistryKeeper)

	bz, err := app.AppCodec().MarshalJSON(&types.QueryEntriesRequest{})
	require.NoError(t, err)

	resBz, err := querier(ctx, []string{types.QueryEntries}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	res := types.QueryEntriesResponse{}

	app.AppCodec().MustUnmarshalJSON(resBz, &res)
	require.Len(t, res.Registry.Entries, 1)
	require.Equal(t, &expectedRegistry, res.Registry)
}
