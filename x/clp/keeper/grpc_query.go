package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/clp/types"
)

var _ types.QueryServer = Keeper{}

func (q Keeper) QueryGetPool(ctx context.Context, req *types.QueryReqGetPoolRequest) (*types.QueryReqGetPoolResponse, error) {
	return nil, nil
}

func (q Keeper) QueryLiquidityProvider(ctx context.Context, req *types.QueryReqLiquidityProviderRequest) (*types.QueryReqLiquidityProviderResponse, error) {

	return nil, nil
}

func (q Keeper) QueryGetAssetList(ctx context.Context, req *types.QueryReqGetAssetListRequest) (*types.QueryReqGetAssetListResponse, error) {
	return nil, nil
}

func (q Keeper) QueryGetLiquidityProviderList(ctx context.Context, req *types.QueryReqGetLiquidityProviderListRequest) (*types.QueryReqGetLiquidityProviderListResponse, error) {
	return nil, nil
}
