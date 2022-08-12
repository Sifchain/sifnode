//go:build !FEATURE_TOGGLE_SDK_045
// +build !FEATURE_TOGGLE_SDK_045

package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	db "github.com/tendermint/tm-db"
)

func (k Keeper) getStoreIterator(ctx sdk.Context) db.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.BlacklistPrefix, nil)
}
