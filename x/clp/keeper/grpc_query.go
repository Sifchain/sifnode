package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sifchain/sifnode/x/clp/types"
)

const MaxPageLimit = 200

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper Keeper
}

var _ types.QueryServer = Querier{}

func (k Querier) GetPool(c context.Context, req *types.PoolReq) (*types.PoolRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	pool, err := k.Keeper.GetPool(ctx, req.Symbol)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.Symbol)
	}
	return &types.PoolRes{
		Pool:             &pool,
		Height:           ctx.BlockHeight(),
		ClpModuleAddress: types.GetCLPModuleAddress().String(),
	}, nil
}

func (k Querier) GetPools(c context.Context, req *types.PoolsReq) (*types.PoolsRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)

	pools, pageRes, err := k.Keeper.GetPoolsPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &types.PoolsRes{
		Pools:            pools,
		Height:           ctx.BlockHeight(),
		ClpModuleAddress: types.GetCLPModuleAddress().String(),
		Pagination:       pageRes,
	}, nil
}

func (k Querier) GetLiquidityProvider(c context.Context, req *types.LiquidityProviderReq) (*types.LiquidityProviderRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	lp, err := k.Keeper.GetLiquidityProvider(ctx, req.Symbol, req.LpAddress)
	if err != nil {
		return nil, err
	}
	pool, err := k.Keeper.GetPool(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	native, external, _, _ := CalculateAllAssetsForLP(pool, lp)
	lpResponse := types.NewLiquidityProviderResponse(lp, ctx.BlockHeight(), native.String(), external.String())
	return &lpResponse, nil
}

func (k Querier) GetAssetList(c context.Context, req *types.AssetListReq) (*types.AssetListRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.LpAddress)
	if err != nil {
		return nil, err
	}
	assetList, _, err := k.Keeper.GetAssetsForLiquidityProviderPaginated(ctx, addr, &query.PageRequest{Limit: MaxPageLimit})
	if err != nil {
		return nil, err
	}
	return &types.AssetListRes{
		Assets: assetList,
	}, nil
}

func (k Querier) GetLiquidityProviderList(c context.Context, req *types.LiquidityProviderListReq) (*types.LiquidityProviderListRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)
	searchingAsset := types.NewAsset(req.Symbol)

	lpList, pageRes, err := k.Keeper.GetLiquidityProvidersForAssetPaginated(ctx, searchingAsset, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.LiquidityProviderListRes{
		LiquidityProviders: lpList,
		Height:             ctx.BlockHeight(),
		Pagination:         pageRes,
	}, nil
}

func (k Querier) GetLiquidityProviders(c context.Context, req *types.LiquidityProvidersReq) (*types.LiquidityProvidersRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)

	lpList, pageRes, err := k.Keeper.GetAllLiquidityProvidersPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &types.LiquidityProvidersRes{
		LiquidityProviders: lpList,
		Height:             ctx.BlockHeight(),
		Pagination:         pageRes,
	}, nil
}
