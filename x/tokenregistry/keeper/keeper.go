package keeper

import (
	"bytes"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/Sifchain/sifnode/x/instrumentation"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

type keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey) types.Keeper {
	return keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// Logger returns a module-specific logger.
func (k keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k keeper) SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountStorePrefix
	store.Set(key, k.cdc.MustMarshal(&gogotypes.BytesValue{Value: adminAccount}))
}

func (k keeper) IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool {
	account := k.GetAdminAccount(ctx)
	if account == nil {
		return false
	}
	return bytes.Equal(account, adminAccount)
}

func (k keeper) GetAdminAccount(ctx sdk.Context) (adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountStorePrefix
	bz := store.Get(key)
	acc := gogotypes.BytesValue{}
	k.cdc.MustUnmarshal(bz, &acc)
	adminAccount = sdk.AccAddress(acc.Value)
	return adminAccount
}

func (k keeper) CheckEntryPermissions(entry *types.RegistryEntry, requiredPermissions []types.Permission) bool {
	for _, requiredPermission := range requiredPermissions {
		var has bool
		for _, allowedPermission := range entry.Permissions {
			if allowedPermission == requiredPermission {
				has = true
				break
			}
		}
		if !has {
			return false
		}
	}
	return true
}

func (k keeper) GetDenomPrefix(ctx sdk.Context, denom string) []byte {
	return append(types.TokenDenomPrefix, []byte(denom)...)
}

// SetToken add a new denom
func (k keeper) SetToken(ctx sdk.Context, entry *types.RegistryEntry) {
	// get a copy to avoid modify input
	tmpCopy := *entry
	tmpCopy.Sanitize()
	key := k.GetDenomPrefix(ctx, tmpCopy.Denom)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&tmpCopy)

	store.Set(key, bz)
}

// RemoveToken remove a token
func (k keeper) RemoveToken(ctx sdk.Context, denom string) {
	key := k.GetDenomPrefix(ctx, denom)
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

// GetRegistryEntry to get a token from the Token Registry without looping
// Returns Error if entry is not found, Panics if entry is found but unable to be Unmarshalled
// Index Starts at 1, returns ordered data, not based upon insertion
func (k keeper) GetRegistryEntry(ctx sdk.Context, denom string) (*types.RegistryEntry, error) {
	var entry types.RegistryEntry
	store := ctx.KVStore(k.storeKey)
	key := k.GetDenomPrefix(ctx, denom)

	bz := store.Get(key)
	if bz == nil {
		return nil, errors.Wrap(errors.ErrKeyNotFound, "registry entry not found")
	}
	k.cdc.MustUnmarshal(bz, &entry)

	return &entry, nil
}

// Iterate over the entire token registry by slicing the query over many transactions
// Limits of 100 or less only
func (k keeper) GetRegistryPaginated(ctx sdk.Context, page uint, limit uint) (types.Registry, error) {
	var entries []*types.RegistryEntry
	store := ctx.KVStore(k.storeKey)
	if limit > 100 {
		return types.Registry{}, errors.Wrap(errors.ErrTxTooLarge, "Registry Requests limited to 100 or less")
	}
	iterator := sdk.KVStorePrefixIteratorPaginated(store, types.TokenDenomPrefix, page, limit)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var registryEntry types.RegistryEntry
		key := iterator.Key()
		bz := store.Get(key)
		k.cdc.MustUnmarshal(bz, &registryEntry)

		entries = append(entries, &registryEntry)
	}

	return types.Registry{
		Entries: entries,
	}, nil
}

// DEPRECATED: Use outside of Genesis and unit/integration tests is Deprecated, DO NOT USE
// Use GetRegistryEntry instead to lookup registry entries by denom
// IF you must fetch the registry use GetRegistryPaginated
// GetRegistry get all token's metadata
func (k keeper) GetRegistry(ctx sdk.Context) types.Registry {
	var entries []*types.RegistryEntry
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.TokenDenomPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var registry types.RegistryEntry
		key := iterator.Key()
		bz := store.Get(key)
		if bz == nil {
			// If some reason an entry is nil continue rather then panic
			continue
		}
		k.cdc.MustUnmarshal(bz, &registry)

		entries = append(entries, &registry)
	}

	return types.Registry{
		Entries: entries,
	}
}

// reset all registry
// DO NOT USE AFTER PEGGY 2.0 MIGRATION OR OUTSIDE GENESIS
func (k keeper) SetRegistry(ctx sdk.Context, wl types.Registry) {
	registry := k.GetRegistry(ctx)
	for _, entry := range registry.Entries {
		k.RemoveToken(ctx, entry.Denom)
	}
	for _, item := range wl.Entries {
		k.SetToken(ctx, item)
	}
}

func (k keeper) GetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor) bool {
	registryEntry, err := k.GetRegistryEntry(ctx, denom)
	if err != nil {
		panic("Invalid Denom for Get Double Peg")
	}
	if result, ok := registryEntry.DoublePeggedNetworkMap[uint32(networkDescriptor)]; ok {
		return result
	}

	return true
}

func (k keeper) SetFirstDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor) {
	firstLockDoublePeg := k.GetFirstLockDoublePeg(ctx, denom, networkDescriptor)
	if firstLockDoublePeg {
		registryEntry, err := k.GetRegistryEntry(ctx, denom)
		if err != nil {
			panic("Invalid Denom for Set Double Peg")
		}
		if registryEntry.DoublePeggedNetworkMap == nil {
			registryEntry.DoublePeggedNetworkMap = make(map[uint32]bool)
		}
		registryEntry.DoublePeggedNetworkMap[uint32(networkDescriptor)] = false
		k.SetToken(ctx, registryEntry)

		instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetFirstDoublePeg, "networkDescriptor", networkDescriptor, "registry", registryEntry)
	}
}

func (k keeper) AddMultipleTokens(ctx sdk.Context, entries []*types.RegistryEntry) {
	for _, entry := range entries {
		k.SetToken(ctx, entry)
	}
}

func (k keeper) RemoveMultipleTokens(ctx sdk.Context, denoms []string) {
	for _, denom := range denoms {
		k.RemoveToken(ctx, denom)
	}
}
