package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (keeper Keeper) SetLiquidityProvider(ctx sdk.Context, lp *types.LiquidityProvider) {
	if !lp.Validate() {
		return
	}
	store := ctx.KVStore(keeper.storeKey)
	key := types.GetLiquidityProviderKey(lp.Asset.Symbol, lp.LiquidityProviderAddress)
	store.Set(key, keeper.cdc.MustMarshal(lp))
}

func (keeper Keeper) GetLiquidityProvider(ctx sdk.Context, symbol string, lpAddress string) (types.LiquidityProvider, error) {
	var lp types.LiquidityProvider
	key := types.GetLiquidityProviderKey(symbol, lpAddress)
	store := ctx.KVStore(keeper.storeKey)
	if !keeper.Exists(ctx, key) {
		return lp, types.ErrLiquidityProviderDoesNotExist
	}
	bz := store.Get(key)
	keeper.cdc.MustUnmarshal(bz, &lp)
	return lp, nil
}

func (keeper Keeper) GetLiquidityProviderIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LiquidityProviderPrefix)
}

func (keeper Keeper) GetAssetsForLiquidityProviderPaginated(ctx sdk.Context, lpAddress sdk.AccAddress,
	pagination *query.PageRequest) ([]*types.Asset, *query.PageResponse, error) {
	var assetList []*types.Asset
	store := ctx.KVStore(keeper.storeKey)
	lpStore := prefix.NewStore(store, types.LiquidityProviderPrefix)
	pageRes, err := query.FilteredPaginate(lpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var lp types.LiquidityProvider
		if len(value) <= 0 {
			return false, nil
		}
		err := keeper.cdc.Unmarshal(value, &lp)
		if err != nil {
			return false, err
		}
		if lp.Asset == nil {
			return false, nil
		}
		if lp.LiquidityProviderAddress != lpAddress.String() {
			return false, nil
		}
		if accumulate {
			assetList = append(assetList, lp.Asset)
		}
		return true, nil
	})
	if err != nil {
		return nil, &query.PageResponse{}, status.Error(codes.Internal, err.Error())
	}
	return assetList, pageRes, nil
}

func (keeper Keeper) DestroyLiquidityProvider(ctx sdk.Context, symbol string, lpAddress string) {
	key := types.GetLiquidityProviderKey(symbol, lpAddress)
	if !keeper.Exists(ctx, key) {
		return
	}
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(key)
}

func (keeper Keeper) GetLiquidityProvidersForAssetPaginated(ctx sdk.Context, asset types.Asset,
	pagination *query.PageRequest) ([]*types.LiquidityProvider, *query.PageResponse, error) {
	var lpList []*types.LiquidityProvider
	store := ctx.KVStore(keeper.storeKey)
	lpStore := prefix.NewStore(store, types.LiquidityProviderPrefix)
	pageRes, err := query.FilteredPaginate(lpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var lp types.LiquidityProvider
		if len(value) <= 0 {
			return false, nil
		}
		err := keeper.cdc.Unmarshal(value, &lp)
		if err != nil {
			return false, err
		}
		if lp.Asset == nil {
			return false, nil
		}
		if !lp.Asset.Equals(asset) {
			return false, nil
		}
		if accumulate {
			lpList = append(lpList, &lp)
		}
		return true, nil
	})
	if err != nil {
		return nil, &query.PageResponse{}, status.Error(codes.Internal, err.Error())
	}
	return lpList, pageRes, nil
}

func (keeper Keeper) GetAllLiquidityProvidersPaginated(ctx sdk.Context,
	pagination *query.PageRequest) ([]*types.LiquidityProvider, *query.PageResponse, error) {
	var lpList []*types.LiquidityProvider
	store := ctx.KVStore(keeper.storeKey)
	lpStore := prefix.NewStore(store, types.LiquidityProviderPrefix)
	pageRes, err := query.Paginate(lpStore, pagination, func(key []byte, value []byte) error {
		var liquidityProvider types.LiquidityProvider
		err := keeper.cdc.Unmarshal(value, &liquidityProvider)
		if err != nil {
			return err
		}
		lpList = append(lpList, &liquidityProvider)
		return nil
	})
	if err != nil {
		return nil, &query.PageResponse{}, status.Error(codes.Internal, err.Error())
	}
	return lpList, pageRes, nil
}
