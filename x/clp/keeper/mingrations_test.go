package keeper_test

import (
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	tkrKeeper "github.com/Sifchain/sifnode/x/tokenregistry/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigrator_MigrateToVer3(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	migrationMap := tkrKeeper.GetDenomMigrationMap()
	expectedResults := make(map[string]bool)
	for peggy1denom, peggy2denom := range migrationMap {
		err := app.ClpKeeper.SetPool(ctx, &types.Pool{
			ExternalAsset:        &types.Asset{Symbol: peggy1denom},
			NativeAssetBalance:   sdk.NewUintFromString("10000000000000000000"),
			ExternalAssetBalance: sdk.NewUintFromString("70000000000000000000"),
			PoolUnits:            sdk.NewUintFromString("10000000000000000000"),
		})
		assert.NoError(t, err)
		expectedResults[peggy2denom] = true
	}

	migrator := keeper.NewMigrator(app.ClpKeeper)
	err := migrator.MigrateToVer3(ctx)
	assert.NoError(t, err)

	pools := app.ClpKeeper.GetPools(ctx)
	for _, pool := range pools {
		assert.True(t, expectedResults[pool.ExternalAsset.Symbol])
	}
}
