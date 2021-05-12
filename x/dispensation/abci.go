package dispensation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock, _ Keeper) {
	// Distribute drops if any are pending
	//_ = k.DistributeDrops(ctx, req.Header.Height)
}
