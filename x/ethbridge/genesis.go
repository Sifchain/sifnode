package ethbridge

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis import genesis data from GenesisState
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	for _, token := range data.PeggyTokens {
		keeper.AddPeggyToken(ctx, token)
	}

	if data.CethReceiverAccount != nil {
		keeper.SetCethReceiverAccount(ctx, data.CethReceiverAccount)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis export data into genesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	tokens := keeper.GetPeggyToken(ctx)
	cethReceiverAccount := keeper.GetCethReceiverAccount(ctx)

	return types.GenesisState{
		PeggyTokens:         tokens,
		CethReceiverAccount: cethReceiverAccount,
	}
}

// ValidateGenesis validates the ethbridge genesis parameters
func ValidateGenesis(data types.GenesisState) error {
	return nil
}
