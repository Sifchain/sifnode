package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"

	"fmt"

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

func (k Keeper) GetLiquidityProviders(ctx sdk.Context) []*types.LiquidityProvider {
	var lpList []*types.LiquidityProvider
	iterator := k.GetLiquidityProviderIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lp types.LiquidityProvider
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &lp)
		lpList = append(lpList, &lp)
	}
	return lpList
}

func (k Keeper) GetLiquidityProviderIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LiquidityProviderPrefix)
}

func (k Keeper) GetAssetsForLiquidityProvider(ctx sdk.Context, lpAddress sdk.AccAddress) []types.Asset {
	var assetList []types.Asset

	iterator := k.GetLiquidityProviderIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var lp types.LiquidityProvider
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &lp)
		if lp.LiquidityProviderAddress == lpAddress.String() {
			assetList = append(assetList, *lp.Asset) //todo: test nil panics
		}
	}

	return assetList
}

func (k Keeper) DestroyLiquidityProvider(ctx sdk.Context, symbol string, lpAddress string) {
	key := types.GetLiquidityProviderKey(symbol, lpAddress)
	if !k.Exists(ctx, key) {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

// Deprecated: GetLiquidityProvidersForAsset use GetLiquidityProvidersForAssetPaginated
func (k Keeper) GetLiquidityProvidersForAsset(ctx sdk.Context, asset types.Asset) []*types.LiquidityProvider {
	var lpList []*types.LiquidityProvider
	iterator := k.GetLiquidityProviderIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lp types.LiquidityProvider
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &lp)
		if lp.Asset.Equals(asset) {
			lpList = append(lpList, &lp)
		}
	}
	return lpList
}

func (k Keeper) GetLiquidityProvidersForAssetPaginated(ctx sdk.Context, asset types.Asset, pagination *query.PageRequest) ([]*types.LiquidityProvider, *query.PageResponse, error) {
	var lpList []*types.LiquidityProvider
	logger := k.Logger(ctx).With("func", "GetLiquidityProvidersForAssetPaginated")

	store := ctx.KVStore(k.storeKey)
	lpStore := prefix.NewStore(store, types.LiquidityProviderPrefix)

	pageRes, err := query.FilteredPaginate(lpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var lp types.LiquidityProvider
		if len(value) <= 0 {
			logger.Info(fmt.Sprintf("cannot unmarshall empty liquidity provider for key %s", key))
			return false, nil
		}

		err := k.cdc.UnmarshalBinaryBare(value, &lp)
		if err != nil {
			return false, err
		}

		if lp.Asset == nil {
			logger.Info(fmt.Sprintf("skipping liquidity provider as asset is nil for key %s", key))
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

func (k Keeper) GetAllLiquidityProviders(ctx sdk.Context) []types.LiquidityProvider {
	var lpList []types.LiquidityProvider
	iterator := k.GetLiquidityProviderIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lp types.LiquidityProvider
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &lp)
		lpList = append(lpList, lp)

	}
	return lpList
}
