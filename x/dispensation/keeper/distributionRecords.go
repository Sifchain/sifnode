package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) SetDistributionRecord(ctx sdk.Context, dr types.DistributionRecord) error {
	if !dr.Validate() {
		return errors.Wrapf(types.ErrInvalid, "unable to set record : %s", dr.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(dr.DistributionName, dr.Address.String())
	store.Set(key, k.cdc.MustMarshalBinaryBare(dr))
	return nil
}

func (k Keeper) GetDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string) (types.DistributionRecord, error) {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(airdropName, recipientAddress)
	if !k.Exists(ctx, key) {
		return dr, errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &dr)
	return dr, nil
}

func (k Keeper) ExistsDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string) bool {
	key := types.GetDistributionRecordKey(airdropName, recipientAddress)
	if k.Exists(ctx, key) {
		return true
	}
	return false
}
