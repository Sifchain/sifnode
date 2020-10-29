package keeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_Errors(t *testing.T) {
	pool := GenerateRandomPool(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	_ = keeper.Logger(ctx)
	pool.ExternalAsset.Ticker = ""
	keeper.SetPool(ctx, pool)
	getpools := keeper.GetPools(ctx)
	assert.Equal(t, len(getpools), 0, "No pool added")

	lp := GenerateRandomLP(1)[0]
	lp.Asset.SourceChain = ""
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	assert.Error(t, err)
	assert.NotEqual(t, getlp, lp)
	assert.NotNil(t, GenerateAddress())
}

func TestKeeper_SetPool(t *testing.T) {

	pool := GenerateRandomPool(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetPool(ctx, pool)
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Ticker)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	assert.Equal(t, keeper.ExistsPool(ctx, pool.ExternalAsset.Ticker), true)
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
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Ticker)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	keeper.DestroyPool(ctx, pool.ExternalAsset.Ticker)
	_, err = keeper.GetPool(ctx, pool.ExternalAsset.Ticker)
	assert.Error(t, err, "Pool should be deleted")
	// This should do nothing.
	keeper.DestroyPool(ctx, pool.ExternalAsset.Ticker)
}

func TestKeeper_SetLiquidityProvider(t *testing.T) {
	lp := GenerateRandomLP(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
	lpList := keeper.GetLiqudityProvidersForAsset(ctx, lp.Asset)
	assert.Equal(t, lp, lpList[0])
}

func TestKeeper_DestroyLiquidityProvider(t *testing.T) {
	lp := GenerateRandomLP(1)[0]
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
	assert.True(t, keeper.GetLiquidityProviderIterator(ctx).Valid())
	keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	_, err = keeper.GetLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	assert.Error(t, err, "LiquidityProvider has been deleted")
	// This should do nothing
	keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	assert.False(t, keeper.GetLiquidityProviderIterator(ctx).Valid())
}
