//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: &Params{
			LeverageMax:          sdk.NewUint(2),
			HealthGainFactor:     sdk.NewDec(1),
			InterestRateMin:      sdk.NewDecWithPrec(5, 3),
			InterestRateMax:      sdk.NewDec(3),
			InterestRateDecrease: sdk.NewDecWithPrec(1, 1),
			InterestRateIncrease: sdk.NewDecWithPrec(1, 1),
			ForceCloseThreshold:  sdk.NewDecWithPrec(1, 1),
			EpochLength:          1,
			Pools:                []string{},
		},
	}
}