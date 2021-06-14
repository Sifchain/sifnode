package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/Sifchain/sifnode/x/dispensation/types/legacy"
)

func MigrateRecords(keeper Keeper) func(ctx sdk.Context, plan upgrade.Plan) {
	// Migrates distribution records, and distributions to new structure.
	return func(ctx sdk.Context, plan upgrade.Plan) {
		UpgradeDistributionRecords(ctx, keeper)
		UpgradeDistributions(ctx, keeper)
	}
}

func UpgradeDistributions(ctx sdk.Context, keeper Keeper) {
	var keysForDeletion []string
	keysForSetting := make(map[string][]byte, 10)

	// Collect legacy distribution records
	iterator := keeper.GetDistributionIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr legacy.Distribution084
		bytesValue := iterator.Value()
		keeper.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)

		upgraded := types.Distribution{
			DistributionName: dr.DistributionName,
			DistributionType: dr.DistributionType,
		}

		key := types.GetDistributionsKey(upgraded.DistributionName, upgraded.DistributionType)
		keysForDeletion = append(keysForDeletion, string(iterator.Key()))
		keysForSetting[string(key)] = keeper.cdc.MustMarshalBinaryBare(upgraded)
	}

	store := ctx.KVStore(keeper.storeKey)
	// Delete old before setting new, in case of key clash.
	for _, key := range keysForDeletion {
		store.Delete([]byte(key))
	}
	for key, value := range keysForSetting {
		store.Set([]byte(key), value)
	}
}

func UpgradeDistributionRecords(ctx sdk.Context, keeper Keeper) {
	var drKeysForDeletion []string
	drKeysForSetting := make(map[string][]byte, 10)

	// Collect legacy distribution records
	iterator := keeper.GetDistributionRecordsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dr legacy.DistributionRecord084
		bytesValue := iterator.Value()
		keeper.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)

		upgraded := types.DistributionRecord{
			DistributionStatus: types.DistributionStatus(dr.ClaimStatus),
			DistributionName:   dr.DistributionName,
			// All distributions so far have been Airdrops.
			DistributionType:            types.Airdrop,
			RecipientAddress:            dr.RecipientAddress,
			Coins:                       dr.Coins,
			DistributionStartHeight:     dr.DistributionStartHeight,
			DistributionCompletedHeight: dr.DistributionCompletedHeight,
		}

		key := types.GetDistributionRecordKey(upgraded.DistributionName, upgraded.RecipientAddress.String(), upgraded.DistributionType.String())
		drKeysForDeletion = append(drKeysForDeletion, string(iterator.Key()))
		drKeysForSetting[string(key)] = keeper.cdc.MustMarshalBinaryBare(upgraded)
	}

	store := ctx.KVStore(keeper.storeKey)
	// Delete old before setting new, in case of key clash.
	for _, key := range drKeysForDeletion {
		store.Delete([]byte(key))
	}
	for key, value := range drKeysForSetting {
		store.Set([]byte(key), value)
	}
}
