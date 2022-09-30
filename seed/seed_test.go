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
		Decimals: 6,
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom: "cusdc", Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
		Decimals: 6,
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom: "clink", Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
		Decimals: 6,
	})
	swapPriceExternal, err := sdk.NewDecFromStr("127.230711704805236001")
	require.NoError(t, err)
	swapPriceNative, err := sdk.NewDecFromStr("0.007859737531926673")
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &clptypes.Pool{
		ExternalAsset:        &clptypes.Asset{Symbol: "ceth"},
		NativeAssetBalance:   sdk.NewUintFromString("75151690408573637255087195"),
		ExternalAssetBalance: sdk.NewUintFromString("590672561692"),
		SwapPriceNative:      &swapPriceNative,
		SwapPriceExternal:    &swapPriceExternal,
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &clptypes.Pool{
		ExternalAsset:        &clptypes.Asset{Symbol: "cusdc"},
		NativeAssetBalance:   sdk.NewUintFromString("75151690408573637255087195"),
		ExternalAssetBalance: sdk.NewUintFromString("590672561692"),
		SwapPriceNative:      &swapPriceNative,
		SwapPriceExternal:    &swapPriceExternal,
	})
	require.NoError(t, err)
	err = app.ClpKeeper.SetPool(ctx, &clptypes.Pool{
		ExternalAsset:        &clptypes.Asset{Symbol: "clink"},
		NativeAssetBalance:   sdk.NewUintFromString("75151690408573637255087195"),
		ExternalAssetBalance: sdk.NewUintFromString("590672561692"),
		SwapPriceNative:      &swapPriceNative,
		SwapPriceExternal:    &swapPriceExternal,
	})
	require.NoError(t, err)

	err = seed.Seed(app.ClpKeeper, app.BankKeeper, app.TokenRegistryKeeper, ctx, 10000, 3)
	require.NoError(t, err)

	lps, err := app.ClpKeeper.GetAllLiquidityProviders(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, 10000, len(lps))
}
