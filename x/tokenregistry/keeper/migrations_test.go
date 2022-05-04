package keeper_test

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/keeper"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	tkrtypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MigrateToVer4(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)

	tt := []struct {
		name        string
		denom       string
		permissions []tkrtypes.Permission
		isMigrating bool
	}{
		{
			name:  "TC1",
			denom: "ceth",
			permissions: []tkrtypes.Permission{
				tkrtypes.Permission_IBCIMPORT,
				tkrtypes.Permission_IBCEXPORT,
				tkrtypes.Permission_CLP,
			},
			isMigrating: true,
		},
		{
			name:  "TC1",
			denom: "cdash",
			permissions: []tkrtypes.Permission{
				tkrtypes.Permission_IBCIMPORT,
				tkrtypes.Permission_IBCEXPORT,
				tkrtypes.Permission_CLP,
			},
			isMigrating: false,
		},
	}
	// Test setup
	for _, s := range tt {
		app.TokenRegistryKeeper.SetToken(ctx, &tkrtypes.RegistryEntry{
			Denom:       s.denom,
			Permissions: s.permissions,
		})
	}
	migrator := keeper.NewMigrator(app.TokenRegistryKeeper)
	migrator.MigrateToVer4(ctx)

	for _, s := range tt {
		tc := s
		t.Run(tc.name, func(t *testing.T) {
			entry, err := app.TokenRegistryKeeper.GetRegistryEntry(
				ctx,
				tc.denom,
			)
			assert.NoError(t, err)
			if tc.isMigrating {
				assert.Equal(t,
					[]tkrtypes.Permission{tkrtypes.Permission_IBCIMPORT},
					entry.Permissions,
				)
			} else {
				assert.Equal(t,
					tc.permissions,
					entry.Permissions,
				)
			}
		})
	}
}
