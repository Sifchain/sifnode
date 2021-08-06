package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This package adds set and get operations for DistributionRecord
// A distribution record is responsible for distributing funds to the recipients

func (k Keeper) SetDistributionRecord(ctx sdk.Context, dr types.DistributionRecord) error {
	if !dr.Validate() {
		return errors.Wrapf(types.ErrInvalid, "unable to set record : %s", dr.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(dr.DistributionStatus, dr.DistributionName, dr.RecipientAddress, dr.DistributionType)
	store.Set(key, k.cdc.MustMarshalBinaryBare(&dr))
	return nil
}

func (k Keeper) GetDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, status types.DistributionStatus, distributionType types.DistributionType) (*types.DistributionRecord, error) {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(status, airdropName, recipientAddress, distributionType)
	if !k.Exists(ctx, key) {
		return &dr, errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &dr)
	return &dr, nil
}

func (k Keeper) ExistsDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, status types.DistributionStatus, distributionType types.DistributionType) bool {
	key := types.GetDistributionRecordKey(status, airdropName, recipientAddress, distributionType)
	return k.Exists(ctx, key)
}

func (k Keeper) GetDistributionRecordsIterator(ctx sdk.Context, status types.DistributionStatus) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetPrefixFromStatus(status)
	if prefix != nil {
		return sdk.KVStorePrefixIterator(store, prefix)
	}
	return nil
}

func (k Keeper) DeleteDistributionRecord(ctx sdk.Context, distributionName string, recipientAddress string, status types.DistributionStatus, distributionType types.DistributionType) error {
	var dr types.DistributionRecord
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(status, distributionName, recipientAddress, distributionType)
	if !k.Exists(ctx, key) {
		return errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	store.Delete(key)
	return nil
}

func (k Keeper) GetRecordsForNameAndStatus(ctx sdk.Context, name string, status types.DistributionStatus) *types.DistributionRecords {
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
	return &res
}

func (k Keeper) GetRecordsForNameAndStatusPaginated(ctx sdk.Context, name string, s types.DistributionStatus, pagination *query.PageRequest) (*types.DistributionRecords, *query.PageResponse, error) {
	var res types.DistributionRecords
	store := ctx.KVStore(k.storeKey)
	recordsStore := prefix.NewStore(store, types.GetPrefixFromStatus(s))
	pageRes, err := query.FilteredPaginate(recordsStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var dr types.DistributionRecord
		if len(value) <= 0 {
			return false, nil
		}
		err := k.cdc.UnmarshalBinaryBare(value, &dr)
		if err != nil {
			return false, err
		}
		if dr.DistributionName != name {
			return false, nil
		}
		if accumulate {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
		return true, nil
	})
	if err != nil {
		return nil, &query.PageResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &res, pageRes, nil
}

func (k Keeper) GetRecordsForNameStatusAndType(ctx sdk.Context, name string, status types.DistributionStatus, distributionType types.DistributionType) *types.DistributionRecords {
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
		if dr.DistributionName == name && dr.DistributionType == distributionType {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
	}
	return &res
}

func (k Keeper) GetRecordsForRecipient(ctx sdk.Context, recipient string) *types.DistributionRecords {
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
		}
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_FAILED)
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
		}
	}
	return &res
}

func (k Keeper) GetLimitedRecordsForStatus(ctx sdk.Context, status types.DistributionStatus) *types.DistributionRecords {
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
		res.DistributionRecords = append(res.DistributionRecords, &dr)
		count++
		if count == types.MaxRecordsPerBlock {
			break
		}
	}
	return &res
}

func (k Keeper) GetLimitedRecordsForRunner(ctx sdk.Context,
	distributionName string,
	authorizedRunner string,
	distributionType types.DistributionType,
	status types.DistributionStatus) *types.DistributionRecords {
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
		if count == types.MaxRecordsPerBlock {
			break
		}
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == distributionName &&
			dr.DistributionStatus == types.DistributionStatus_DISTRIBUTION_STATUS_PENDING &&
			dr.AuthorizedRunner == authorizedRunner &&
			dr.DistributionType == distributionType {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
			count = count + 1
		}
	}
	return &res
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
		err := k.cdc.UnmarshalBinaryBare(bytesValue, &dr)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Unmarshal failed for record bytes : %s ", bytesValue))
			// Not panicking here .
			// Records data is not that important . We can ignore a record if it is causing an issue for chain upgrade .
			// Logging data out for investigation
			continue
		}
		res.DistributionRecords = append(res.DistributionRecords, &dr)
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		res.DistributionRecords = append(res.DistributionRecords, &dr)
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_FAILED)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		res.DistributionRecords = append(res.DistributionRecords, &dr)
	}
	return &res
}

func (k Keeper) GetRecordsForName(ctx sdk.Context, name string) *types.DistributionRecords {
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
		if dr.DistributionName == name {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}

	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == name {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
	}
	iterator = k.GetDistributionRecordsIterator(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_FAILED)
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == name {
			res.DistributionRecords = append(res.DistributionRecords, &dr)
		}
	}
	return &res
}

func (k Keeper) ChangeRecordStatus(ctx sdk.Context, dr types.DistributionRecord, height int64, newStatus types.DistributionStatus) error {
	oldStatus := dr.DistributionStatus
	dr.DistributionStatus = newStatus
	dr.DistributionCompletedHeight = height
	// Setting to completed prefix
	err := k.SetDistributionRecord(ctx, dr)
	if err != nil {
		return errors.Wrapf(types.ErrDistribution, "error setting distribution record  : %s", dr.String())
	}
	// Deleting from old prefix
	err = k.DeleteDistributionRecord(ctx, dr.DistributionName, dr.RecipientAddress, oldStatus, dr.DistributionType) // Delete the record in the pending prefix so the iteration is cheaper.
	if err != nil {
		return errors.Wrapf(types.ErrDistribution, "error deleting distribution record  : %s", dr.String())
	}
	return nil
}
