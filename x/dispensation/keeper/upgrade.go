package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/Sifchain/sifnode/x/dispensation/types/legacy"
)

func Upgrade086(keeper Keeper) func(ctx sdk.Context, plan upgrade.Plan) {

	// Migrates distribution records to new structure.
	return func(ctx sdk.Context, plan upgrade.Plan) {

		// Collect legacy distribution records
		iterator := keeper.GetDistributionRecordsIterator(ctx)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			var dr legacy.DistributionRecord084
			bytesValue := iterator.Value()
			keeper.cdc.MustUnmarshalBinaryBare(bytesValue, &dr)

			upgraded := types.DistributionRecord{
				DistributionStatus:          types.DistributionStatus(dr.ClaimStatus),
				DistributionName:            dr.DistributionName,
				DistributionType:            types.DistributionTypeUnknown,
				RecipientAddress:            dr.RecipientAddress,
				Coins:                       dr.Coins,
				DistributionStartHeight:     dr.DistributionStartHeight,
				DistributionCompletedHeight: dr.DistributionCompletedHeight,
			}

			key := iterator.Key()
			store := ctx.KVStore(keeper.storeKey)
			store.Set(key, keeper.cdc.MustMarshalBinaryBare(upgraded))
		}
	}
}
