package keeper

import (
	"strings"

	"github.com/Sifchain/sifnode/x/instrumentation"
	"go.uber.org/zap"

	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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

// Fetches token metadata if it exists
func (k keeper) GetTokenMetadata(ctx sdk.Context, denomHash string) (types.TokenMetadata, bool) {

	entry, err := k.GetRegistryEntry(ctx, denomHash)
	if errors.IsOf(err, errors.ErrKeyNotFound) {
		return types.TokenMetadata{}, false
	}
	if err != nil {
		panic("Unahandled Registry Error")
	}

	// This is commented out because it is superceded by whats in develop, this change makes testing easier
	// if !entry.IsWhitelisted {
	// 	ctx.Logger().Debug(instrumentation.PeggyTestMarker, "It is not whitelisted", zap.Reflect("entry", entry))
	// 	instrumentation.PeggyCheckpoint(ctx.Logger(), instrument)

	// 	return types.TokenMetadata{}, false
	// }

	metadata := types.TokenMetadata{
		Decimals:          entry.Decimals,
		Name:              entry.DisplayName,
		Symbol:            entry.DisplaySymbol,
		TokenAddress:      entry.Address,
		NetworkDescriptor: entry.Network,
	}

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.GetTokenMetadata, "denomHash", denomHash, "entry", zap.Reflect("entry", entry), "metadata", zap.Reflect("metadata", metadata))

	return metadata, true
}

// AddTokenMetadata adds new token metadata information if the token does not exist in the keeper.
// If it already exists, it just returns the denomHash.
func (k keeper) AddTokenMetadata(ctx sdk.Context, metadata types.TokenMetadata) string {
	denomHash := ethbridgetypes.GetDenom(
		metadata.NetworkDescriptor,
		ethbridgetypes.NewEthereumAddress(metadata.TokenAddress),
	)

	// Verify the Registry Entry is empty before adding token metadata
	// If it is not, simply return the current denomHash without updating
	// If any other error is returned, panic.
	entry, err := k.GetRegistryEntry(ctx, denomHash)
	// If entry was found since no error was returned
	if err == nil {
		return denomHash
		// If Error was reported, verify its only of type Key Not Found, otherwise panic
	} else if !errors.IsOf(err, errors.ErrKeyNotFound) {
		panic("Unexpected error from GetRegistryEntry")
	}

	entry = &types.RegistryEntry{}
	entry.Decimals = metadata.Decimals
	entry.DisplayName = metadata.Name
	entry.DisplaySymbol = metadata.Symbol
	entry.Address = metadata.TokenAddress
	entry.Network = metadata.NetworkDescriptor
	entry.Denom = denomHash

	k.SetToken(ctx, entry)

	instrumentation.PeggyCheckpoint(k.Logger(ctx), instrumentation.AddTokenMetadata, "entry", entry)

	return denomHash
}

func (k keeper) AddIBCTokenMetadata(ctx sdk.Context, metadata types.TokenMetadata, cosmosSender sdk.AccAddress) string {
	logger := k.Logger(ctx)
	if !IsIBCToken(metadata.Name) {
		logger.Error("Token is not IBC, cannot modify metadata manually")
		return ""
	}

	if !k.GetAdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_TOKENREGISTRY, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return ""
	}

	return k.AddTokenMetadata(ctx, metadata)
}
