package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper interface {
	IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool
	SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress)
	CheckDenomPermissions(entry *RegistryEntry, permissions []Permission) bool
	GetDenom(wl Registry, denom string) *RegistryEntry
	SetToken(ctx sdk.Context, entry *RegistryEntry)
	RemoveToken(ctx sdk.Context, denom string)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
	GetRegistry(ctx sdk.Context) Registry
}
