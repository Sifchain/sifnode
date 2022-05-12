package keeper

import (
	"strings"

	"github.com/Sifchain/sifnode/x/instrumentation"
	"go.uber.org/zap"

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

// AddTokenMetadata adds new token metadata information if the token does not exist in the keeper, or it does exist and IsWhitelisted is false.
func (k keeper) AddTokenMetadata(ctx sdk.Context, metadata types.TokenMetadata) string {
	denomHash := ethbridgetypes.GetDenom(
		metadata.NetworkDescriptor,
		ethbridgetypes.NewEthereumAddress(metadata.TokenAddress),
	)

	entry, err := k.GetRegistryEntry(ctx, denomHash)
	if err != nil {
		entry = &types.RegistryEntry{}
	}

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

	if !k.IsAdminAccount(ctx, cosmosSender) {
		logger.Error("cosmos sender is not admin account.")
		return ""
	}

	return k.AddTokenMetadata(ctx, metadata)
}
