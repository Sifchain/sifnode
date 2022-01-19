package keeper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Sifchain/sifnode/x/instrumentation"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/tendermint/libs/log"
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

func (k keeper) GetEntry(wl types.Registry, denom string) (*types.RegistryEntry, error) {
	for i := range wl.Entries {
		e := wl.Entries[i]
		if e != nil && strings.EqualFold(e.Denom, denom) {
			return wl.Entries[i], nil
		}
	}
	return nil, errors.Wrap(errors.ErrKeyNotFound, "registry entry not found")
}

func (k keeper) SetToken(ctx sdk.Context, entry *types.RegistryEntry) {
	wl := k.GetRegistry(ctx)
	for i := range wl.Entries {
		if wl.Entries[i] != nil && strings.EqualFold(wl.Entries[i].Denom, entry.Denom) {
			wl.Entries[i] = entry
			k.SetRegistry(ctx, wl)
			return
		}
	}
	wl.Entries = append(wl.Entries, entry)
	k.SetRegistry(ctx, wl)
}

func (k keeper) RemoveToken(ctx sdk.Context, denom string) {
	registry := k.GetRegistry(ctx)
	updated := make([]*types.RegistryEntry, 0)
	for _, t := range registry.Entries {
		if t != nil && !strings.EqualFold(t.Denom, denom) {
			updated = append(updated, t)
		}
	}
	k.SetRegistry(ctx, types.Registry{
		Entries: updated,
	})
}

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
		k.cdc.MustUnmarshal(bz, &registry)

		entries = append(entries, &registry)
	}

	return types.Registry{
		Entries: entries,
	}
}

// SetRegistry add a bunch of tokens
func (k keeper) SetRegistry(ctx sdk.Context, wl types.Registry) {

	for _, item := range wl.Entries {
		k.SetToken(ctx, item)
	}
}

func (k keeper) GetDenomPrefix(ctx sdk.Context, denom string) []byte {
	return append(types.TokenDenomPrefix, []byte(denom)...)
}

func (k keeper) GetDenom(ctx sdk.Context, denom string) types.RegistryEntry {

	var entry types.RegistryEntry
	key := k.GetDenomPrefix(ctx, denom)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	k.cdc.MustUnmarshal(bz, &entry)

	return entry
}

func (k keeper) GetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor) bool {
	registry := k.GetDenom(ctx, denom)
	if result, ok := registry.DoublePeggedNetworkMap[uint32(networkDescriptor)]; ok {
		return result
	}

	return true
}

func (k keeper) SetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor) {
	firstLockDoublePeg := k.GetFirstLockDoublePeg(ctx, denom, networkDescriptor)
	if firstLockDoublePeg {
		registry := k.GetDenom(ctx, denom)
		registry.DoublePeggedNetworkMap[uint32(networkDescriptor)] = false
		k.SetToken(ctx, &registry)

		instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetFirstLockDoublePeg, "networkDescriptor", networkDescriptor, "registry", registry)
	}
}
