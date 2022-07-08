package types

import (
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/tendermint/tendermint/libs/log"

	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper interface {
	StoreKey() sdk.StoreKey
	GetAdminKeeper() adminkeeper.Keeper
	CheckEntryPermissions(entry *RegistryEntry, permissions []Permission) bool
	SetToken(ctx sdk.Context, entry *RegistryEntry)
	RemoveToken(ctx sdk.Context, denom string)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
	GetRegistry(ctx sdk.Context) Registry // Deprecated DO NOT USE
	GetRegistryPaginated(ctx sdk.Context, page uint, limit uint) (Registry, error)
	GetRegistryEntry(ctx sdk.Context, denom string) (*RegistryEntry, error)
	SetRegistry(ctx sdk.Context, registry Registry)
	Logger(ctx sdk.Context) log.Logger
	GetTokenMetadata(ctx sdk.Context, denomHash string) (TokenMetadata, bool)
	AddTokenMetadata(ctx sdk.Context, metadata TokenMetadata) string
	AddIBCTokenMetadata(ctx sdk.Context, metadata TokenMetadata, cosmosSender sdk.AccAddress) string
	GetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor) bool
	SetFirstDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor)
	AddMultipleTokens(ctx sdk.Context, entries []*RegistryEntry)
	RemoveMultipleTokens(ctx sdk.Context, denoms []string)
}

type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authtypes.AccountI
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}
