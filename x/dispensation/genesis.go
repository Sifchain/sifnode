package dispensation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// TODO Add import and export state
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) (res []abci.ValidatorUpdate) {
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) *GenesisState {
	return &GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
