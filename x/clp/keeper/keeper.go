package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the clp store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	bankKeeper types.BankKeeper
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankkeeper types.BankKeeper) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		bankKeeper: bankkeeper,
		//paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	if !pool.Validate() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolKey(pool.ExternalAsset.Ticker, pool.ExternalAsset.SourceChain)
	store.Set(key, k.cdc.MustMarshalBinaryBare(pool))
}
func (k Keeper) GetPool(ctx sdk.Context, ticker string, sourceChain string) (types.Pool, error) {
	var pool types.Pool
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolKey(ticker, sourceChain)
	if !k.Exists(ctx, key) {
		return pool, types.ErrPoolDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &pool)
	return pool, nil
}

func (k Keeper) GetPools(ctx sdk.Context) []types.Pool {
	var poolList []types.Pool
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

func (k Keeper) DestroyPool(ctx sdk.Context, ticker string, sourceChain string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolKey(ticker, sourceChain)
	if !k.Exists(ctx, key) {
		return
	}
	store.Delete(key)
}

func (k Keeper) GetPoolsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.PoolPrefix)
}

func (k Keeper) SetLiquidityProvider(ctx sdk.Context, lp types.LiquidityProvider) {
	if !lp.Validate() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityProviderKey(lp.Asset.Ticker, lp.LiquidityProviderAddress)
	store.Set(key, k.cdc.MustMarshalBinaryBare(lp))
}

func (k Keeper) GetLiquidityProvider(ctx sdk.Context, ticker string, lpAddress string) (types.LiquidityProvider, error) {
	var lp types.LiquidityProvider
	key := types.GetLiquidityProviderKey(ticker, lpAddress)
	store := ctx.KVStore(k.storeKey)
	if !k.Exists(ctx, key) {
		return lp, types.LiquidityProviderDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &lp)
	return lp, nil
}

func (k Keeper) DestroyLiquidityProvider(ctx sdk.Context, ticker string, lpAddress string) {
	key := types.GetLiquidityProviderKey(ticker, lpAddress)
	if !k.Exists(ctx, key) {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

func (k Keeper) GetLiquidityProviderIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LiquidityProviderPrefix)
}

func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
