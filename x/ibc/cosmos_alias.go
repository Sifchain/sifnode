package ibc

import (
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/core"
	ibc_keeper "github.com/cosmos/cosmos-sdk/x/ibc/core/keeper"
)

type (
	Keeper = ibc_keeper.Keeper
)

var (
	NewCosmosAppModule = ibc.NewAppModule
)

type (
	CosmosAppModule      = ibc.AppModule
	CosmosAppModuleBasic = ibc.AppModuleBasic
)
