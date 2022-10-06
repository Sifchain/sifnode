package keeper

import (
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	db "github.com/tendermint/tm-db"
)

func (k Keeper) IsBlacklisted(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(append(types.BlacklistPrefix, []byte(address)...))
}

func (k Keeper) SetBlacklist(ctx sdk.Context, msg *types.MsgSetBlacklist) error {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return err
	}

	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_ETHBRIDGE, from) {
		return oracletypes.ErrNotAdminAccount
	}

	store := ctx.KVStore(k.storeKey)
	// Process removals
	var removals []string
	iter := k.getStoreIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) > 1 {
			address := string(key[1:])
			remains := false
			for _, current := range msg.Addresses {
				if current == address {
					remains = true
				}
			}

			if !remains {
				removals = append(removals, address)
			}
		}
	}
	err = iter.Close()
	if err != nil {
		return err
	}

	for _, address := range removals {
		store.Delete(append(types.BlacklistPrefix, []byte(address)...))
	}
	for _, address := range msg.Addresses {
		store.Set(append(types.BlacklistPrefix, []byte(address)...), []byte(address))
	}

	return nil
}

func (k Keeper) GetBlacklist(ctx sdk.Context) []string {
	var addresses []string
	iter := k.getStoreIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		address := string(key[1:])
		addresses = append(addresses, address)
	}
	return addresses
}

func (k Keeper) getStoreIterator(ctx sdk.Context) db.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.BlacklistPrefix)
}
