package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pkg/errors"
)

//CreateAndDistributeDrops creates new drop Records . These records are then used to facilitate distribution
// Each Recipient and DropName generate a unique Record
func (k Keeper) CreateDrops(ctx sdk.Context, output []bank.Output, name string) error {
	return errors.New("Dispensation module is disabled")
	for _, receiver := range output {
		distributionRecord := types.NewDistributionRecord(name, distributionType, receiver.Address, receiver.Coins, ctx.BlockHeight(), -1)
		if k.ExistsDistributionRecord(ctx, name, receiver.Address.String()) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address.String())
			if err != nil {
				return errors.Wrapf(types.ErrDistribution, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			distributionRecord.Add(oldRecord)
		}
		distributionRecord.DistributionStatus = types.Pending
		err := k.SetDistributionRecord(ctx, distributionRecord)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", distributionRecord.String())
		}
		// Lock the user claim so that the user cannot delete the claim while the distribution is in progress.
		// Claim will not exist if its not a LM/VS drop
		// IF it is a LM/VS drop the associated claim must always exist .
		// The users of this module need to make sure they are submitting the proper distribution type when distributing rewards
		// The same user might be eligible for Airdrop/LM/VS rewards . Based on Distribution type submitted the appropriate claim will be locked.
		if distributionType == types.LiquidityMining || distributionType == types.ValidatorSubsidy {
			err := k.LockClaim(ctx, receiver.Address.String(), distributionType)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to verify associated claim for address : %s", receiver.Address.String()))
			}
		}
	}
	return nil
}

// DistributeDrops is called at the beginning of every block .
// It checks if any pending records are present , if there are it completes the top 10
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64) error {
	return errors.New("Dispensation module is disabled")
	pendingRecords := k.GetPendingRecordsLimited(ctx, 10)
	for _, record := range pendingRecords {
		err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, record.RecipientAddress, record.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "for address  : %s", record.RecipientAddress.String())
		}
		record.DistributionStatus = types.Completed
		record.DistributionCompletedHeight = height
		err = k.SetDistributionRecord(ctx, record)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", record.String())
		}
		// Use record details to delete associated claim
		// The claim should always be locked at this point in time .
		if record.DistributionType == types.LiquidityMining || record.DistributionType == types.ValidatorSubsidy {
			k.DeleteClaim(ctx, record.RecipientAddress.String(), record.DistributionType)
		}
		ctx.Logger().Info(fmt.Sprintf("Distributed to : %s | At height : %d | Amount :%s \n", record.RecipientAddress.String(), height, record.Coins.String()))
	}
	return nil
}

// AccumulateDrops collects funds from a senders account and transfers it to the Dispensation module account
func (k Keeper) AccumulateDrops(ctx sdk.Context, input []bank.Input) error {
	for _, fundingInput := range input {
		err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, fundingInput.Address, types.ModuleName, fundingInput.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedInputs, "for address  : %s", fundingInput.Address.String())
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
