package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (k Keeper) SetAirdropRecord(ctx sdk.Context, ar types.AirdropRecord) error {
	if !ar.Validate() {
		return errors.Wrapf(types.ErrInvalid, "Record Details : %s", ar.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetAirdropRecordKey(ar.AirdropName)
	store.Set(key, k.cdc.MustMarshalBinaryBare(ar))
	return nil
}

func (k Keeper) GetAirdropRecord(ctx sdk.Context, airdropName string) (types.AirdropRecord, error) {
	var ar types.AirdropRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetAirdropRecordKey(airdropName)
	if !k.Exists(ctx, key) {
		return ar, errors.Wrapf(types.ErrInvalid, "Record Does not Exist : %s", ar.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &ar)
	return ar, nil
}

func (k Keeper) ExistsAirdrop(ctx sdk.Context, airdropName string) bool {
	key := types.GetAirdropRecordKey(airdropName)
	if k.Exists(ctx, key) {
		return true
	}
	return false
}
