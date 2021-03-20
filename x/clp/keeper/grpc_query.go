package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/clp/types"
)

var _ types.QueryServer = Keeper{}

func (q Keeper) QueryGetPool(ctx context.Context, req *types.QueryGetPoolRequest) (*types.QueryGetPoolResponse, error) {
	return nil, nil
}

func (q Keeper) QueryLiquidityProvider(ctx context.Context, req *types.QueryLiquidityProviderRequest) (*types.QueryLiquidityProviderResponse, error) {

	return nil, nil
}

func (q Keeper) QueryGetAssetList(ctx context.Context, req *types.QueryGetAssetListRequest) (*types.QueryGetAssetListResponse, error) {
	return nil, nil
}

func (q Keeper) QueryGetLiquidityProviderList(ctx context.Context, req *types.QueryGetLiquidityProviderListRequest) (*types.QueryGetLiquidityProviderListResponse, error) {
	return nil, nil
}
