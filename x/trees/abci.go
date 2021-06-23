package trees

import (
	"github.com/Sifchain/sifnode/x/trees/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.FillLimitOrders(ctx)
}
