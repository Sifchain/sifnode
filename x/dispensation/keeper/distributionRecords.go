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
	if !k.Exists(ctx, key) {
		return dr, errors.Wrapf(types.ErrInvalid, "record Does not exist : %s", dr.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &dr)
	return dr, nil
}

func (k Keeper) DeleteDistributionRecord(ctx sdk.Context, dr types.DistributionRecord) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionRecordKey(dr.DistributionName, dr.RecipientAddress.String(), dr.DistributionType.String())
	store.Delete(key)
}

func (k Keeper) ExistsDistributionRecord(ctx sdk.Context, airdropName string, recipientAddress string, distributionType string) bool {
	key := types.GetDistributionRecordKey(airdropName, recipientAddress, distributionType)
	if k.Exists(ctx, key) {
		return true
	}
	return false
}

func (k Keeper) GetDistributionRecordsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DistributionRecordPrefix)
}

func (k Keeper) GetRecordsForNameAll(ctx sdk.Context, name string) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
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
func (k Keeper) GetRecordsForNamePending(ctx sdk.Context, name string) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == name && dr.DistributionStatus == types.Pending {
			res = append(res, dr)
		}
	}
	return res
}

func (k Keeper) GetRecordsForNameCompleted(ctx sdk.Context, name string) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionName == name && dr.DistributionStatus == types.Completed {
			res = append(res, dr)
		}
	}
	return res
}

func (k Keeper) GetRecordsForRecipient(ctx sdk.Context, recipient sdk.AccAddress) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
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

func (k Keeper) GetPendingRecordsLimited(ctx sdk.Context, limit int) types.DistributionRecords {
	var res types.DistributionRecords
	iterator := k.GetDistributionRecordsIterator(ctx)
	count := 0
	defer iterator.Close()
	// Todo : Change the set completed from BlockBeginner to move the records to a different prefix (Or Prune it ? ).So that we can avoid extra iterations.
	// Todo : Extra iteration might be a major issue later .
	// This is performance fix and does not affect functionality. This fix has been done in the .42 version of the module
	for ; iterator.Valid(); iterator.Next() {
		var dr types.DistributionRecord
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)
		if dr.DistributionStatus == types.Pending {
			res = append(res, dr)
			count++
		}
		if count == limit {
			break
		}
	}
	return res
}
