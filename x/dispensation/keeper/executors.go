package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pkg/errors"
)

//CreateDrops creates new drop Records . These records are then used to facilitate distribution
// Each Recipient and DropName generate a unique Record
func (k Keeper) CreateDrops(ctx sdk.Context, output []bank.Output, name string, distributionType types.DistributionType) error {
	for _, receiver := range output {
		distributionRecord := types.NewDistributionRecord(name, distributionType, receiver.Address, receiver.Coins, ctx.BlockHeight(), -1)
		if k.ExistsDistributionRecord(ctx, name, receiver.Address.String(), distributionType.String()) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address.String(), distributionType.String())
			if err != nil {
				return errors.Wrapf(types.ErrDistribution, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			if oldRecord.DistributionStatus == types.Pending {
				distributionRecord = distributionRecord.Add(oldRecord)
			}

		}
		distributionRecord.DistributionStatus = types.Pending
		err := k.SetDistributionRecord(ctx, distributionRecord)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", distributionRecord.String())
		}
	}
	return nil
}

// DistributeDrops is called at the beginning of every block .
// It checks if any pending records are present , if there are it completes the top 10
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64, DistributionName string) error {
	// TODO replace 10 with a variable declared in genesis or a constant in keys.go.
	pendingRecords := k.GetRecordsForNamePendingLimited(ctx, DistributionName, 10)
	for _, record := range pendingRecords {
		err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, record.RecipientAddress, record.Coins)
		if err != nil {
			err := errors.Wrapf(err, "Distribution failed for address  : %s", record.RecipientAddress.String())
			ctx.Logger().Error(err.Error())
			err = k.MoveRecordToFailed(ctx, record)
			if err != nil {
				return errors.Wrapf(err, "Unable to set Distribution Records to Failed : %s", record.String())
			}
			continue
		}
		record.DistributionStatus = types.Completed
		record.DistributionCompletedHeight = height
		err = k.SetDistributionRecord(ctx, record)
		if err != nil {
			err := errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", record.String())
			ctx.Logger().Error(err.Error())
			// If the SetDistributionRecord returns error , that would mean the required amount was transferred to the user , but the record was not set to completed .
			// In this case we try to take the funds back from the user , and attempt the withdrawal later .
			err = k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, record.RecipientAddress, types.ModuleName, record.Coins)
			if err != nil {
				return errors.Wrapf(err, "Unable to set Distribution Records to completed : %s", record.String())

			}
			continue
		}
		// Use record details to delete associated claim
		// The claim should always be locked at this point in time .
		if record.DoesClaimExist() {
			k.DeleteClaim(ctx, record.RecipientAddress.String(), record.DistributionType)
		}
		ctx.Logger().Info(fmt.Sprintf("Distributed to : %s | At height : %d | Amount :%s \n", record.RecipientAddress.String(), height, record.Coins.String()))
	}
	return nil
}

// AccumulateDrops collects funds from a senders account and transfers it to the Dispensation module account
func (k Keeper) AccumulateDrops(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) error {
	err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, amount)
	if err != nil {
		return errors.Wrapf(types.ErrFailedInputs, "for address  : %s", addr.String())
	}
	return nil

}

// Verify if the distribution is correct
// The verification is the for distributionName + distributionType
func (k Keeper) VerifyAndSetDistribution(ctx sdk.Context, distributionName string,
	distributionType types.DistributionType, runner sdk.AccAddress) error {

	if k.ExistsDistribution(ctx, distributionName, distributionType) {
		return errors.Wrapf(types.ErrDistribution, "airdrop with same name already exists : %s ", distributionName)
	}
	// Create distribution only if a distribution with the same name does not exist
	err := k.SetDistribution(ctx, types.NewDistribution(distributionType, distributionName, runner))
	if err != nil {
		return errors.Wrapf(types.ErrDistribution, "unable to set airdrop :  %s ", distributionName)
	}

	return nil
}
