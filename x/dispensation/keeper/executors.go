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
func (k Keeper) CreateDrops(ctx sdk.Context, output []banktypes.Output, name string) error {
	for _, receiver := range output {
		distributionRecord := types.NewDistributionRecord(name, "", receiver.Coins, sdk.NewInt(ctx.BlockHeight()), sdk.NewInt(-1))
		if k.ExistsDistributionRecord(ctx, name, receiver.Address) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address)
			if err != nil {
				return errors.Wrapf(types.ErrDistribution, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			distributionRecord.Add(oldRecord)
		}
		distributionRecord.ClaimStatus = sdk.NewUint(1)
		err := k.SetDistributionRecord(ctx, distributionRecord)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", distributionRecord.String())
		}
	}
	return nil
}

//DistributeDrops is called at the beginning of every block .
// It checks if any pending records are present , if there are it completes the top 10
func (k Keeper) DistributeDrops(ctx sdk.Context, height int64) error {
	pendingRecords := k.GetPendingRecordsLimited(ctx, 10)
	for _, record := range pendingRecords.DistributionRecords {
		err := k.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(record.RecipientAddress), record.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "for address  : %s", record.RecipientAddress)
		}
		record.ClaimStatus = sdk.NewUint(2)
		record.DistributionCompletedHeight = sdk.NewInt(height)
		err = k.SetDistributionRecord(ctx, *record)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", record)
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
