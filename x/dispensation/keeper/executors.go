package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pkg/errors"
)

func (k Keeper) CreateAndDistributeDrops(ctx sdk.Context, output []bank.Output, airDropName string) error {
	for _, receiver := range output {
		err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver.Address, receiver.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedOutputs, "for address  : %s", receiver.Address.String())
		}
		distributionRecord := types.NewDistributionRecord(airDropName, receiver.Address, receiver.Coins)
		if k.ExistsDistributionRecord(ctx, airDropName, receiver.Address.String()) {
			oldRecord, err := k.GetDistributionRecord(ctx, airDropName, receiver.Address.String())
			if err != nil {
				return errors.Wrapf(types.ErrAirdrop, "failed appending record for : %s", distributionRecord.Address)
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

func (k Keeper) AccumulateDrops(ctx sdk.Context, input []bank.Input) error {
	for _, fundingInput := range input {
		err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, fundingInput.Address, types.ModuleName, fundingInput.Coins)
		if err != nil {
			return errors.Wrapf(types.ErrFailedInputs, "for address  : %s", fundingInput.Address.String())
		}
	}
	return nil
}

func (k Keeper) VerifyAirdrop(ctx sdk.Context, airDropName string) error {
	if k.ExistsAirdrop(ctx, airDropName) {
		return errors.Wrapf(types.ErrAirdrop, "airdrop with same name already exists : %s ", airDropName)
	}
	err := k.SetAirdropRecord(ctx, types.NewAirdropRecord(airDropName))
	if err != nil {
		return errors.Wrapf(types.ErrAirdrop, "unable to set airdrop :  %s ", airDropName)

	}
	return nil
}
