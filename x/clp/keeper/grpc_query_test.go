package keeper_test

import (
	"context"
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func TestQuerier_GetPool(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	querier.GetPool(ctx, nil)
}

func TestQuerier_GetPools(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	querier.GetPools(ctx, nil)
}

func TestQuerier_GetPools_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.PoolsReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	querier.GetPools(ctx, req)
}

func TestQuerier_GetLiquidityProvider(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	querier.GetLiquidityProvider(ctx, nil)
}

func TestQuerier_GetLiquidityProviderData(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	querier.GetLiquidityProviderData(ctx, nil)
}

func TestQuerier_GetAssetList(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	querier.GetAssetList(ctx, nil)
}

func TestQuerier_GetAssetList_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.AssetListReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	querier.GetAssetList(ctx, req)
}

func TestQuerier_GetLiquidityProviderList(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	querier.GetLiquidityProviderList(ctx, nil)
}

func TestQuerier_GetLiquidityProviderList_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.LiquidityProviderListReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	querier.GetLiquidityProviderList(ctx, req)
}

func TestQuerier_GetLiquidityProviders(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	querier.GetLiquidityProviders(ctx, nil)
}

func TestQuerier_GetLiquidityProviders_ReachedLimit(t *testing.T) {
	var ctx context.Context
	querier := clpkeeper.Querier{}

	req := &types.LiquidityProvidersReq{
		Pagination: &query.PageRequest{
			Limit: clpkeeper.MaxPageLimit + 1,
		},
	}

	querier.GetLiquidityProviders(ctx, req)
}
