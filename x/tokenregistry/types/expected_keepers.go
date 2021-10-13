package types

import (
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper interface {
	IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool
	SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress)
	IsDenomWhitelisted(ctx sdk.Context, denom string) bool
	CheckDenomPermissions(ctx sdk.Context, denom string, permissions []Permission) bool
	GetDenom(ctx sdk.Context, denom string) RegistryEntry
	SetToken(ctx sdk.Context, entry *RegistryEntry)
	RemoveToken(ctx sdk.Context, denom string)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
	GetDenomWhitelist(ctx sdk.Context) Registry
	Logger(ctx sdk.Context) log.Logger
	GetTokenMetadata(ctx sdk.Context, denomHash string) (TokenMetadata, bool)
	AddTokenMetadata(ctx sdk.Context, metadata TokenMetadata) string
	AddIBCTokenMetadata(ctx sdk.Context, metadata TokenMetadata, cosmosSender sdk.AccAddress) string
	GetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor) bool
	SetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor)
}

type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authtypes.AccountI
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}
