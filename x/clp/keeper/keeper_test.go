package keeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_SetPool(t *testing.T) {

	pool := GenerateRandomPool(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetPool(ctx, pool)
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Ticker, pool.ExternalAsset.SourceChain)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
}

func TestKeeper_GetPools(t *testing.T) {
	pools := GenerateRandomPool(10)
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	for _, pool := range pools {
		keeper.SetPool(ctx, pool)
	}
	getpools := keeper.GetPools(ctx)
	assert.Greater(t, len(getpools), 0, "More than one pool added")
	assert.LessOrEqual(t, len(getpools), len(pools), "Set pool will ignore duplicates")
}

func TestKeeper_DestroyPool(t *testing.T) {
	pool := GenerateRandomPool(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetPool(ctx, pool)
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Ticker, pool.ExternalAsset.SourceChain)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	keeper.DestroyPool(ctx, pool.ExternalAsset.Ticker, pool.ExternalAsset.SourceChain)
	_, err = keeper.GetPool(ctx, pool.ExternalAsset.Ticker, pool.ExternalAsset.SourceChain)
	assert.Error(t, err, "Pool should be deleted")
}

func TestKeeper_SetLiquidityProvider(t *testing.T) {
	lp := GenerateRandomLP(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
}

func TestKeeper_DestroyLiquidityProvider(t *testing.T) {
	lp := GenerateRandomLP(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
	keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
	_, err = keeper.GetLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
	assert.Error(t, err, "LiquidityProvider has been deleted")
}
