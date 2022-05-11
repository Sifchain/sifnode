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

func (m Migrator) MigrateToVer4(ctx sdk.Context) {
	denomMigrationsList := GetDenomMigrationList()
	for _, migration := range denomMigrationsList {
		entry, err := m.keeper.GetRegistryEntry(ctx, migration.denom)
		if err != nil {
			panic(err)
		}
		peggyTwoEntry := entry
		entry.Permissions = []tkrtypes.Permission{
			tkrtypes.Permission_IBCIMPORT,
		}
		m.keeper.SetToken(ctx, entry)
		peggyTwoEntry.Denom = "sifBridge" + migration.evmChainID + migration.tokenAddress
		peggyTwoEntry.Peggy_2Denom = entry.Denom
		m.keeper.SetToken(ctx, peggyTwoEntry)
	}
}

func (k keeper) DeleteOldAdminAccount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := tkrtypes.AdminAccountStorePrefix
	store.Delete(key)
}

type DenomMigrator struct {
	denom        string
	evmChainID   string
	tokenAddress string
}

func GetDenomMigrationList() []DenomMigrator {
	return []DenomMigrator{
		{
			denom:        "ceth",
			evmChainID:   "0001",
			tokenAddress: "0x0000000000000000000000000000000000000000",
		},
	}
}
