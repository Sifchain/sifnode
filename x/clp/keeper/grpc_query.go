package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sifchain/sifnode/x/clp/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
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

	ctx := sdk.UnwrapSDKContext(c)

	addr, err := sdk.AccAddressFromBech32(req.LpAddress)
	if err != nil {
		return nil, err
	}

	assetList := k.GetAssetsForLiquidityProvider(ctx, addr)

	al := make([]*types.Asset, len(assetList))

	for i := range assetList {
		asset := assetList[i]
		al = append(al, &asset)
	}

	return &types.AssetListRes{
		Assets: al,
	}, nil
}

func (k Querier) GetLiquidityProviderList(c context.Context, req *types.LiquidityProviderListReq) (*types.LiquidityProviderListRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	searchingAsset := types.NewAsset(req.Symbol)
	lpList, pageRes, err := k.GetLiquidityProvidersForAssetPaginated(ctx, searchingAsset, req.Pagination)
	if err != nil {
		return nil, err
	}

	liquidityProviders := make([]*types.LiquidityProvider, len(lpList))

	for i, lp := range lpList {
		lp := lp
		liquidityProviders[i] = &lp
	}

	return &types.LiquidityProviderListRes{
		LiquidityProviders: liquidityProviders,
		Height:             ctx.BlockHeight(),
		Pagination:         pageRes,
	}, nil
}

func (k Querier) GetLiquidityProviders(c context.Context, req *types.LiquidityProvidersReq) (*types.LiquidityProvidersRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var lpl []*types.LiquidityProvider
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, types.LiquidityProviderPrefix)

	pageRes, err := query.FilteredPaginate(valStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var lp types.LiquidityProvider
		err := k.cdc.UnmarshalBinaryBare(value, &lp)
		if err != nil {
			return false, err
		}

		if accumulate {
			lpl = append(lpl, &lp)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.LiquidityProvidersRes{
		LiquidityProviders: lpl,
		Height:             ctx.BlockHeight(),
		Pagination:         pageRes,
	}, nil
}
