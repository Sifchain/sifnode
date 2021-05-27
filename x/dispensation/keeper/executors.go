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
func (k Keeper) CreateDrops(ctx sdk.Context, output []banktypes.Output, name string, distributionType types.DistributionType) error {
	for _, receiver := range output {
		distributionRecord := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, distributionType, name, receiver.Address, receiver.Coins, ctx.BlockHeight(), -1)
		if k.ExistsDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, distributionRecord.DistributionType) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, distributionRecord.DistributionType)
			if err != nil {
				return errors.Wrapf(types.ErrDistribution, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			distributionRecord = distributionRecord.Add(*oldRecord)
		}
		distributionRecord.DistributionStatus = types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
		err := k.SetDistributionRecord(ctx, distributionRecord)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", distributionRecord.String())
		}
	}
	return nil
}

// DistributeDrops is called at the beginning of every block .
// It checks if any pending records are present , if there are it completes the top 10
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64) error {
	pendingRecords := k.GetRecordsLimitedForStatus(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
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
		err = k.DeleteDistributionRecord(ctx, record.DistributionName, record.RecipientAddress, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType) // Delete the record in the pending prefix so the iteration is cheaper.
		if err != nil {
			return errors.Wrapf(types.ErrDistribution, "error deleting pending distibution record  : %s", record.String())
		}
		// Use record details to delete associated claim
		// The claim should always be locked at this point in time .
		if record.DoesClaimExist() {
			k.DeleteClaim(ctx, record.RecipientAddress, record.DistributionType)
		}
		ctx.Logger().Info(fmt.Sprintf("Distributed to : %s | At height : %d | Amount :%s \n", record.RecipientAddress, height, record.Coins.String()))
	}
	return nil
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
