package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper interface {
	IsAdminAccount(ctx sdk.Context, moduleName AdminType, adminAccount sdk.AccAddress) bool
	SetAdminAccount(ctx sdk.Context, account *AdminAccount)
	CheckEntryPermissions(entry *RegistryEntry, permissions []Permission) bool
	GetEntry(registry Registry, denom string) (*RegistryEntry, error)
	SetToken(ctx sdk.Context, entry *RegistryEntry)
	RemoveToken(ctx sdk.Context, denom string)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
	GetRegistry(ctx sdk.Context) Registry
	SetRegistry(ctx sdk.Context, registry Registry)
}
