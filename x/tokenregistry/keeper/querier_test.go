package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/tokenregistry/keeper"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryEntries(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:    "rowan",
		Decimals: 18,
	})
	querier := keeper.NewLegacyQuerier(app.TokenRegistryKeeper)
	bz, err := app.AppCodec().MarshalJSON(&types.QueryEntriesRequest{})
	require.NoError(t, err)
	resBz, err := querier(ctx, []string{types.QueryEntries}, abci.RequestQuery{Data: bz})
	require.Nil(t, resBz)
	require.ErrorIs(t, err, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "Token Registry Legacy Querier No Longer Available"))
}
