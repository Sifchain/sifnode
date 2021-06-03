package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
)

func GetCursorKey(name string) []byte {
	return append(dispensationtypes.CursorPrefix, []byte(name)...)
}

func (k Keeper) SetCursor(ctx sdk.Context, name string, position []byte) error {
	store := ctx.KVStore(k.storeKey)
	key := GetCursorKey(name)
	store.Set(key, k.cdc.MustMarshalBinaryBare(position))
	return nil
}

func (k Keeper) GetCursor(ctx sdk.Context, name string) []byte {
	var position []byte
	store := ctx.KVStore(k.storeKey)
	key := GetCursorKey(name)
	k.cdc.MustUnmarshalBinaryBare(store.Get(key), &position)
	return position
}
