package seed_test

import (
	"testing"

	"github.com/Sifchain/sifnode/seed"
	"github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSeed(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom: "ceth", Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom: "cusdc", Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom: "clink", Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
	})
	err := app.ClpKeeper.SetPool(ctx, &clptypes.Pool{
		ExternalAsset:        &clptypes.Asset{Symbol: "ceth"},
		NativeAssetBalance:   sdk.NewUint(1000000000),
		ExternalAssetBalance: sdk.NewUint(1000000000),
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &clptypes.Pool{
		ExternalAsset:        &clptypes.Asset{Symbol: "cusdc"},
		NativeAssetBalance:   sdk.NewUint(1000000000),
		ExternalAssetBalance: sdk.NewUint(1000000000),
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &clptypes.Pool{
		ExternalAsset:        &clptypes.Asset{Symbol: "clink"},
		NativeAssetBalance:   sdk.NewUint(1000000000),
		ExternalAssetBalance: sdk.NewUint(1000000000),
	})
	require.NoError(t, err)

	err = seed.Seed(app.ClpKeeper, app.BankKeeper, ctx, 10000, []string{"ceth", "cusdc", "clink"})
	require.NoError(t, err)

	lps, err := app.ClpKeeper.GetAllLiquidityProviders(ctx)
	require.NoError(t, err)
	require.Equal(t, 30000, len(lps))
}
