package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) InitGenesis(sdk.Context, types.GenesisState) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
