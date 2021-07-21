package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

/*

type BankKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
	IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool)
}

*/

type Keeper interface {
	IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool
	IsDenomWhitelisted(ctx sdk.Context, denom string) bool
	GetDenom(ctx sdk.Context, denom string) DenomWhitelistEntry
	SetDenom(ctx sdk.Context, denom string, exp int64)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
}
