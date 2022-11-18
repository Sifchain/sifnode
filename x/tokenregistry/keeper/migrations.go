package keeper

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	tkrtypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	tkrtypes.Keeper
}

func NewMigrator(keeper tkrtypes.Keeper) Migrator {
	return Migrator{keeper}
}

func (m Migrator) MigrateToVer4(ctx sdk.Context) error {
	store := ctx.KVStore(m.StoreKey())
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

	denomMigrationMap := getDenomMigrationMap()
	for peggy1denom, peggy2denom := range denomMigrationMap {
		m.Keeper.SetPeggy2Denom(ctx, peggy1denom, peggy2denom)
	}

	return nil
}

func getDenomMigrationMap() map[string]string {
	migrationMap := map[string]string{}
	input, err := ioutil.ReadFile(denomMigrationFilePath())
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(input, &migrationMap)
	if err != nil {
		panic(err)
	}
	return migrationMap
}

func denomMigrationFilePath() string {
	fp, err := filepath.Abs("../../../smart-contracts/data/denom_mapping_peggy1_to_peggy2.json")
	if err != nil {
		panic(err)
	}
	return fp
}
