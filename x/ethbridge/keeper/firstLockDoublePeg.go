package keeper

import (
	"strconv"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetFirstLockDoublePeg get current global nonce and update it
func (k Keeper) GetFirstLockDoublePeg(ctx sdk.Context, denom string) bool {
	prefix := k.GetFirstLockDoublePegPrefix(ctx, denom)
	store := ctx.KVStore(k.storeKey)

	if !k.ExistsFirstLockDoublePeg(ctx, prefix) {
		// global nonce start from 1
		return false
	}
	value := store.Get(prefix)
	result, err := strconv.ParseBool(string(value))
	if err != nil {
		return false
	}

	return result
}

// SetFirstLockDoublePeg set denom double pegged
func (k Keeper) SetFirstLockDoublePeg(ctx sdk.Context, denom string) {
	prefix := k.GetFirstLockDoublePegPrefix(ctx, denom)
	store := ctx.KVStore(k.storeKey)
	b := []byte{}
	b = strconv.AppendBool(b, true)
	store.Set(prefix, b)
}

// ExistsFirstLockDoublePeg check if the global nonce exists
func (k Keeper) ExistsFirstLockDoublePeg(ctx sdk.Context, prefix []byte) bool {
	if !k.Exists(ctx, prefix) {
		return false
	}
	return true
}

// GetFirstLockDoublePegPrefix compute the prefix
func (k Keeper) GetFirstLockDoublePegPrefix(ctx sdk.Context, denom string) []byte {

	return append(types.FirstLockDoublePegPrefix, []byte(denom)[:]...)
}
