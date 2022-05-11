package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracleTypes "github.com/Sifchain/sifnode/x/oracle/types"
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
		peggy2denom := types.GetDenom(migration.networkDescriptor, migration.tokenContractAddress)
		entry.Permissions = []tkrtypes.Permission{
			tkrtypes.Permission_IBCIMPORT,
		}
		entry.Peggy_2Denom = peggy2denom
		m.keeper.SetToken(ctx, entry)
		peggyTwoEntry.Denom = peggy2denom
		m.keeper.SetToken(ctx, peggyTwoEntry)
	}
}

func (k keeper) DeleteOldAdminAccount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := tkrtypes.AdminAccountStorePrefix
	store.Delete(key)
}

type DenomMigrator struct {
	denom                string
	networkDescriptor    oracleTypes.NetworkDescriptor
	tokenContractAddress types.EthereumAddress
}

func GetDenomMigrationList() []DenomMigrator {
	return []DenomMigrator{
		{
			denom:                "ceth",
			networkDescriptor:    oracleTypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM,
			tokenContractAddress: types.NewEthereumAddress("0x0000000000000000000000000000000000000000"),
		},
	}
}
