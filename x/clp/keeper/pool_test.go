package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"

	"github.com/Sifchain/sifnode/x/clp/test"
)

func TestKeeper_SetPool_ValidatePool(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	err := clpKeeper.SetPool(ctx, &pool)
	assert.NoError(t, err)
	getpool, err := clpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	assert.Equal(t, clpKeeper.ExistsPool(ctx, pool.ExternalAsset.Symbol), true)
	boolean := clpKeeper.ValidatePool(pool)
	assert.True(t, boolean)
}

func TestKeeper_GetPools(t *testing.T) {
	pools := test.GenerateRandomPool(10)
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	for i := range pools {
		pool := pools[i]
		err := clpKeeper.SetPool(ctx, &pool)
		assert.NoError(t, err)
	}
	getpools, _, err := clpKeeper.GetPoolsPaginated(ctx, &query.PageRequest{})
	assert.NoError(t, err)
	assert.Greater(t, len(getpools), 0, "More than one pool added")
	assert.LessOrEqual(t, len(getpools), len(pools), "Set pool will ignore duplicates")
}

func TestKeeper_DestroyPool(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	err := clpKeeper.SetPool(ctx, &pool)
	assert.NoError(t, err)
	getpool, err := clpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	err = clpKeeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	_, err = clpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.Error(t, err, "Pool should be deleted")
	// This should do nothing.
	err = clpKeeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	assert.Error(t, err)
}
