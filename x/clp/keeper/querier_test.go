package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestQueryErrorPool(t *testing.T) {
	cdc, app, ctx := createTestInput()

	keeper := app.ClpKeeper
	//Set Data
	pool, _, _ := SetData(keeper, ctx)
	querier := clpkeeper.NewQuerier(keeper, cdc)
	//Test Pool
	queryPool := types.PoolReq{
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
	cdc, app, ctx := createTestInput()

	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	pool, _, _ := SetData(keeper, ctx)
	querier := clpkeeper.NewQuerier(keeper, cdc)
	//Test Pool
	queryPool := types.PoolReq{
		Symbol: pool.ExternalAsset.Symbol,
	}
	qp, errRes := cdc.MarshalJSON(queryPool)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	qpool, err := querier(ctx, []string{types.QueryPool}, query)
	assert.NoError(t, err)
	var p types.PoolResponse
	err = cdc.UnmarshalJSON(qpool, &p)
	assert.NoError(t, err)
	assert.Equal(t, pool.ExternalAsset, p.ExternalAsset)
}

func TestQueryErrorPools(t *testing.T) {
	cdc, app, ctx := createTestInput()

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clpkeeper.NewQuerier(app.ClpKeeper, cdc)
	query.Path = ""
	query.Data = nil
	//Test Pools
	_, err := querier(ctx, []string{types.QueryPools}, query)
	assert.Error(t, err)
}

func TestQueryGetPools(t *testing.T) {
	cdc, app, ctx := createTestInput()

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	_, pools, _ := SetData(app.ClpKeeper, ctx)
	querier := clpkeeper.NewQuerier(app.ClpKeeper, cdc)
	query.Path = ""
	query.Data = nil
	//Test Pools
	qpools, err := querier(ctx, []string{types.QueryPools}, query)
	assert.NoError(t, err)
	var poolsRes types.PoolsResponse

	err = cdc.UnmarshalJSON(qpools, &poolsRes)
	assert.NoError(t, err)
	assert.Greater(t, len(poolsRes.Pools), 0, "More than one pool added")
	assert.LessOrEqual(t, len(poolsRes.Pools), len(pools), "Set pool will ignore duplicates")
}

func TestQueryErrorLiquidityProvider(t *testing.T) {
	cdc, app, ctx := createTestInput()

	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clpkeeper.NewQuerier(keeper, cdc)
	_, err := querier(ctx, []string{types.QueryLiquidityProvider}, query)
	assert.Error(t, err)
	//Set Data
	_, _, lp := SetData(keeper, ctx)
	//Test Get Liquidity Provider

	addr, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	assert.NoError(t, err)

	queryLp := types.LiquidityProviderReq{
		Symbol:    "", //lp.Asset.Ticker,
		LpAddress: addr.String(),
	}
	qlp, errRes := cdc.MarshalJSON(queryLp)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qlp
	_, err = querier(ctx, []string{types.QueryLiquidityProvider}, query)
	assert.Error(t, err)
}

func TestQueryGetLiquidityProvider(t *testing.T) {
	cdc, app, ctx := createTestInput()

	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	_, _, lp := SetData(keeper, ctx)
	querier := clpkeeper.NewQuerier(keeper, cdc)
	//Test Get Liquidity Provider
	addr, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	assert.NoError(t, err)

	queryLp := types.LiquidityProviderReq{
		Symbol:    lp.Asset.Symbol,
		LpAddress: addr.String(),
	}
	qlp, errRes := cdc.MarshalJSON(queryLp)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qlp
	qliquidityprovider, err := querier(ctx, []string{types.QueryLiquidityProvider}, query)
	assert.NoError(t, err)
	var l types.LiquidityProviderRes
	err = cdc.UnmarshalJSON(qliquidityprovider, &l)
	assert.NoError(t, err)
	assert.Equal(t, lp.Asset, l.LiquidityProvider.Asset)

}

func SetData(keeper clpkeeper.Keeper, ctx sdk.Context) (types.Pool, []types.Pool, types.LiquidityProvider) {
	pool := test.GenerateRandomPool(1)[0]
	err := keeper.SetPool(ctx, &pool)
	if err != nil {
		ctx.Logger().Error("Unable to set pool")
	}
	pools := test.GenerateRandomPool(10)
	for i := range pools {
		p := pools[i]
		err = keeper.SetPool(ctx, &p)
		if err != nil {
			ctx.Logger().Error("Unable to set pool")
		}
	}
	lp := test.GenerateRandomLP(1)[0]
	keeper.SetLiquidityProvider(ctx, &lp)
	return pool, pools, lp
}
