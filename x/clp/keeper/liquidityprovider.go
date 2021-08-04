package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetLiquidityProvider(ctx sdk.Context, lp *types.LiquidityProvider) {
	if !lp.Validate() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityProviderKey(lp.Asset.Symbol, lp.LiquidityProviderAddress)
	store.Set(key, k.cdc.MustMarshalBinaryBare(lp))
}

func (k Keeper) GetLiquidityProvider(ctx sdk.Context, symbol string, lpAddress string) (types.LiquidityProvider, error) {
	var lp types.LiquidityProvider
	key := types.GetLiquidityProviderKey(symbol, lpAddress)
	store := ctx.KVStore(k.storeKey)
	if !k.Exists(ctx, key) {
		return lp, types.ErrLiquidityProviderDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &lp)
	return lp, nil
}

func (k Keeper) GetLiquidityProviderIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LiquidityProviderPrefix)
}

func (k Keeper) GetAssetsForLiquidityProviderPaginated(ctx sdk.Context, lpAddress sdk.AccAddress,
	pagination *query.PageRequest) ([]*types.Asset, *query.PageResponse, error) {

	var assetList []*types.Asset
	store := ctx.KVStore(k.storeKey)
	lpStore := prefix.NewStore(store, types.LiquidityProviderPrefix)
	pageRes, err := query.FilteredPaginate(lpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var lp types.LiquidityProvider
		if len(value) <= 0 {
			return false, nil
		}
		err := k.cdc.UnmarshalBinaryBare(value, &lp)
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

func (k Keeper) DestroyLiquidityProvider(ctx sdk.Context, symbol string, lpAddress string) {
	key := types.GetLiquidityProviderKey(symbol, lpAddress)
	if !k.Exists(ctx, key) {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

func (k Keeper) GetLiquidityProvidersForAssetPaginated(ctx sdk.Context, asset types.Asset,
	pagination *query.PageRequest) ([]*types.LiquidityProvider, *query.PageResponse, error) {

	var lpList []*types.LiquidityProvider
	store := ctx.KVStore(k.storeKey)
	lpStore := prefix.NewStore(store, types.LiquidityProviderPrefix)
	pageRes, err := query.FilteredPaginate(lpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var lp types.LiquidityProvider
		if len(value) <= 0 {
			return false, nil
		}
		err := k.cdc.UnmarshalBinaryBare(value, &lp)
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

func (k Keeper) GetAllLiquidityProvidersPaginated(ctx sdk.Context,
	pagination *query.PageRequest) ([]*types.LiquidityProvider, *query.PageResponse, error) {

	var lpList []*types.LiquidityProvider
	store := ctx.KVStore(k.storeKey)
	lpStore := prefix.NewStore(store, types.LiquidityProviderPrefix)
	pageRes, err := query.Paginate(lpStore, pagination, func(key []byte, value []byte) error {
		var liquidityProvider types.LiquidityProvider
		err := k.cdc.UnmarshalBinaryBare(value, &liquidityProvider)
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
