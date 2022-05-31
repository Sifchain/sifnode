package keeper

import (
	"encoding/json"
	tkrtypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io/ioutil"
	"path/filepath"
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
	denomMigrations := GetDenomMigrationMap()
	for peggy1denom, peggy2denom := range denomMigrations {
		entry, err := m.keeper.GetRegistryEntry(ctx, peggy1denom)
		if err != nil {
			panic(err)
		}
		peggyTwoEntry := entry
		entry.Permissions = []tkrtypes.Permission{
			tkrtypes.Permission_IBCIMPORT,
		}
		entry.Peggy_2Denom = peggy2denom
		m.keeper.SetToken(ctx, entry)
		peggyTwoEntry.Denom = peggy2denom
		peggyTwoEntry.Peggy_2Denom = peggy2denom
		m.keeper.SetToken(ctx, peggyTwoEntry)
	}
}

func (k keeper) DeleteOldAdminAccount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := tkrtypes.AdminAccountStorePrefix
	store.Delete(key)
}

func GetDenomMigrationMap() map[string]string {
	migrationMap := map[string]string{}
	input, err := ioutil.ReadFile(DenomMigrationFilePath())
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(input, &migrationMap)
	if err != nil {
		panic(err)
	}
	return migrationMap
}

func DenomMigrationFilePath() string {
	fp, err := filepath.Abs("../../../smart-contracts/data/denom_mapping_peggy1_to_peggy2.json")
	if err != nil {
		panic(err)
	}
	return fp
}
