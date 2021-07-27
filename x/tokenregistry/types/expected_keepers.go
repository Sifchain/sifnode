package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper interface {
	IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool
	IsDenomWhitelisted(ctx sdk.Context, denom string) bool
	GetDenom(ctx sdk.Context, denom string) RegistryEntry
	SetToken(ctx sdk.Context, entry *RegistryEntry)
	RemoveToken(ctx sdk.Context, denom string)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
	GetDenomWhitelist(ctx sdk.Context) Registry
}
