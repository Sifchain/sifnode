package keeper_test

import (
	"testing"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_CheckDenomPermissions(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "rowan",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_CLP},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "t1",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_UNSPECIFIED},
	})
	// Duplicate permission is interpreted correctly.
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "t2",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_IBCEXPORT, types.Permission_IBCEXPORT},
	})
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	entry, err := app.TokenRegistryKeeper.GetEntry(registry, "rowan")
	assert.NoError(t, err)
	entry2, err := app.TokenRegistryKeeper.GetEntry(registry, "t2")
	assert.NoError(t, err)
	assert.True(t, app.TokenRegistryKeeper.CheckEntryPermissions(entry, []types.Permission{types.Permission_CLP}))
	assert.False(t, app.TokenRegistryKeeper.CheckEntryPermissions(entry, []types.Permission{types.Permission_IBCIMPORT}))
	assert.False(t, app.TokenRegistryKeeper.CheckEntryPermissions(entry, []types.Permission{types.Permission_CLP, types.Permission_IBCIMPORT}))
	assert.False(t, app.TokenRegistryKeeper.CheckEntryPermissions(entry2, []types.Permission{types.Permission_IBCEXPORT, types.Permission_IBCIMPORT}))
	assert.True(t, app.TokenRegistryKeeper.CheckEntryPermissions(entry, []types.Permission{}))
}

func TestKeeper_AddRemoveRegisterAll(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "rowan",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_CLP},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "t1",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_UNSPECIFIED},
	})

	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "t2",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_IBCEXPORT, types.Permission_IBCEXPORT},
	})
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	assert.Equal(t, len(registry.Entries), 3)

	entries := []*types.RegistryEntry{}
	entries = append(entries, &types.RegistryEntry{
		Denom:       "t2",
		Decimals:    19,
		Permissions: []types.Permission{types.Permission_IBCEXPORT, types.Permission_IBCEXPORT},
	})

	entries = append(entries, &types.RegistryEntry{
		Denom:       "t3",
		Decimals:    19,
		Permissions: []types.Permission{types.Permission_IBCEXPORT, types.Permission_IBCEXPORT},
	})

	// add entries
	app.TokenRegistryKeeper.AddMultipleTokens(ctx, entries)

	registry = app.TokenRegistryKeeper.GetRegistry(ctx)

	assert.Equal(t, len(registry.Entries), 4)

	// remove entries
	denom := []string{"t2", "t3", "t4"}
	app.TokenRegistryKeeper.RemoveMultipleTokens(ctx, denom)

	registry = app.TokenRegistryKeeper.GetRegistry(ctx)

	assert.Equal(t, len(registry.Entries), 2)
}

func TestKeeper_SetFirstLockDoublePeg(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)

	denom := "denom"
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       denom,
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_CLP},
	})
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_GANACHE

	assert.True(t, app.TokenRegistryKeeper.GetFirstLockDoublePeg(ctx, denom, networkDescriptor))
	app.TokenRegistryKeeper.SetFirstLockDoublePeg(ctx, denom, networkDescriptor)
	assert.False(t, app.TokenRegistryKeeper.GetFirstLockDoublePeg(ctx, denom, networkDescriptor))

}
