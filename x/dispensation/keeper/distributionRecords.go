package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// This package adds set and get operations for DistributionRecord
// A distribution record is responsible for distributing funds to the recipients

func (k Keeper) SetDistributionRecord(ctx sdk.Context, dr types.DistributionRecord) error {
	if !dr.Validate() {
		return errors.Wrapf(types.ErrInvalid, "unable to set record : %s", dr.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(dr.ClaimStatus, dr.DistributionName, dr.RecipientAddress.String())
	store.Set(key, k.cdc.MustMarshalBinaryBare(dr))
	return nil
}

func (k Keeper) GetDistributionRecord(ctx sdk.Context, distributionName string, recipientAddress string, status types.ClaimStatus) (types.DistributionRecord, error) {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(status, distributionName, recipientAddress)
	if !k.Exists(ctx, key) {
		return dr, errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &dr)
	return dr, nil
}

func (k Keeper) DeleteDistributionRecord(ctx sdk.Context, distributionName string, recipientAddress string, status types.ClaimStatus) error {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(status, distributionName, recipientAddress)
	if !k.Exists(ctx, key) {
		return errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	store.Delete(key)
	return nil
}

func (k Keeper) ExistsDistributionRecord(ctx sdk.Context, distributionName string, recipientAddress string, status types.ClaimStatus) bool {
	key := types.GetDistributionRecordKey(status, distributionName, recipientAddress)
	if k.Exists(ctx, key) {
		return true
	}
	return false
}

func (k Keeper) GetDistributionRecordsIterator(ctx sdk.Context, status types.ClaimStatus) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	switch status {
	case types.Pending:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixPending)
	case types.Completed:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixCompleted)
	default:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixCompleted) // Have not decided to return completed or pending here
	}
}

func (k Keeper) GetRecordsForName(ctx sdk.Context, name string, status types.ClaimStatus) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, status)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == name {
			res = append(res, dr)
		}
	}
	return res
}

func (k Keeper) GetRecordsForRecipient(ctx sdk.Context, recipient sdk.AccAddress) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, types.Pending)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.RecipientAddress.Equals(recipient) {
			res = append(res, dr)
		}
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.Completed)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.RecipientAddress.Equals(recipient) {
			res = append(res, dr)
		}
	}
	return res
}

func (k Keeper) GetRecordsLimited(ctx sdk.Context, status types.ClaimStatus) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, status)
	count := 0
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.ClaimStatus == types.Pending {
			res = append(res, dr)
			count++
		}
		if count == types.MAX_RECORDS_PER_BLOCK {
			break
		}
	}
	return res
}
