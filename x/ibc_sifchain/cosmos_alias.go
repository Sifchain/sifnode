package ibc_sifchain

import (
	ibc_transfer "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
	ibc_trasfer_keeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
)

type (
	Keeper = ibc_trasfer_keeper.Keeper
)

var (
	NewCosmosAppModule = ibc_transfer.NewAppModule
)

type (
	CosmosAppModule      = ibc_transfer.AppModule
	CosmosAppModuleBasic = ibc_transfer.AppModuleBasic
)
