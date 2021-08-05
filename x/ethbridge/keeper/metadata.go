package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// Searches the keeper to determine if a specific token has
// been stored before
func (k Keeper) ExistsTokenMetadata(ctx sdk.Context, denomHash string) bool {
	return k.Exists(ctx, []byte(denomHash))
}
