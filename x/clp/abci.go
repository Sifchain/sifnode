package clp

import (
	"github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) []abci.ValidatorUpdate {
	params := keeper.GetRewardsParams(ctx)
	pools := keeper.GetPools(ctx)
	currentPeriod := keeper.GetCurrentRewardPeriod(ctx, params)
	if currentPeriod != nil && !currentPeriod.Allocation.IsZero() {
		err := keeper.DistributeDepthRewards(ctx, currentPeriod, pools)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}
