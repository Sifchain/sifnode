package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// This package adds set and get operations for distribution lists

func (k Keeper) SetDistributionList(ctx sdk.Context, ar types.DistributionList) error {
	if !ar.Validate() {
		return errors.Wrapf(types.ErrInvalid, "Record Details : %s", ar.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetAirdropRecordKey(ar.DistributionName)
	store.Set(key, k.cdc.MustMarshalBinaryBare(ar))
	return nil
}

func (k Keeper) GetDistributionList(ctx sdk.Context, name string) (types.DistributionList, error) {
	var ar types.DistributionList
	store := ctx.KVStore(k.storeKey)
	key := types.GetAirdropRecordKey(name)
	if !k.Exists(ctx, key) {
		return ar, errors.Wrapf(types.ErrInvalid, "Record Does not Exist : %s", ar.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &ar)
	return ar, nil
}

func (k Keeper) ExistsDistributionList(ctx sdk.Context, name string) bool {
	key := types.GetAirdropRecordKey(name)
	if k.Exists(ctx, key) {
		return true
	}
	return false
}
