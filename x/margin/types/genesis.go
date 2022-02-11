package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: &Params{
			LeverageMax:          sdk.NewUint(1),
			HealthGainFactor:     sdk.NewDec(1),
			InterestRateMin:      sdk.NewDec(0.005),
			InterestRateMax:      sdk.NewDec(3),
			InterestRateDecrease: sdk.NewDec(0.1),
			InterestRateIncrease: sdk.NewDec(0.1),
			ForceCloseThreshold:  sdk.NewDec(0.1),
			EpochLength:          1,
			Pools:                []string{},
		},
	}
}
