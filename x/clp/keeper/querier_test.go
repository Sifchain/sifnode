package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestNewQuerier(t *testing.T) {
	cdc := codec.New()
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	//Set Data
	pool, pools, lp := SetData(keeper, ctx)
	querier := NewQuerier(keeper)

	//Test Pool
	queryPool := types.QueryReqGetPool{
		Ticker:      pool.ExternalAsset.Ticker,
		SourceChain: pool.ExternalAsset.SourceChain,
	}
	qp, errRes := cdc.MarshalJSON(queryPool)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	qpool, err := querier(ctx, []string{"pool"}, query)
	assert.NoError(t, err)
	var p types.Pool
	err = keeper.cdc.UnmarshalJSON(qpool, &p)
	assert.NoError(t, err)
	assert.Equal(t, pool.ExternalAsset, p.ExternalAsset)

	//Test Pools
	query.Path = ""
	query.Data = nil
	qpools, err := querier(ctx, []string{"allpools"}, query)
	assert.NoError(t, err)
	var poolist []types.Pool
	err = keeper.cdc.UnmarshalJSON(qpools, &poolist)
	assert.NoError(t, err)
	assert.Greater(t, len(poolist), 0, "More than one pool added")
	assert.LessOrEqual(t, len(poolist), len(pools), "Set pool will ignore duplicates")

	//Test Get Liquidity Provider
	queryLp := types.QueryReqLiquidityProvider{
		Ticker: lp.Asset.Ticker,
		Ip:     lp.LiquidityProviderAddress,
	}
	qlp, errRes := cdc.MarshalJSON(queryLp)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qlp
	qliquidityprovider, err := querier(ctx, []string{"liquidityProvider"}, query)
	assert.NoError(t, err)
	var l types.LiquidityProvider
	err = keeper.cdc.UnmarshalJSON(qliquidityprovider, &l)
	assert.NoError(t, err)
	assert.Equal(t, lp.Asset, l.Asset)

}

func SetData(keeper Keeper, ctx sdk.Context) (types.Pool, []types.Pool, types.LiquidityProvider) {
	pool := generateRandomPool(1)[0]
	keeper.SetPool(ctx, pool)

	pools := generateRandomPool(10)
	for _, p := range pools {
		keeper.SetPool(ctx, p)
	}

	lp := generateRandomLP(1)[0]
	keeper.SetLiquidityProvider(ctx, lp)
	return pool, pools, lp

}
