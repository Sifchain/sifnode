package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// TODO review default param values
		Params: &Params{
			LeverageMax:          sdk.NewUint(1),
			HealthGainFactor:     sdk.NewDec(1),
			InterestRateMin:      sdk.NewDec(1),
			InterestRateMax:      sdk.NewDec(1),
			InterestRateDecrease: sdk.NewDec(1),
			InterestRateIncrease: sdk.NewDec(1),
			EpochLength:          1,
			// TODO start with empty slice for pools
			Pools: []string{"ceth"},
		},
	}
}
