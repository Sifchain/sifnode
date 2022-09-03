package keeper

import (
	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

type keeper struct {
	cdc         codec.BinaryCodec
	storeKey    sdk.StoreKey
	adminKeeper adminkeeper.Keeper
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, adminKeeper adminkeeper.Keeper) types.Keeper {
	return keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		adminKeeper: adminKeeper,
	}
}

func (k keeper) StoreKey() sdk.StoreKey {
	return k.storeKey
}

func (k keeper) GetAdminKeeper() adminkeeper.Keeper {
	return k.adminKeeper
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
		if e != nil && e.Denom == denom {
			return wl.Entries[i], nil
		}
	}
	return nil, errors.Wrap(errors.ErrKeyNotFound, "registry entry not found")
}

func (k keeper) SetToken(ctx sdk.Context, entry *types.RegistryEntry) {
	wl := k.GetRegistry(ctx)
	for i := range wl.Entries {
		if wl.Entries[i] != nil && types.StringCompare(wl.Entries[i].Denom, entry.Denom) {
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
		if t != nil && !types.StringCompare(t.Denom, denom) {
			updated = append(updated, t)
		}
	}
	k.SetRegistry(ctx, types.Registry{
		Entries: updated,
	})
}

func (k keeper) SetRegistry(ctx sdk.Context, wl types.Registry) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&wl)
	store.Set(types.WhitelistStorePrefix, bz)
}

func (k keeper) GetRegistry(ctx sdk.Context) types.Registry {
	var whitelist types.Registry
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.WhitelistStorePrefix)
	if len(bz) == 0 {
		return types.Registry{}
	}
	k.cdc.MustUnmarshal(bz, &whitelist)
	return whitelist
}
