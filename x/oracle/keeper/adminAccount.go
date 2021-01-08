package keeper

import (
	"bytes"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(adminAccount))
}

func (k Keeper) IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool {
	account := k.GetAdminAccount(ctx)
	return bytes.Equal(account, adminAccount)
}

func (k Keeper) GetAdminAccount(ctx sdk.Context) (adminAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.AdminAccountPrefix
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &adminAccount)
	return
}
