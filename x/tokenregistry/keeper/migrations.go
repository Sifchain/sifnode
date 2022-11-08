package keeper

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	tkrtypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper tkrtypes.Keeper
}

func NewMigrator(keeper tkrtypes.Keeper) Migrator {
	return Migrator{keeper}
}

func (m Migrator) MigrateToVer4(ctx sdk.Context) error {
	store := ctx.KVStore(m.keeper.StoreKey())
	iterator := sdk.KVStorePrefixIterator(store, []byte{0x02})
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	m.MigrateDenomToVer4(ctx)
	return nil
}

func (m Migrator) MigrateDenomToVer4(ctx sdk.Context) {
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
