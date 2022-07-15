package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetPool(ctx sdk.Context, pool *types.Pool) error {
	if !pool.Validate() {
		return types.ErrUnableToSetPool
	}
	store := ctx.KVStore(k.storeKey)
	key, err := types.GetPoolKey(pool.ExternalAsset.Symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return err
	}
	store.Set(key, k.cdc.MustMarshal(pool))
	return nil
}

func (k Keeper) ValidatePool(pool types.Pool) bool {
	if !pool.Validate() {
		return false
	}
	_, err := types.GetPoolKey(pool.ExternalAsset.Symbol, types.GetSettlementAsset().Symbol)
	return err == nil
}
func (k Keeper) GetPool(ctx sdk.Context, symbol string) (types.Pool, error) {
	var pool types.Pool
	store := ctx.KVStore(k.storeKey)
	key, err := types.GetPoolKey(symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return pool, err
	}
	if !k.Exists(ctx, key) {
		return pool, types.ErrPoolDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, nil
}

func (k Keeper) ExistsPool(ctx sdk.Context, symbol string) bool {
	key, err := types.GetPoolKey(symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return false
	}
	return k.Exists(ctx, key)
}

// GetPools Use GetPoolsPaginated for RPC queries
func (k Keeper) GetPools(ctx sdk.Context) []*types.Pool {
	var poolList []*types.Pool
	iterator := k.GetPoolsIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &pool)
		poolList = append(poolList, &pool)
	}
	return poolList
}

func (k Keeper) GetPoolsPaginated(ctx sdk.Context, pagination *query.PageRequest) ([]*types.Pool, *query.PageResponse, error) {
	var poolList []*types.Pool
	store := ctx.KVStore(k.storeKey)
	poolStore := prefix.NewStore(store, types.PoolPrefix)
	pageRes, err := query.Paginate(poolStore, pagination, func(key []byte, value []byte) error {
		var pool types.Pool
		err := k.cdc.Unmarshal(value, &pool)
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

func (k Keeper) DestroyPool(ctx sdk.Context, symbol string) error {
	store := ctx.KVStore(k.storeKey)
	key, err := types.GetPoolKey(symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return err
	}
	if !k.Exists(ctx, key) {
		return types.ErrPoolDoesNotExist
	}
	store.Delete(key)
	return nil
}

func (k Keeper) GetPoolsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.PoolPrefix)
}

func (k Keeper) SendRowanFromPool(ctx sdk.Context, pool *types.Pool, amount sdk.Uint, recipient sdk.AccAddress) error {
	err := k.SendRowanFromPoolNoPoolUpdate(ctx, pool, amount, recipient)
	if err != nil {
		return err
	}

	return k.SetPool(ctx, pool)
}

func (k Keeper) SendRowanFromPoolNoPoolUpdate(ctx sdk.Context, pool *types.Pool, amount sdk.Uint, recipient sdk.AccAddress) error {
	if pool.NativeAssetBalance.LT(amount) {
		return fmt.Errorf("pool balance too low for transfer. Has %s but transfer wants %s", pool.NativeAssetBalance, amount)
	}

	coin := sdk.NewCoin(types.NativeSymbol, sdk.NewIntFromBigInt(amount.BigInt()))
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(amount)
	return nil
}

func (k Keeper) RemoveRowanFromPool(ctx sdk.Context, pool *types.Pool, amount sdk.Uint) error {
	if pool.NativeAssetBalance.LT(amount) {
		return fmt.Errorf("pool balance too low for transfer. Has %s but transfer wants %s", pool.NativeAssetBalance, amount)
	}

	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(amount)
	return k.SetPool(ctx, pool)
}
