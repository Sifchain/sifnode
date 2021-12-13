package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)

	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
