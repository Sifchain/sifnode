package trees

import (
	"github.com/Sifchain/sifnode/x/trees/keeper"
	"github.com/Sifchain/sifnode/x/trees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data types.GenesisState) {

	return types.NewGenesisState()
}
