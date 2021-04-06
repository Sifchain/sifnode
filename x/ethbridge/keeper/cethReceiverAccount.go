package keeper

import (
	"bytes"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetCethReceiverAccount(ctx sdk.Context, cethReceiverAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.CethReceiverAccountPrefix
	// TODO Wrap in proto
	store.Set(key, k.cdc.MustMarshalBinaryBare(cethReceiverAccount))
}

func (k Keeper) IsCethReceiverAccount(ctx sdk.Context, cethReceiverAccount sdk.AccAddress) bool {
	account := k.GetCethReceiverAccount(ctx)
	return bytes.Equal(account, cethReceiverAccount)
}

func (k Keeper) IsCethReceiverAccountSet(ctx sdk.Context) bool {
	account := k.GetCethReceiverAccount(ctx)
	return account != nil
}

func (k Keeper) GetCethReceiverAccount(ctx sdk.Context) (cethReceiverAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.CethReceiverAccountPrefix
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &cethReceiverAccount)
	return
}
