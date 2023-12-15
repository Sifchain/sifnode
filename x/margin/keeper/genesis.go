package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) []abci.ValidatorUpdate {
	// Set initial Margin parameters
	k.SetParams(ctx, data.Params)

	// Set all the mtps
	for _, mtp := range data.MtpList {
		err := k.SetMTP(ctx, mtp)
		if err != nil {
			panic(err)
		}
	}

	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// Retrieve the Margin parameters
	params := k.GetParams(ctx)

	// Retrieve all the mtps
	mtps := k.GetAllMTPS(ctx)

	return &types.GenesisState{
		Params:  &params,
		MtpList: mtps,
	}
}
