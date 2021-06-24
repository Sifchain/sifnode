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
<<<<<<< HEAD
	key := types.GetDistributionRecordKey(dr.DistributionStatus, dr.DistributionName, dr.RecipientAddress)
	store.Set(key, k.cdc.MustMarshalBinaryBare(&dr))
	return nil
}

func (k Keeper) GetDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, status types.DistributionStatus) (*types.DistributionRecord, error) {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(status, airdropName, recipientAddress)
=======
	key := types.GetDistributionRecordKey(dr.DistributionName, dr.RecipientAddress.String(), dr.DistributionType.String())
	store.Set(key, k.cdc.MustMarshalBinaryBare(dr))
	return nil
}

func (k Keeper) SetDistributionRecordFailed(ctx sdk.Context, dr types.DistributionRecord) error {
	if !dr.Validate() {
		return errors.Wrapf(types.ErrInvalid, "unable to set record : %s", dr.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordFailedKey(dr.DistributionName, dr.RecipientAddress.String(), dr.DistributionType.String())
	store.Set(key, k.cdc.MustMarshalBinaryBare(dr))
	return nil
}

func (k Keeper) MoveRecordToFailed(ctx sdk.Context, dr types.DistributionRecord) error {
	err := k.SetDistributionRecordFailed(ctx, dr)
	if err != nil {
		return err
	}
	k.DeleteDistributionRecord(ctx, dr)
	return nil
}

func (k Keeper) GetDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, distributionType string) (types.DistributionRecord, error) {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(airdropName, recipientAddress, distributionType)
>>>>>>> develop
	if !k.Exists(ctx, key) {
		return &dr, errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &dr)
	return &dr, nil
}

<<<<<<< HEAD
func (k Keeper) ExistsDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, status types.DistributionStatus) bool {
	key := types.GetDistributionRecordKey(status, airdropName, recipientAddress)
	return k.Exists(ctx, key)
}

func (k Keeper) GetDistributionRecordsIterator(ctx sdk.Context, status types.DistributionStatus) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	switch status {
	case types.DistributionStatus_DISTRIBUTION_STATUS_PENDING:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixPending)
	case types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED:
		return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixCompleted)
	default:
		return nil
=======
func (k Keeper) DeleteDistributionRecord(ctx sdk.Context, dr types.DistributionRecord) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(dr.DistributionName, dr.RecipientAddress.String(), dr.DistributionType.String())
	store.Delete(key)
}

func (k Keeper) ExistsDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, distributionType string) bool {
	key := types.GetDistributionRecordKey(airdropName, recipientAddress, distributionType)
	if k.Exists(ctx, key) {
		return true
>>>>>>> develop
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

<<<<<<< HEAD
func (k Keeper) GetRecordsForName(ctx sdk.Context, name string, status types.DistributionStatus) types.DistributionRecords {
=======
func (k Keeper) GetDistributionRecordsIteratorFailed(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefixFailed)
}

func (k Keeper) GetRecordsForNameAll(ctx sdk.Context, name string) types.DistributionRecords {
>>>>>>> develop
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, status)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic("Failed to close iterator")
		}
	}(iterator)
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

<<<<<<< HEAD
func (k Keeper) GetRecordsForRecipient(ctx sdk.Context, recipient string) types.DistributionRecords {
=======
func (k Keeper) GetRecordsForNameAllFailed(ctx sdk.Context, name string) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIteratorFailed(ctx)
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

func (k Keeper) GetRecordsForNameAndType(ctx sdk.Context, name string, drType types.DistributionType) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == name && dr.DistributionType == drType {
			res = append(res, dr)
		}
	}
	return res
}

// The two queries have been replaced with a single query with status as a field in the .42 version
func (k Keeper) GetRecordsForNamePending(ctx sdk.Context, distributionName string) types.DistributionRecords {
>>>>>>> develop
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic("Failed to close iterator")
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
<<<<<<< HEAD
		if dr.RecipientAddress == recipient {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic("Failed to close iterator")
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.RecipientAddress == recipient {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
=======
		if dr.DistributionName == distributionName && dr.DistributionStatus == types.Pending {
			res = append(res, dr)
>>>>>>> develop
		}
	}
	return res
}

<<<<<<< HEAD
func (k Keeper) GetRecordsLimited(ctx sdk.Context, status types.DistributionStatus) types.DistributionRecords {
=======
func (k Keeper) GetRecordsForNameCompleted(ctx sdk.Context, distributionName string) types.DistributionRecords {
>>>>>>> develop
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, status)
	count := 0
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic("Failed to close iterator")
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
<<<<<<< HEAD
		res.DistributionRecords = append(res.DistributionRecords, &dr)
		count++
		if count == types.MaxRecordsPerBlock {
			break
=======
		if dr.DistributionName == distributionName && dr.DistributionStatus == types.Completed {
			res = append(res, dr)
>>>>>>> develop
		}
	}
	return res
}

func (k Keeper) GetRecords(ctx sdk.Context) *types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic("Failed to close iterator")
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		res.DistributionRecords = append(res.DistributionRecords, &dr)
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		res.DistributionRecords = append(res.DistributionRecords, &dr)
	}
	return &res
}

func (k Keeper) GetRecordsForNamePendingLimited(ctx sdk.Context, distributionName string, limit int, runner sdk.AccAddress, distributionType types.DistributionType) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	count := 0
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if count == limit {
			break
		}
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == distributionName &&
			dr.DistributionStatus == types.Pending &&
			dr.AuthorizedRunner.Equals(runner) &&
			dr.DistributionType == distributionType {
			res = append(res, dr)
			count = count + 1
		}
	}
	return res
}

func (k Keeper) GetRecords(ctx sdk.Context) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	defer iterator.Close()
<<<<<<< HEAD
=======
	// Todo : Change the set completed from BlockBeginner to move the records to a different prefix (Or Prune it ? ).So that we can avoid extra iterations.
	// Todo : Extra iteration might be a major issue later .
	// This is performance fix and does not affect functionality. This fix has been done in the .42 version of the module
>>>>>>> develop
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		res = append(res, dr)
	}
	return res
}

func (k Keeper) GetRecordsForNamePendingLimited(ctx sdk.Context, distributionName string, limit int, runner sdk.AccAddress, distributionType types.DistributionType) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	count := 0
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if count == limit {
			break
		}
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == distributionName &&
			dr.DistributionStatus == types.Pending &&
			dr.AuthorizedRunner.Equals(runner) &&
			dr.DistributionType == distributionType {
			res = append(res, dr)
			count = count + 1
		}
	}
	return res
}

func (k Keeper) GetRecords(ctx sdk.Context) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		res = append(res, dr)
	}
	return res
}
