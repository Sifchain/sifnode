package keeper_test

import (
	"testing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/clp"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestQueryErrorPool(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	cdc := app.Codec()
	keeper := app.ClpKeeper
	//Set Data
	pool, _, _ := SetData(keeper, ctx)
	querier := clpkeeper.NewQuerier(keeper)
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
	err = keeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	require.NoError(t, err)
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
	querier := clpkeeper.NewQuerier(keeper)
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
	var p types.PoolResponse
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
	querier := clpkeeper.NewQuerier(keeper)
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
	querier := clpkeeper.NewQuerier(keeper)
	query.Path = ""
	query.Data = nil
	//Test Pools
	qpools, err := querier(ctx, []string{types.QueryPools}, query)
	assert.NoError(t, err)
	var poolsRes types.PoolsResponse

	err = keeper.Codec().UnmarshalJSON(qpools, &poolsRes)
	assert.NoError(t, err)
	assert.Greater(t, len(poolsRes.Pools), 0, "More than one pool added")
	assert.LessOrEqual(t, len(poolsRes.Pools), len(pools), "Set pool will ignore duplicates")
}

func TestQueryErrorLiquidityProvider(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	cdc := app.Codec()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clpkeeper.NewQuerier(keeper)
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
	querier := clpkeeper.NewQuerier(keeper)
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
	var l types.LiquidityProviderResponse
	err = keeper.Codec().UnmarshalJSON(qliquidityprovider, &l)
	assert.NoError(t, err)
	assert.Equal(t, lp.Asset, l.Asset)

}

func SetData(keeper clpkeeper.Keeper, ctx sdk.Context) (types.Pool, []types.Pool, types.LiquidityProvider) {
	pool := test.GenerateRandomPool(1)[0]
	err := keeper.SetPool(ctx, pool)
	if err != nil {
		ctx.Logger().Error("Unable to set pool")
	}
	pools := test.GenerateRandomPool(10)
	for _, p := range pools {
		err = keeper.SetPool(ctx, p)
		if err != nil {
			ctx.Logger().Error("Unable to set pool")
		}
	}
	lp := test.GenerateRandomLP(1)[0]
	keeper.SetLiquidityProvider(ctx, lp)
	return pool, pools, lp
}
