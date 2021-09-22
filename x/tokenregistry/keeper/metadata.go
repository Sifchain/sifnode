package keeper

import (
	"strings"

	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"

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
func (k keeper) GetTokenMetadata(ctx sdk.Context, denomHash string) (types.TokenMetadata, bool) {

	entry := k.GetDenom(ctx, denomHash)

	if !entry.IsWhitelisted {
		return types.TokenMetadata{}, false
	}

	metadata := types.TokenMetadata{
		Decimals:          entry.Decimals,
		Name:              entry.DisplayName,
		Symbol:            entry.DisplaySymbol,
		TokenAddress:      entry.Address,
		NetworkDescriptor: entry.Network,
	}
	return metadata, true
}

// Add new token metadata information
func (k keeper) AddTokenMetadata(ctx sdk.Context, metadata types.TokenMetadata) string {
	denomHash := ethbridgetypes.GetDenomHash(
		metadata.NetworkDescriptor,
		metadata.TokenAddress,
		metadata.Decimals,
		metadata.Name,
		metadata.Symbol,
	)

	entry := k.GetDenom(ctx, denomHash)

	if entry.IsWhitelisted {
		entry.Decimals = metadata.Decimals
		entry.DisplayName = metadata.Name
		entry.DisplaySymbol = metadata.Symbol
		entry.Address = metadata.TokenAddress
		entry.Network = metadata.NetworkDescriptor
		entry.Denom = denomHash

		k.SetToken(ctx, &entry)
	}

	return denomHash
}

func (k keeper) AddIBCTokenMetadata(ctx sdk.Context, metadata types.TokenMetadata, cosmosSender sdk.AccAddress) string {
	logger := k.Logger(ctx)
	if !IsIBCToken(metadata.Name) {
		logger.Error("Token is not IBC, cannot modify metadata manually")
		return ""
	}

	if !k.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return ""
	}

	return k.AddTokenMetadata(ctx, metadata)
}
