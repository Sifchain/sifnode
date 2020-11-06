package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the clp store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	BankKeeper types.BankKeeper
	paramstore params.Subspace
}

// NewKeeper creates a clp keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankkeeper types.BankKeeper, paramstore params.Subspace) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		BankKeeper: bankkeeper,
		paramstore: paramstore.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) error {
	if !pool.Validate() {
		return types.ErrUnableToSetPool
	}
	store := ctx.KVStore(k.storeKey)
	key, err := types.GetPoolKey(pool.ExternalAsset.Ticker, types.GetSettlementAsset().Ticker)
	if err != nil {
		return err
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(pool))
	return nil
}
func (k Keeper) GetPool(ctx sdk.Context, ticker string) (types.Pool, error) {
	var pool types.Pool
	store := ctx.KVStore(k.storeKey)
	key, err := types.GetPoolKey(ticker, types.GetSettlementAsset().Ticker)
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

func (k Keeper) ExistsPool(ctx sdk.Context, ticker string) bool {
	key, err := types.GetPoolKey(ticker, types.GetSettlementAsset().Ticker)
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

func (k Keeper) DestroyPool(ctx sdk.Context, ticker string) error {
	store := ctx.KVStore(k.storeKey)
	key, err := types.GetPoolKey(ticker, types.GetSettlementAsset().Ticker)
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

func (k Keeper) SetLiquidityProvider(ctx sdk.Context, lp types.LiquidityProvider) {
	if !lp.Validate() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityProviderKey(lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	store.Set(key, k.cdc.MustMarshalBinaryBare(lp))
}

func (k Keeper) GetLiquidityProvider(ctx sdk.Context, ticker string, lpAddress string) (types.LiquidityProvider, error) {
	var lp types.LiquidityProvider
	key := types.GetLiquidityProviderKey(ticker, lpAddress)
	store := ctx.KVStore(k.storeKey)
	if !k.Exists(ctx, key) {
		return lp, types.ErrLiquidityProviderDoesNotExist
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

func (k Keeper) GetLiqudityProvidersForAsset(ctx sdk.Context, asset types.Asset) []types.LiquidityProvider {
	var lpList []types.LiquidityProvider
	iterator := k.GetLiquidityProviderIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lp types.LiquidityProvider
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &lp)
		if lp.Asset == asset {
			lpList = append(lpList, lp)
		}
	}
	return lpList
}

func (k Keeper) GetLiquidityProviderIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LiquidityProviderPrefix)
}

func (k Keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}
