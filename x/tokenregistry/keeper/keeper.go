package keeper

import (
	"bytes"
	"fmt"
	"strings"

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
	k.Logger(ctx).Info("GetRegistry", "registry", whitelist)

	for _, i := range whitelist.Entries {
		k.Logger(ctx).Info("GetRegistry", "entry", i)

	}
	return whitelist
}

func (k keeper) GetDenomPrefix(ctx sdk.Context, denom string) []byte {
	return append(types.TokenDenomPrefix, []byte(denom)...)
}

func (k keeper) GetDenom(ctx sdk.Context, denom string) types.RegistryEntry {
	result := types.RegistryEntry{}
	registry := k.GetRegistry(ctx)
	entry, _ := k.GetEntry(registry, denom)
	if entry != nil {
		result = *entry
	}

	return result
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
		if registry.DoublePeggedNetworkMap == nil {
			registry.DoublePeggedNetworkMap = make(map[uint32]bool)
		}
		registry.DoublePeggedNetworkMap[uint32(networkDescriptor)] = false
		k.SetToken(ctx, &registry)

		instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetFirstLockDoublePeg, "networkDescriptor", networkDescriptor, "registry", registry)
	}
}

// TODO get the denom temporarily, will add a map to keep the data from network+address to denom
// after confirmed we really need this to identify the denom after receive the burn event from Ethereum
func (k keeper) GetDenomFromContract(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, contract string) (string, error) {
	entries := k.GetRegistry(ctx)
	for _, entry := range entries.Entries {
		if entry.Address == contract && entry.Network == networkDescriptor {
			return entry.Denom, nil
		}
	}
	errorMsg := fmt.Sprintf("denom not found for %s in %s", contract, networkDescriptor)
	return "", errors.Wrap(errors.ErrKeyNotFound, errorMsg)
}
