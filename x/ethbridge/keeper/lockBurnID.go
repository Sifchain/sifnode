package keeper

import (
	"errors"
	"fmt"
	"strconv"

	// "github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InsertNewID add a token into peggy token list
func (k Keeper) InsertNewID(ctx sdk.Context, lockBurnID string) error {
	key := []byte(lockBurnID)
	if k.Exists(ctx, key) {
		return errors.New("lock burn ID already in store")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(key, k.cdc.MustMarshalBinaryBare(false))
	return nil
}

// SetLockBurnID return if token in peggy token list
func (k Keeper) SetLockBurnID(ctx sdk.Context, lockBurnID string) (bool, error) {
	oldValue, err := k.GetLockBurnID(ctx, lockBurnID)
	if err != nil {
		return false, err
	}
	if oldValue {
		return false, nil
	}
	key := []byte(lockBurnID)
	store := ctx.KVStore(k.storeKey)
	store.Set(key, k.cdc.MustMarshalBinaryBare(true))
	return true, nil
}

// GetLockBurnID get lock burn ID from store
func (k Keeper) GetLockBurnID(ctx sdk.Context, lockBurnID string) (bool, error) {
	key := []byte(lockBurnID)
	if !k.Exists(ctx, key) {
		return false, errors.New("lock burn ID not in store")
	}

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	value := false
	k.cdc.MustUnmarshalBinaryBare(bz, &value)
	return value, nil
}

// BuildLockBurnID return lockBurnID from address and its sequence number
func BuildLockBurnID(cosmosSender fmt.Stringer, cosmosSenderSequence uint64) string {
	lockBurnID := cosmosSender.String() + strconv.FormatUint(cosmosSenderSequence, 10)
	return lockBurnID
}
