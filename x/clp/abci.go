package clp

import (
	"github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	params := k.GetParams(ctx)
	pools := k.GetPools(ctx)
	currentPeriod := k.GetCurrentRewardPeriod(ctx, params)
	if currentPeriod != nil {
		err := k.DistributeDepthRewards(ctx, currentPeriod, pools)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}
