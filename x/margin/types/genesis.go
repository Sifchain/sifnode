package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// TODO review default param values
		Params: &Params{
			LeverageMax: sdk.NewUint(1),
		},
	}
}
