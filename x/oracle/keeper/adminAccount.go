package keeper

import (
	"bytes"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountPrefix
	store.Set(key, k.Cdc.MustMarshalBinaryBare(adminAccount))
}

func (k Keeper) IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool {
	account := k.GetAdminAccount(ctx)
	if account == nil {
		return false
	}
	return bytes.Equal(account, adminAccount)
}

func (k Keeper) GetAdminAccount(ctx sdk.Context) (adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountPrefix
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	k.Cdc.MustUnmarshalBinaryBare(bz, &adminAccount)
	return
}
