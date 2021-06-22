package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
)

//CreateAndDistributeDrops creates new drop Records . These records are then used to facilitate distribution
// Each Recipient and DropName generate a unique Record
func (k Keeper) CreateDrops(ctx sdk.Context, output []banktypes.Output, name string, distributionType types.DistributionType, authorisedRunner string) error {
	for _, receiver := range output {
		distributionRecord := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, distributionType, name, receiver.Address, receiver.Coins, ctx.BlockHeight(), -1, authorisedRunner)
		if k.ExistsDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, distributionRecord.DistributionType) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, distributionRecord.DistributionType)
			if err != nil {
				return errors.Wrapf(types.ErrDistribution, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			distributionRecord = distributionRecord.Add(*oldRecord)
		}
		err := k.SetDistributionRecord(ctx, distributionRecord)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", distributionRecord.String())
		}
	}
	return nil
}

// DistributeDrops is called at the beginning of every block .
// It checks if any pending records are present , if there are it completes the top 10
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64, distributionName string, authorisedRunner string, distributionType types.DistributionType) (*types.DistributionRecords, error) {
	pendingRecords := k.GetLimitedRecordsForRunner(ctx, distributionName, authorisedRunner, distributionType, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	for _, record := range pendingRecords.DistributionRecords {
		recipientAddress, err := sdk.AccAddressFromBech32(record.RecipientAddress)
		if err != nil {
			err := errors.Wrapf(err, "Invalid address for distribute : %s", record.RecipientAddress)
			ctx.Logger().Error(err.Error())
			continue
		}
		err = k.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddress, record.Coins)
		if err != nil {
			err := errors.Wrapf(types.ErrFailedOutputs, "for address  : %s", record.RecipientAddress)
			ctx.Logger().Error(err.Error())
			err = k.ChangeRecordStatus(ctx, *record, height, types.DistributionStatus_DISTRIBUTION_STATUS_FAILED)
			if err != nil {
				panic(fmt.Sprintf("Unable to set Distribution Records to Failed : %s", record.String()))
			}
			continue
		}

		err = k.ChangeRecordStatus(ctx, *record, height, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
		if err != nil {
			err := errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", record.String())
			ctx.Logger().Error(err.Error())
			// If the SetDistributionRecord returns error , that would mean the required amount was transferred to the user , but the record was not set to completed .
			// In this case we try to take the funds back from the user , and attempt the withdrawal later .
			err = k.GetBankKeeper().SendCoinsFromAccountToModule(ctx, recipientAddress, types.ModuleName, record.Coins)
			if err != nil {
				panic(fmt.Sprintf("Unable to set Distribution Records to completed : %s", record.String()))
			}
			continue
		}
		// Use record details to delete associated claim
		if record.DoesTypeSupportClaim() {
			k.DeleteClaim(ctx, record.RecipientAddress, record.DistributionType)
		}
		ctx.Logger().Info(fmt.Sprintf("Distributed to : %s | At height : %d | Amount :%s \n", record.RecipientAddress, height, record.Coins.String()))
	}
	return pendingRecords, nil
}

// AccumulateDrops collects funds from a senders account and transfers it to the Dispensation module account
func (k Keeper) AccumulateDrops(ctx sdk.Context, addr string, amount sdk.Coins) error {
	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return errors.Wrapf(err, "Invalid address for distribute : %s", addr)
	}
	err = k.GetBankKeeper().SendCoinsFromAccountToModule(ctx, address, types.ModuleName, amount)
	if err != nil {
		return errors.Wrapf(types.ErrFailedInputs, "for address  : %s", addr)
	}
	return nil
}

// Verify if the distribution is correct
// The verification is the for distributionName + distributionType
func (k Keeper) VerifyAndSetDistribution(ctx sdk.Context, distributionName string, distributionType types.DistributionType) error {
	if k.ExistsDistribution(ctx, distributionName, distributionType) {
		return errors.Wrapf(types.ErrDistribution, "airdrop with same name already exists : %s ", distributionName)
	}
	// Create distribution only if a distribution with the same name does not exist
	err := k.SetDistribution(ctx, types.NewDistribution(distributionType, distributionName))
	if err != nil {
		return errors.Wrapf(types.ErrDistribution, "unable to set airdrop :  %s ", distributionName)

	}
	return nil
}
