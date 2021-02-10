package clp

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, bankKeeper types.BankKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	keeper.SetParams(ctx, data.Params)
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return GenesisState{}
}

// ValidateGenesis validates the clp genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
