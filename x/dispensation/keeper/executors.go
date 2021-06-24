package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
)

//CreateDrops creates new drop Records . These records are then used to facilitate distribution
// Each Recipient and DropName generate a unique Record
<<<<<<< HEAD
	for _, receiver := range output {
		distributionRecord := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, distributionType, name, receiver.Address, receiver.Coins, ctx.BlockHeight(), -1)
		if k.ExistsDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
			if err != nil {
				return errors.Wrapf(types.ErrDistribution, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			distributionRecord.Add(*oldRecord)
		}
		distributionRecord.DistributionStatus = types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
		err := k.SetDistributionRecord(ctx, distributionRecord)
=======
func (k Keeper) CreateDrops(ctx sdk.Context, output []bank.Output, name string, distributionType types.DistributionType, runner sdk.AccAddress) error {
	for _, receiver := range output {
		distributionRecord := types.NewDistributionRecord(name, distributionType, receiver.Address, receiver.Coins, ctx.BlockHeight(), -1, runner)
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
>>>>>>> develop
}

// DistributeDrops is called at the beginning of every block .
// It checks if any pending records are present , if there are it completes the top 10
<<<<<<< HEAD
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64) error {
	pendingRecords := k.GetRecordsLimited(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	for _, record := range pendingRecords.DistributionRecords {
		recipientAddress, err := sdk.AccAddressFromBech32(record.RecipientAddress)
		if err != nil {
			return errors.Wrapf(err, "Invalid address for distribute : %s", record.RecipientAddress)
		}
		err = k.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddress, record.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "for address  : %s", record.RecipientAddress)
		}
		record.DistributionStatus = types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED
		record.DistributionCompletedHeight = height
		// Setting to completed prefix
		err = k.SetDistributionRecord(ctx, *record)
		if err != nil {
			return errors.Wrapf(types.ErrDistribution, "error setting distibution record  : %s", record.String())
		}
		// Deleting from Pending prefix
		err = k.DeleteDistributionRecord(ctx, record.DistributionName, record.RecipientAddress, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING) // Delete the record in the pending prefix so the iteration is cheaper.
		if err != nil {
			return errors.Wrapf(types.ErrDistribution, "error deleting pending distibution record  : %s", record.String())
		}
		// Use record details to delete associated claim
		// The claim should always be locked at this point in time .
		if record.DistributionType == types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING || record.DistributionType == types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY {
			k.DeleteClaim(ctx, record.RecipientAddress, record.DistributionType)
		}
		ctx.Logger().Info(fmt.Sprintf("Distributed to : %s | At height : %d | Amount :%s \n", record.RecipientAddress, height, record.Coins.String()))
	}
	return nil
}

// AccumulateDrops collects funds from a senders account and transfers it to the Dispensation module account
func (k Keeper) AccumulateDrops(ctx sdk.Context, input []banktypes.Input) error {
	for _, fundingInput := range input {
		err := k.GetBankKeeper().SendCoinsFromAccountToModule(ctx, sdk.AccAddress(fundingInput.Address), types.ModuleName, fundingInput.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedInputs, "for address  : %s", fundingInput.Address)
		}
=======
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64, DistributionName string, runner sdk.AccAddress, distributionType types.DistributionType) (types.DistributionRecords, error) {
	// TODO replace 10 with a variable declared in genesis or a constant in keys.go.
	pendingRecords := k.GetRecordsForNamePendingLimited(ctx, DistributionName, 10, runner, distributionType)
	for _, record := range pendingRecords {
		err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, record.RecipientAddress, record.Coins)
		if err != nil {
			err := errors.Wrapf(err, "Distribution failed for address  : %s", record.RecipientAddress.String())
			ctx.Logger().Error(err.Error())
			err = k.MoveRecordToFailed(ctx, record)
			if err != nil {
				return pendingRecords, errors.Wrapf(err, "Unable to set Distribution Records to Failed : %s", record.String())
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
				return pendingRecords, errors.Wrapf(err, "Unable to set Distribution Records to completed : %s", record.String())
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
	return pendingRecords, nil
}

// AccumulateDrops collects funds from a senders account and transfers it to the Dispensation module account
func (k Keeper) AccumulateDrops(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) error {
	err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, amount)
	if err != nil {
		return errors.Wrapf(types.ErrFailedInputs, "for address  : %s", addr.String())
>>>>>>> develop
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
