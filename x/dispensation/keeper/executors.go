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
		if k.ExistsDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
			if err != nil {
				return errors.Wrapf(types.ErrDistribution, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			distributionRecord.Add(*oldRecord)
		}
		distributionRecord.DistributionStatus = types.DistributionStatus_DISTRIBUTION_STATUS_PENDING
		err := k.SetDistributionRecord(ctx, distributionRecord)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", distributionRecord.String())
		}
		// Lock the user claim so that the user cannot delete the claim while the distribution is in progress.
		// Claim will not exist if its not a LM/VS drop
		// IF it is a LM/VS drop the associated claim must always exist .
		// The users of this module need to make sure they are submitting the proper distribution type when distributing rewards
		// The same user might be eligible for Airdrop/LM/VS rewards . Based on Distribution type submitted the appropriate claim will be locked.
		if distributionType == types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING || distributionType == types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY {
			err := k.LockClaim(ctx, receiver.Address, distributionType)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to verify associated claim for address : %s", receiver.Address))
			}
		}
	}
	return nil
}

// DistributeDrops is called at the beginning of every block .
// It checks if any pending records are present , if there are it completes the top 10
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
	}
	return nil
}

// Verify if the distribution is correct
// The verification is the for distributionName + distributionType
func (k Keeper) VerifyDistribution(ctx sdk.Context, distributionName string, distributionType types.DistributionType) error {
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
