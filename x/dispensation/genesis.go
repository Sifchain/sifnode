package dispensation

import (
	"fmt"
<<<<<<< HEAD

=======
>>>>>>> develop
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
<<<<<<< HEAD
	if data.DistributionRecords != nil {
		for _, record := range data.DistributionRecords.DistributionRecords {
			err := keeper.SetDistributionRecord(ctx, *record)
			if err != nil {
				panic(fmt.Sprintf("Error setting distribution record during init genesis : %s", record.String()))
			}
		}
	}
	if data.Distributions != nil {
		for _, dist := range data.Distributions.Distributions {
			err := keeper.SetDistribution(ctx, *dist)
			if err != nil {
				panic(fmt.Sprintf("Error setting distribution during init genesis : %s", dist.String()))
			}
		}
	}

=======
	for _, record := range data.DistributionRecords {
		err := keeper.SetDistributionRecord(ctx, record)
		if err != nil {
			panic(fmt.Sprintf("Error setting distribution record during init genesis : %s", record.String()))
		}
	}
	for _, dist := range data.Distributions {
		err := keeper.SetDistribution(ctx, dist)
		if err != nil {
			panic(fmt.Sprintf("Error setting distribution during init genesis : %s", dist.String()))
		}
	}
	for _, claim := range data.Claims {
		err := keeper.SetClaim(ctx, claim)
		if err != nil {
			panic(fmt.Sprintf("Error setting claim during init genesis : %s", claim.String()))
		}
	}
>>>>>>> develop
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return GenesisState{
		Distributions:       keeper.GetDistributions(ctx),
		DistributionRecords: keeper.GetRecords(ctx),
<<<<<<< HEAD
=======
		Claims:              keeper.GetClaims(ctx),
>>>>>>> develop
	}
}

func ValidateGenesis(data GenesisState) error {
<<<<<<< HEAD
	if data.DistributionRecords != nil {
		for _, record := range data.DistributionRecords.DistributionRecords {
			if !record.Validate() {
				return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Record is invalid : %s", record.String()))
			}
		}

	}
	if data.Distributions != nil {
		for _, dist := range data.Distributions.Distributions {
			if !dist.Validate() {
				return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Distribution is invalid : %s", dist.String()))
			}
		}
	}

=======
	for _, record := range data.DistributionRecords {
		if !record.Validate() {
			return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Record is invalid : %s", record.String()))
		}
	}
	for _, dist := range data.Distributions {
		if !dist.Validate() {
			return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Distribution is invalid : %s", dist.String()))
		}
	}
	for _, claim := range data.Claims {
		if !claim.Validate() {
			return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Claim is invalid : %s", claim.String()))
		}
	}
>>>>>>> develop
	return nil
}
