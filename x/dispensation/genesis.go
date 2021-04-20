package dispensation

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// TODO Add import and export state
func InitGenesis(ctx sdk.Context, keeper Keeper, data cdctypes.Any) (res []abci.ValidatorUpdate) {
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) *cdctypes.Any {
	return &cdctypes.Any{}
}

func ValidateGenesis(data cdctypes.Any) error {
	return nil
}
