package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddPeggyToken add a token into peggy token list
func (k Keeper) AddPeggyToken(ctx sdk.Context, token string) {
	if k.ExistsPeggyToken(ctx, token) {
		return
	}
	tokens := k.GetPeggyToken(ctx)

	store := ctx.KVStore(k.storeKey)
	key := types.PeggyTokenKeyPrefix

	tokens = append(tokens, token)
	store.Set(key, k.cdc.MustMarshalBinaryBare(tokens))
}

// ExistsPeggyToken return if token in peggy token list
func (k Keeper) ExistsPeggyToken(ctx sdk.Context, token string) bool {
	tokens := k.GetPeggyToken(ctx)
	for _, value := range tokens {
		if value == token {
			return true
		}
	}
	return false
}

// GetPeggyToken get peggy token list
func (k Keeper) GetPeggyToken(ctx sdk.Context) (tokens []string) {
	if !k.Exists(ctx, types.PeggyTokenKeyPrefix) {
		return make([]string, 0)
	}
	store := ctx.KVStore(k.storeKey)
	key := types.PeggyTokenKeyPrefix
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &tokens)
	return
}
