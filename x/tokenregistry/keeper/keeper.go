package keeper

import (
	"bytes"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

type keeper struct {
	cdc codec.BinaryMarshaler

	storeKey sdk.StoreKey
	// bankKeeper    types.BankKeeper
}

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey) types.Keeper {
	return keeper{
		// bankKeeper: bankKeeper,
		cdc:      cdc,
		storeKey: storeKey,
	}
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

func (k keeper) GetDenom(ctx sdk.Context, denom string) types.RegistryEntry {
	wl := k.GetDenomWhitelist(ctx)

	for i := range wl.Entries {
		if strings.EqualFold(wl.Entries[i].Denom, strings.ToLower(denom)) &&
			wl.Entries[i] != nil {
			return *wl.Entries[i]
		}
	}

	return types.RegistryEntry{
		IsWhitelisted: false,
		Denom:         denom,
	}
}

func (k keeper) SetDenom(ctx sdk.Context, denom string, decimals int64) {
	wl := k.GetDenomWhitelist(ctx)

	denom = strings.ToLower(denom)

	var exists bool
	for i := range wl.Entries {
		if wl.Entries[i].Denom == denom {
			exists = true
			wl.Entries[i].Decimals = decimals
		}
	}

	if !exists {
		wl.Entries = append(wl.Entries, &types.RegistryEntry{
			IsWhitelisted: true,
			Denom:         denom,
			Decimals:      decimals,
		})
	}

	k.SetDenomWhitelist(ctx, wl)
}

func (k keeper) SetDenomWhitelist(ctx sdk.Context, wl types.Registry) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&wl)

	store.Set(types.WhitelistStorePrefix, bz)
}

func (k keeper) GetDenomWhitelist(ctx sdk.Context) types.Registry {
	var whitelist types.Registry
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.WhitelistStorePrefix)
	if len(bz) == 0 {
		return types.Registry{}
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &whitelist)

	return whitelist
}
