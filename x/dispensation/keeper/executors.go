package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pkg/errors"
)

//CreateAndDistributeDrops creates new drop Records . These records are then used to facilitate distribution
// Each Recipient and DropName generate a unique Record

func (k Keeper) CreateAndDistributeDrops(ctx sdk.Context, output []bank.Output, name string) error {
	for _, receiver := range output {
		err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver.Address, receiver.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "for address  : %s", receiver.Address.String())
		}
		distributionRecord := types.NewDistributionRecord(name, receiver.Address, receiver.Coins)
		if k.ExistsDistributionRecord(ctx, name, receiver.Address.String()) {
			oldRecord, err := k.GetDistributionRecord(ctx, name, receiver.Address.String())
			if err != nil {
				return errors.Wrapf(types.ErrAirdrop, "failed appending record for : %s", distributionRecord.RecipientAddress)
			}
			distributionRecord.Add(oldRecord)
		}
		err = k.SetDistributionRecord(ctx, distributionRecord)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "error setting distibution record  : %s", distributionRecord.String())
		}
	}
	return nil
}

// Accumulate Drops from the sender accounts
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
func (k Keeper) VerifyDistribution(ctx sdk.Context, name string, t types.DistributionType) error {
	if k.ExistsDistribution(ctx, name) {
		return errors.Wrapf(types.ErrAirdrop, "airdrop with same name already exists : %s ", name)
	}
	// Create distribution only if a distribution with the same name does not exist
	err := k.SetDistribution(ctx, types.NewDistribution(t, name))
	if err != nil {
		return errors.Wrapf(types.ErrAirdrop, "unable to set airdrop :  %s ", name)

	}
	return nil
}
