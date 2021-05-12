package dispensation

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// Distribute drops if any are pending
	err := k.DistributeDrops(ctx, req.Header.Height)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Error Distributing : %s | Height : %d ", err, req.Header.Height))
		return
	}
}
