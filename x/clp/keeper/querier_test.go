package keeper_test

import (
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestQueryErrorPool(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	cdc := app.Codec()
	keeper := app.ClpKeeper
	//Set Data
	pool, _, _ := SetData(keeper, ctx)
	querier := clp.NewQuerier(keeper)
	//Test Pool
	queryPool := types.QueryReqGetPool{
		Symbol: pool.ExternalAsset.Symbol,
	}
	qp, errRes := cdc.MarshalJSON(queryPool)
	require.NoError(t, errRes)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, err := querier(ctx, []string{"bogus"}, query)
	assert.Error(t, err)
	_, err = querier(ctx, []string{types.QueryPool}, query)
	assert.Error(t, err)
	keeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	query.Path = ""
	query.Data = qp
	_, err = querier(ctx, []string{types.QueryPool}, query)
	// Should fail after it is deleted.
	assert.Error(t, err)
}

func TestQueryGetPool(t *testing.T) {
	//cdc := codec.New()
	app, ctx := test.CreateTestApp(false)
	cdc := app.Codec()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	pool, _, _ := SetData(keeper, ctx)
	querier := clp.NewQuerier(keeper)
	//Test Pool
	queryPool := types.QueryReqGetPool{
		Symbol: pool.ExternalAsset.Symbol,
	}
	qp, errRes := cdc.MarshalJSON(queryPool)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	qpool, err := querier(ctx, []string{types.QueryPool}, query)
	assert.NoError(t, err)
	var p types.Pool
	err = keeper.Codec().UnmarshalJSON(qpool, &p)
	assert.NoError(t, err)
	assert.Equal(t, pool.ExternalAsset, p.ExternalAsset)
}

func TestQueryErrorPools(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clp.NewQuerier(keeper)
	query.Path = ""
	query.Data = nil
	//Test Pools
	_, err := querier(ctx, []string{types.QueryPools}, query)
	assert.Error(t, err)
}

func TestQueryGetPools(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	_, pools, _ := SetData(keeper, ctx)
	querier := clp.NewQuerier(keeper)
	query.Path = ""
	query.Data = nil
	//Test Pools
	qpools, err := querier(ctx, []string{types.QueryPools}, query)
	assert.NoError(t, err)
	var poolist []types.Pool
	err = keeper.Codec().UnmarshalJSON(qpools, &poolist)
	assert.NoError(t, err)
	assert.Greater(t, len(poolist), 0, "More than one pool added")
	assert.LessOrEqual(t, len(poolist), len(pools), "Set pool will ignore duplicates")
}

func TestQueryErrorLiquidityProvider(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	cdc := app.Codec()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clp.NewQuerier(keeper)
	_, err := querier(ctx, []string{types.QueryLiquidityProvider}, query)
	assert.Error(t, err)
	//Set Data
	_, _, lp := SetData(keeper, ctx)
	//Test Get Liquidity Provider

	queryLp := types.QueryReqLiquidityProvider{
		Symbol:    "", //lp.Asset.Ticker,
		LpAddress: lp.LiquidityProviderAddress,
	}
	qlp, errRes := cdc.MarshalJSON(queryLp)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qlp
	_, err = querier(ctx, []string{types.QueryLiquidityProvider}, query)
	assert.Error(t, err)
}

func TestQueryGetLiquidityProvider(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	cdc := app.Codec()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	_, _, lp := SetData(keeper, ctx)
	querier := clp.NewQuerier(keeper)
	//Test Get Liquidity Provider
	queryLp := types.QueryReqLiquidityProvider{
		Symbol:    lp.Asset.Symbol,
		LpAddress: lp.LiquidityProviderAddress,
	}
	qlp, errRes := cdc.MarshalJSON(queryLp)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qlp
	qliquidityprovider, err := querier(ctx, []string{types.QueryLiquidityProvider}, query)
	assert.NoError(t, err)
	var l types.LiquidityProvider
	err = keeper.Codec().UnmarshalJSON(qliquidityprovider, &l)
	assert.NoError(t, err)
	assert.Equal(t, lp.Asset, l.Asset)

}

func SetData(keeper clp.Keeper, ctx sdk.Context) (types.Pool, []types.Pool, types.LiquidityProvider) {
	pool := test.GenerateRandomPool(1)[0]
	keeper.SetPool(ctx, pool)
	pools := test.GenerateRandomPool(10)
	for _, p := range pools {
		keeper.SetPool(ctx, p)
	}
	lp := test.GenerateRandomLP(1)[0]
	keeper.SetLiquidityProvider(ctx, lp)
	return pool, pools, lp
}
