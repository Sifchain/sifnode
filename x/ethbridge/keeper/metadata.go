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
func (k Keeper) GetTokenMetadata(ctx sdk.Context, denomHash string) (types.TokenMetadata, bool) {
	if !k.ExistsTokenMetadata(ctx, denomHash) {
		return types.TokenMetadata{}, false
	}
	store := ctx.KVStore(k.storeKey)
	encodedMetadata := store.Get([]byte(denomHash))
	tokenMetadata := types.TokenMetadata{}
	k.cdc.MustUnmarshalBinaryBare(encodedMetadata, &tokenMetadata)
	return tokenMetadata, true
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

func (k Keeper) AddIBCTokenMetadata(ctx sdk.Context, metadata types.TokenMetadata, cosmosSender sdk.AccAddress) string {
	logger := k.Logger(ctx)
	if !IsIBCToken(metadata.Name) {
		logger.Error("Token is not IBC, cannot modify metadata manually")
		return ""
	}

	if !k.oracleKeeper.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return ""
	}

	denom := k.AddTokenMetadata(ctx, metadata)

	return denom
}

// Deletes token metadata for IBC tokens only; returns true on success
func (k Keeper) DeleteTokenMetadata(ctx sdk.Context, cosmosSender sdk.AccAddress, denomHash string) bool {
	logger := k.Logger(ctx)

	// Check if token is IBC token or not, refuse to delete non-IBC tokens
	metadata, success := k.GetTokenMetadata(ctx, denomHash)
	// Check if metadata exists before attempting to delete
	if !success {
		return false
	}

	if !IsIBCToken(metadata.Name) {
		return false
	}

	if !k.oracleKeeper.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
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
