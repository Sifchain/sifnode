package keeper_test

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.True(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{types.Permission_CLP}))
	assert.False(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{types.Permission_IBCIMPORT}))
	assert.False(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{types.Permission_CLP, types.Permission_IBCIMPORT}))
	assert.False(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "t2", []types.Permission{types.Permission_IBCEXPORT, types.Permission_IBCIMPORT}))
	assert.True(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{}))

}
