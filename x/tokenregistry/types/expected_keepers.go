package types

import (
	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
}
