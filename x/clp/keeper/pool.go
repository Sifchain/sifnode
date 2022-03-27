package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (keeper Keeper) SetPool(ctx sdk.Context, pool *types.Pool) error {
	if !pool.Validate() {
		return types.ErrUnableToSetPool
	}
	store := ctx.KVStore(keeper.storeKey)
	key, err := types.GetPoolKey(pool.ExternalAsset.Symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return err
	}
	store.Set(key, keeper.cdc.MustMarshal(pool))
	return nil
}

func (keeper Keeper) ValidatePool(pool types.Pool) bool {
	if !pool.Validate() {
		return false
	}
	_, err := types.GetPoolKey(pool.ExternalAsset.Symbol, types.GetSettlementAsset().Symbol)
	return err == nil
}
func (keeper Keeper) GetPool(ctx sdk.Context, symbol string) (types.Pool, error) {
	var pool types.Pool
	store := ctx.KVStore(keeper.storeKey)
	key, err := types.GetPoolKey(symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return pool, err
	}
	if !keeper.Exists(ctx, key) {
		return pool, types.ErrPoolDoesNotExist
	}
	bz := store.Get(key)
	keeper.cdc.MustUnmarshal(bz, &pool)
	return pool, nil
}

func (keeper Keeper) ExistsPool(ctx sdk.Context, symbol string) bool {
	key, err := types.GetPoolKey(symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return false
	}
	return keeper.Exists(ctx, key)
}

func (keeper Keeper) GetPools(ctx sdk.Context) []*types.Pool {
	var poolList []*types.Pool
	iterator := keeper.GetPoolsIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		bytesValue := iterator.Value()
		keeper.cdc.MustUnmarshal(bytesValue, &pool)
		poolList = append(poolList, &pool)
	}
	return poolList
}

func (keeper Keeper) GetPoolsPaginated(ctx sdk.Context, pagination *query.PageRequest) ([]*types.Pool, *query.PageResponse, error) {
	var poolList []*types.Pool
	store := ctx.KVStore(keeper.storeKey)
	poolStore := prefix.NewStore(store, types.PoolPrefix)
	pageRes, err := query.Paginate(poolStore, pagination, func(key []byte, value []byte) error {
		var pool types.Pool
		err := keeper.cdc.Unmarshal(value, &pool)
		if err != nil {
			return err
		}
		poolList = append(poolList, &pool)
		return nil
	})
	if err != nil {
		return nil, &query.PageResponse{}, status.Error(codes.Internal, err.Error())
	}
	return poolList, pageRes, nil
}

func (keeper Keeper) DestroyPool(ctx sdk.Context, symbol string) error {
	store := ctx.KVStore(keeper.storeKey)
	key, err := types.GetPoolKey(symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return err
	}
	if !keeper.Exists(ctx, key) {
		return types.ErrPoolDoesNotExist
	}
	store.Delete(key)
	return nil
}

func (keeper Keeper) GetPoolsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.PoolPrefix)
}
