package keeper_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkQuery "github.com/cosmos/cosmos-sdk/types/query"

	sifapp "github.com/Sifchain/sifnode/app"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// createTestInput Returns a simapp with custom StakingKeeper
// to avoid messing with the hooks.
func createTestInput() (*codec.LegacyAmino, *sifapp.SifchainApp, sdk.Context) { //nolint
	sifapp.SetConfig(false)
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.ClpKeeper = clpkeeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.BankKeeper,
		app.AccountKeeper,
		app.TokenRegistryKeeper,
		app.AdminKeeper,
		app.MintKeeper,
		func() margintypes.Keeper { return app.MarginKeeper },
		app.GetSubspace(types.ModuleName),
	)
	return app.LegacyAmino(), app, ctx
}

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
	var p types.PoolRes
	err = cdc.UnmarshalJSON(qpool, &p)
	assert.NoError(t, err)
	assert.Equal(t, pool.ExternalAsset, p.Pool.ExternalAsset)
}

func TestQueryEmptyPools(t *testing.T) {
	cdc, app, ctx := createTestInput()
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clpkeeper.NewQuerier(app.ClpKeeper, cdc)
	query.Path = ""
	query.Data = nil
	//Test Pools
	qpools, err := querier(ctx, []string{types.QueryPools}, query)
	assert.NoError(t, err)
	var poolsRes types.PoolsRes
	err = cdc.UnmarshalJSON(qpools, &poolsRes)
	assert.NoError(t, err)
	assert.Empty(t, poolsRes.Pools)
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
	var poolsRes types.PoolsRes
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

func TestQueryGetLiquidityProviderData(t *testing.T) {
	cdc, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	tokens := []string{"cada", "cbch", "cbnb", "cbtc", "ceos", "ceth", "ctrx", "cusdt"}
	pools, lpList := test.GeneratePoolsAndLPs(keeper, ctx, tokens)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clpkeeper.NewQuerier(keeper, cdc)
	addr, err := sdk.AccAddressFromBech32(lpList[0].LiquidityProviderAddress)
	assert.NoError(t, err)
	queryLp := types.LiquidityProviderDataReq{
		LpAddress: addr.String(),
	}
	qlp, errRes := cdc.MarshalJSON(queryLp)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qlp
	qliquidityprovider, err := querier(ctx, []string{types.QueryLiquidityProviderData}, query)
	assert.NoError(t, err)
	var l types.LiquidityProviderDataRes
	err = cdc.UnmarshalJSON(qliquidityprovider, &l)
	assert.NoError(t, err)
	assert.Equal(t, len(pools), len(l.LiquidityProviderData))
	assert.Equal(t, len(lpList), len(l.LiquidityProviderData))
	assert.Equal(t, lpList[0].LiquidityProviderAddress, l.LiquidityProviderData[0].LiquidityProvider.LiquidityProviderAddress)
	for i := range l.LiquidityProviderData {
		lpData := l.LiquidityProviderData[i]
		assert.Contains(t, lpList, *lpData.LiquidityProvider)
		assert.Equal(t, lpList[0].LiquidityProviderAddress, lpData.LiquidityProvider.LiquidityProviderAddress)
		assert.Equal(t, test.TrimFirstRune(tokens[i]), lpData.LiquidityProvider.Asset.Symbol)
		assert.Equal(t, fmt.Sprint(100*uint64(i+1)), lpData.ExternalAssetBalance)
		assert.Equal(t, fmt.Sprint(1000*uint64(i+1)), lpData.NativeAssetBalance)
	}
}

func TestQueryGetLiquidityProviderDataFailure(t *testing.T) {
	cdc, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	tokens := []string{"cada", "cbch", "cbnb", "cbtc", "ceos", "ceth", "ctrx", "cusdt"}
	_, lpList := test.GeneratePoolsAndLPs(keeper, ctx, tokens)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := clpkeeper.NewQuerier(keeper, cdc)
	addr, err := sdk.AccAddressFromBech32(lpList[0].LiquidityProviderAddress)
	assert.NoError(t, err)
	query.Path = ""
	query.Data = nil
	_, err = querier(ctx, []string{types.QueryLiquidityProviderData}, query)
	assert.Error(t, err)
	queryLp := types.LiquidityProviderDataReq{
		LpAddress: addr.String(),
		Pagination: &sdkQuery.PageRequest{
			Limit: 300,
		},
	}
	qlp, errRes := cdc.MarshalJSON(queryLp)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qlp
	_, err = querier(ctx, []string{types.QueryLiquidityProviderData}, query)
	assert.Error(t, err)
}

func TestQueryAssetList(t *testing.T) {
	cdc, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	_, _, lp := SetData(keeper, ctx)
	querier := clpkeeper.NewQuerier(keeper, cdc)
	req := types.AssetListReq{
		LpAddress: lp.LiquidityProviderAddress,
	}
	queryData, err := cdc.MarshalJSON(req)
	require.NoError(t, err)
	query.Data = queryData
	resBz, err := querier(ctx, []string{types.QueryAssetList}, query)
	require.NoError(t, err)
	res := types.AssetListRes{}
	err = cdc.UnmarshalJSON(resBz, &res.Assets)
	require.NoError(t, err)
	require.Equal(t, []*types.Asset{lp.Asset}, res.Assets)
}

func TestQueryLPList(t *testing.T) {
	cdc, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	_, _, lp := SetData(keeper, ctx)
	querier := clpkeeper.NewQuerier(keeper, cdc)
	req := types.LiquidityProviderListReq{
		Symbol: lp.Asset.Symbol,
	}
	queryData, errRes := cdc.MarshalJSON(req)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = queryData
	qLpList, err := querier(ctx, []string{types.QueryLPList}, query)
	assert.NoError(t, err)
	res := types.LiquidityProviderListRes{}
	err = cdc.UnmarshalJSON(qLpList, &res.LiquidityProviders)
	assert.NoError(t, err)
	require.Equal(t, []*types.LiquidityProvider{&lp}, res.LiquidityProviders)
}

func TestQueryAllLPs(t *testing.T) {
	cdc, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	_, _, lp := SetData(keeper, ctx)
	querier := clpkeeper.NewQuerier(keeper, cdc)
	query.Path = ""
	query.Data = nil
	resBz, err := querier(ctx, []string{types.QueryAllLP}, query)
	assert.NoError(t, err)
	var lpRes types.LiquidityProvidersRes
	err = cdc.UnmarshalJSON(resBz, &lpRes.LiquidityProviders)
	assert.NoError(t, err)
	assert.Equal(t, []*types.LiquidityProvider{&lp}, lpRes.LiquidityProviders)
}

func TestQueryPmtpParams(t *testing.T) {
	cdc, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	//Set Data
	querier := clpkeeper.NewQuerier(keeper, cdc)
	query.Path = ""
	query.Data = nil
	resBz, err := querier(ctx, []string{types.QueryPmtpParams}, query)
	assert.NoError(t, err)
	var pmtpParamsRes types.PmtpParamsRes
	err = cdc.UnmarshalJSON(resBz, &pmtpParamsRes)
	assert.NoError(t, err)
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
	keeper.SetLiquidityProvider(ctx, lp)
	return pool, pools, *lp
}
