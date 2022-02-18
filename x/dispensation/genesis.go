package dispensation

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"math/rand"
	"strconv"
	"time"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	CreatePerfTestData(ctx, keeper)
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

func CreatePerfTestData(ctx sdk.Context, keeper Keeper) {
	keeper.SetMaxRecordsPerBlock(ctx, types.MaxRecordsPerBlock{
		MaxRecords: types.MaxRecordsPerBlockConst,
	})
	distributionName := "test_dist"
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	runner := "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	err := keeper.SetDistribution(ctx, types.NewDistribution(
		distributionType,
		distributionName,
		runner,
	))
	if err != nil {
		panic(err)
	}
	// Calculate using N/2 (2a + (N-1)D)
	// N = 1000 (Max runs per block ) / 10 (Step increment )
	outputList := CreatOutputListGen(50500, "1000") //1252500
	totalOutput, err := dispensationUtils.TotalOutput(outputList)
	if err != nil {
		panic(err)
	}
	err = keeper.GetBankKeeper().MintCoins(ctx, types.ModuleName, totalOutput)
	if err != nil {
		panic(err)
	}
	err = keeper.CreateDrops(ctx, outputList, distributionName, distributionType, runner)
	if err != nil {
		panic(err)
	}
}

func CreatOutputListGen(count int, rowanAmount string) []banktypes.Output {
	outputList := make([]banktypes.Output, count)
	amount, ok := sdk.NewIntFromString(rowanAmount)
	if !ok {
		panic("Unable to generate rowan amount")
	}
	coin := sdk.NewCoins(sdk.NewCoin("rowan", amount), sdk.NewCoin("ceth", amount), sdk.NewCoin("catk", amount))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		address := sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i))))
		out := banktypes.NewOutput(address, sdk.NewCoins(coin[0]))
		outputList[i] = out
	}
	return outputList
}
