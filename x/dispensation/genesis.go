package dispensation

import (
	"fmt"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	keeper.SetMintController(ctx,
		types.MintController{TotalCounter: sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.ZeroInt())})
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

	if data.Claims != nil {
		for _, claim := range data.Claims.UserClaims {
			err := keeper.SetClaim(ctx, *claim)
			if err != nil {
				panic(fmt.Sprintf("Error setting claim during init genesis : %s", claim.String()))
			}
		}
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return GenesisState{
		Distributions:       keeper.GetDistributions(ctx),
		DistributionRecords: keeper.GetRecords(ctx),
		Claims:              keeper.GetClaims(ctx),
	}
}

func ValidateGenesis(data GenesisState) error {
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
	if data.Claims != nil {
		for _, claim := range data.Claims.UserClaims {
			if !claim.Validate() {
				return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Claim is invalid : %s", claim.String()))
			}
		}
	}

	return nil
}
