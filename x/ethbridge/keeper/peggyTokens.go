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
	tokens.Tokens = append(tokens.Tokens, token)
	store.Set(key, k.cdc.MustMarshal(&tokens))
}

// ExistsPeggyToken return if token in peggy token list
func (k Keeper) ExistsPeggyToken(ctx sdk.Context, token string) bool {
	tokens := k.GetPeggyToken(ctx)
	for _, value := range tokens.Tokens {
		if value == token {
			return true
		}
	}
	return false
}

// GetPeggyToken get peggy token list
func (k Keeper) GetPeggyToken(ctx sdk.Context) types.PeggyTokens {
	if !k.Exists(ctx, types.PeggyTokenKeyPrefix) {
		return types.PeggyTokens{}
	}
	store := ctx.KVStore(k.storeKey)
	key := types.PeggyTokenKeyPrefix
	bz := store.Get(key)
	tokens := types.PeggyTokens{}
	k.cdc.MustUnmarshal(bz, &tokens)
	return tokens
}
