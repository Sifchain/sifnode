package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type KeeperI interface {
	InitGenesis(sdk.Context, types.GenesisState) []abci.ValidatorUpdate
	ExportGenesis(sdk.Context) *types.GenesisState
}
