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
	CheckEntryPermissions(entry *RegistryEntry, permissions []Permission) bool
	GetEntry(registry Registry, denom string) (*RegistryEntry, error)
	SetToken(ctx sdk.Context, entry *RegistryEntry)
	RemoveToken(ctx sdk.Context, denom string)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
	GetRegistry(ctx sdk.Context) Registry
	SetRegistry(ctx sdk.Context, registry Registry)
	Logger(ctx sdk.Context) log.Logger
	GetTokenMetadata(ctx sdk.Context, denomHash string) (TokenMetadata, bool)
	AddTokenMetadata(ctx sdk.Context, metadata TokenMetadata) string
	AddIBCTokenMetadata(ctx sdk.Context, metadata TokenMetadata, cosmosSender sdk.AccAddress) string
	GetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor) bool
	SetFirstLockDoublePeg(ctx sdk.Context, denom string, networkDescriptor oracletypes.NetworkDescriptor)
	GetDenomFromContract(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, contract string) (string, error)
}

type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authtypes.AccountI
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}
