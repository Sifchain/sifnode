package keeper

import (
	"bytes"
	"fmt"

	"github.com/Sifchain/sifnode/x/instrumentation"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/tendermint/libs/log"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

type keeper struct {
	cdc codec.BinaryMarshaler

	storeKey sdk.StoreKey
}

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey) types.Keeper {
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
	store.Set(key, k.cdc.MustMarshalBinaryBare(&gogotypes.BytesValue{Value: adminAccount}))
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
	k.cdc.MustUnmarshalBinaryBare(bz, &acc)

	adminAccount = sdk.AccAddress(acc.Value)

	return adminAccount
}

func (k keeper) IsDenomWhitelisted(ctx sdk.Context, denom string) bool {
	d := k.GetDenom(ctx, denom)

	return d.IsWhitelisted
}

func (k keeper) CheckDenomPermissions(ctx sdk.Context, denom string, requiredPermissions []types.Permission) bool {
	d := k.GetDenom(ctx, denom)

	for _, requiredPermission := range requiredPermissions {
		var has bool
		for _, allowedPermission := range d.Permissions {
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

func (k keeper) GetDenom(ctx sdk.Context, denom string) types.RegistryEntry {

	var entry types.RegistryEntry
	key := k.GetDenomPrefix(ctx, denom)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	k.cdc.MustUnmarshalBinaryBare(bz, &entry)

	return entry
}

// SetToken add a new denom
func (k keeper) SetToken(ctx sdk.Context, entry *types.RegistryEntry) {
	entry.Sanitize()
	key := k.GetDenomPrefix(ctx, entry.Denom)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(entry)

	store.Set(key, bz)
}

// RemoveToken remove a token
func (k keeper) RemoveToken(ctx sdk.Context, denom string) {
	key := k.GetDenomPrefix(ctx, denom)
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
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
		k.cdc.MustUnmarshalBinaryBare(bz, &registry)

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

// Exists chec if the key existed in db.
func (k keeper) Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
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

		instrumentation.PeggyCheckpoint(ctx.Logger(), "SetFirstLockDoublePeg", "networkDescriptor", networkDescriptor, "registry", registry)
	}
}
