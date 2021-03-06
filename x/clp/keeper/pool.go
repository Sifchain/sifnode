package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) error {
	if !pool.Validate() {
		return types.ErrUnableToSetPool
	}
	store := ctx.KVStore(k.storeKey)
	key, err := types.GetPoolKey(pool.ExternalAsset.Symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return err
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(pool))
	return nil
}

func (k Keeper) ValidatePool(pool types.Pool) bool {
	if !pool.Validate() {
		return false
	}
	_, err := types.GetPoolKey(pool.ExternalAsset.Symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return false
	}
	return true
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
	k.cdc.MustUnmarshalBinaryBare(bz, &pool)
	return pool, nil
}

func (k Keeper) ExistsPool(ctx sdk.Context, symbol string) bool {
	key, err := types.GetPoolKey(symbol, types.GetSettlementAsset().Symbol)
	if err != nil {
		return false
	}
	return k.Exists(ctx, key)
}

func (k Keeper) GetPools(ctx sdk.Context) types.Pools {
	var poolList types.Pools
	iterator := k.GetPoolsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &pool)
		poolList = append(poolList, pool)
	}
	return poolList
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
