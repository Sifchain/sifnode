package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

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

func (k keeper) SetAdminAccount(ctx sdk.Context, account *types.AdminAccount) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAdminAccountKey(*account)
	store.Set(key, k.cdc.MustMarshal(account))
}

func (k keeper) IsAdminAccount(ctx sdk.Context, adminType types.AdminType, adminAccount sdk.AccAddress) bool {
	accounts := k.GetAdminAccountsForType(ctx, adminType)
	if len(accounts.AdminAccounts) == 0 {
		return false
	}
	for _, account := range accounts.AdminAccounts {
		if strings.EqualFold(account.AdminAddress, adminAccount.String()) {
			return true
		}
	}
	return false
}

func (k keeper) GetAdminAccountIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AdminAccountStorePrefix)
}

func (k keeper) GetAdminAccountsForType(ctx sdk.Context, adminType types.AdminType) *types.AdminAccounts {
	var res types.AdminAccounts
	iterator := k.GetAdminAccountIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var al types.AdminAccount
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &al)
		if al.AdminType == adminType {
			res.AdminAccounts = append(res.AdminAccounts, &al)
		}
	}
	return &res
}

func (k keeper) GetAdminAccounts(ctx sdk.Context) *types.AdminAccounts {
	var res types.AdminAccounts
	iterator := k.GetAdminAccountIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var al types.AdminAccount
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &al)
		res.AdminAccounts = append(res.AdminAccounts, &al)
	}
	return &res
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
