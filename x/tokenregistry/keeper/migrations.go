package keeper

import (
	tkrtypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper tkrtypes.Keeper
}

func NewMigrator(keeper tkrtypes.Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) MigrateToVer2(ctx sdk.Context) error {
	registry := m.keeper.GetRegistry(ctx)
	for _, entry := range registry.Entries {
		if entry.Decimals > 9 && m.keeper.CheckEntryPermissions(entry, []tkrtypes.Permission{tkrtypes.Permission_CLP, tkrtypes.Permission_IBCEXPORT}) {
			entry.Permissions = append(entry.Permissions, tkrtypes.Permission_IBCIMPORT)
			entry.IbcCounterpartyDenom = ""
		}
	}
	m.keeper.SetRegistry(ctx, registry)
	return nil
}

func (m Migrator) MigrateToVer3(ctx sdk.Context) error {
	accounts := tkrtypes.InitialAdminAccounts()
	for _, account := range accounts.AdminAccounts {
		m.keeper.SetAdminAccount(ctx, account)
	}
	return nil
}
