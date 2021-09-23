package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_CheckDenomPermissions(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)

	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "rowan",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_PERMISSION_CLP},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "t1",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_PERMISSION_UNSPECIFIED},
	})
	// Duplicate permission is interpreted correctly.
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:       "t2",
		Decimals:    18,
		Permissions: []types.Permission{types.Permission_PERMISSION_IBCEXPORT, types.Permission_PERMISSION_IBCEXPORT},
	})
	assert.True(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{types.Permission_PERMISSION_CLP}))
	assert.False(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{types.Permission_PERMISSION_IBCIMPORT}))
	assert.False(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{types.Permission_PERMISSION_CLP, types.Permission_PERMISSION_IBCIMPORT}))
	assert.False(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "t2", []types.Permission{types.Permission_PERMISSION_IBCEXPORT, types.Permission_PERMISSION_IBCIMPORT}))
	assert.True(t, app.TokenRegistryKeeper.CheckDenomPermissions(ctx, "rowan", []types.Permission{}))

}
