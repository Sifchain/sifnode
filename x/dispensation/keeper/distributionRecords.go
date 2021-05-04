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
	key := types.GetDistributionRecordKey(dr.DistributionStatus, dr.DistributionName, dr.RecipientAddress)
	store.Set(key, k.cdc.MustMarshalBinaryBare(&dr))
	return nil
}

func (k Keeper) GetDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, status types.DistributionStatus) (*types.DistributionRecord, error) {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(status, airdropName, recipientAddress)
	if !k.Exists(ctx, key) {
		return &dr, errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &dr)
	return &dr, nil
}

func (k Keeper) ExistsDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, status types.DistributionStatus) bool {
	key := types.GetDistributionRecordKey(status, airdropName, recipientAddress)
	if k.Exists(ctx, key) {
		return true
	}
	return false
}

func (k Keeper) GetDistributionRecordsIterator(ctx sdk.Context, status types.DistributionStatus) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	switch status {
	case types.DistributionStatus_DISTRIBUTION_STATUS_PENDING:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixPending)
	case types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixCompleted)
	default:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixCompleted)
	}
}

func (k Keeper) DeleteDistributionRecord(ctx sdk.Context, distributionName string, recipientAddress string, status types.DistributionStatus) error {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(status, distributionName, recipientAddress)
	if !k.Exists(ctx, key) {
		return errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	store.Delete(key)
	return nil
}

func (k Keeper) GetRecordsForName(ctx sdk.Context, name string, status types.DistributionStatus) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, status)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == name {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
	}
	return res
}

func (k Keeper) GetRecordsForRecipient(ctx sdk.Context, recipient string) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.RecipientAddress == recipient {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.RecipientAddress == recipient {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
	}
	return res
}

func (k Keeper) GetRecordsLimited(ctx sdk.Context, status types.DistributionStatus) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, status)
	count := 0
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionStatus == types.DistributionStatus_DISTRIBUTION_STATUS_PENDING {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
			count++
		}
		if count == types.MaxRecordsPerBlock {
			break
		}
	}
	return res
}
