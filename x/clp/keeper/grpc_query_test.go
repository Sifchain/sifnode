package keeper_test

import (
	"context"
	"errors"
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestQuerier_GetPool(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetPool(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetPools(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetPools(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetPools_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.PoolsReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetPools(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProvider(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProvider(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviderData(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProviderData(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetAssetList(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetAssetList(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetAssetList_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.AssetListReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetAssetList(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviderList(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProviderList(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviderList_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.LiquidityProviderListReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetLiquidityProviderList(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviders(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	_, err := querier.GetLiquidityProviders(ctx, nil)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}

func TestQuerier_GetLiquidityProviders_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.LiquidityProvidersReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	_, err := querier.GetLiquidityProviders(ctx, req)
	require.Error(t, err, errors.New("rpc error: code = InvalidArgument desc = empty request"))
}
