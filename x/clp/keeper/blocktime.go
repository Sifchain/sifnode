package keeper

import (
	"errors"
	"fmt"
	"time"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetBlockTime(ctx sdk.Context, now time.Time) {
	store := ctx.KVStore(k.storeKey)
	bytes, err := now.MarshalBinary()
	if err != nil {
		k.Logger(ctx).Info(fmt.Sprint("Error marshalling block time: ", err))
		k.Logger(ctx).Info("Next reported block time will be off")
		return
	}

	store.Set(types.BlockTimePrefix, bytes)
}

func (k Keeper) GetBlockTime(ctx sdk.Context) (*time.Time, error) {
	t := time.Time{}
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.BlockTimePrefix)
	if bytes == nil {
		return nil, errors.New("no block time found in store")
	}

	err := t.UnmarshalBinary(bytes)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
