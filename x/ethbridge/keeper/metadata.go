package keeper

import (
	"strings"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Verifies if token name is IBC token
func IsIBCToken(name string) bool {
	parsedName := strings.Split(name, "/")
	if len(parsedName) < 1 {
		return false
	}
	if parsedName[0] != "ibc" {
		return false
	}
	return true
}

// Fetches token meteadata if it exists
func (k Keeper) GetTokenMetadata(ctx sdk.Context, denomHash string) types.TokenMetadata {
	if !k.ExistsTokenMetadata(ctx, denomHash) {
		return types.TokenMetadata{}
	}
	store := ctx.KVStore(k.storeKey)
	encodedMetadata := store.Get([]byte(denomHash))
	tokenMetadata := types.TokenMetadata{}
	k.cdc.MustUnmarshalBinaryBare(encodedMetadata, &tokenMetadata)
	return tokenMetadata
}

// Add new token metadata information
func (k Keeper) AddTokenMetadata(ctx sdk.Context, metadata types.TokenMetadata) string {
	denomHash := types.GetDenomHash(
		metadata.NetworkDescriptor,
		metadata.TokenAddress,
		metadata.Decimals,
		metadata.Name,
		metadata.Symbol,
	)
	key := []byte(denomHash)
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshalBinaryBare(&metadata)
	store.Set(key, value)
	return denomHash
}

// Deletes token metadata for IBC tokens only; returns true on success
func (k Keeper) DeleteTokenMetadata(ctx sdk.Context, denomHash string) bool {
	// Check if metadata exists first
	if !k.ExistsTokenMetadata(ctx, denomHash) {
		return false
	}
	// Check if token is IBC token or not, refuse to delete non-IBC tokens
	metadata := k.GetTokenMetadata(ctx, denomHash)
	if !IsIBCToken(metadata.Name) {
		return false
	}
	// If we made it this far, we have an IBC token, lets delete it
	key := []byte(denomHash)
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
	return true
}

// Searches the keeper to determine if a specific token has
// been stored before
func (k Keeper) ExistsTokenMetadata(ctx sdk.Context, denomHash string) bool {
	return k.Exists(ctx, []byte(denomHash))
}
