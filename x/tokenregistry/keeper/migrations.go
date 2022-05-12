package keeper

import (
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

	return nil
}
