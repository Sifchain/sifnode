package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) QueueRemoval(ctx sdk.Context, msg *types.MsgRemoveLiquidity) {
	queue := k.GetRemovalQueue(ctx)
	request := types.RemovalRequest{
		Id:  queue.Id + 1,
		Msg: msg,
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRemovalRequestKey(request), k.cdc.MustMarshal(&request))

	queue.Count += 1
	queue.Id += 1
	k.SetRemovalQueue(ctx, queue)
}

func (k Keeper) GetRemovalQueue(ctx sdk.Context) types.RemovalQueue {
	var queue types.RemovalQueue
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RemovalQueuePrefix)
	k.cdc.MustUnmarshal(bz, &queue)
	return queue
}

func (k Keeper) SetRemovalQueue(ctx sdk.Context, queue types.RemovalQueue) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RemovalQueuePrefix, k.cdc.MustMarshal(&queue))
}

func (k Keeper) GetRemovalQueueIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.RemovalRequestPrefix)
}

func (k Keeper) ProcessRemovalQueue(ctx sdk.Context, msg *types.MsgAddLiquidity, unitsToDistribute sdk.Uint) {
	perRequestUnits := unitsToDistribute.Quo(sdk.NewUint(uint64(k.GetRemovalQueue(ctx).Count)))

	it := k.GetRemovalQueueIterator(ctx)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		var request types.RemovalRequest
		k.cdc.MustUnmarshal(it.Value(), &request)

		lp, err := k.GetLiquidityProvider(ctx, request.Msg.ExternalAsset.Symbol, request.Msg.Signer)
		if err != nil {
			continue
		}

		requestUnits := ConvWBasisPointsToUnits(lp.LiquidityProviderUnits, request.Msg.WBasisPoints)
		withdrawUnits := sdk.MinUint(requestUnits, perRequestUnits)
		withdrawWBasisPoints := ConvUnitsToWBasisPoints(lp.LiquidityProviderUnits, withdrawUnits)

		// Reuse removal logic using withdrawWBasisPoints
		/*k.ProcessRemoveLiqiduityMsg(ctx, types.MsgRemoveLiquidity{
			Signer:        request.Msg.Signer,
			ExternalAsset: request.Msg.ExternalAsset,
			WBasisPoints:  withdrawWBasisPoints,
			Asymmetry:     request.Msg.Asymmetry,
		})*/

		// Update the queued request
		k.SetProcessedRemovalRequest(ctx, request, withdrawWBasisPoints /*, rowanValue*/)
	}
}

func (k Keeper) SetProcessedRemovalRequest(ctx sdk.Context, request types.RemovalRequest, pointsProcessed sdk.Int /*, rowanRemoved sdk.Uint*/) {
	request.Msg.WBasisPoints = request.Msg.WBasisPoints.Sub(pointsProcessed)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRemovalRequestKey(request), k.cdc.MustMarshal(&request))
	//k.DecrementRemovalQueueTotalValue(ctx, rowanRemoved)
	if request.Msg.WBasisPoints.LT(sdk.ZeroInt()) {
		k.DequeueRemovalRequest(ctx, request)
	}

	// emitProcessedRemovalRequest(request, pointsProcessed)
}

func (k Keeper) DequeueRemovalRequest(ctx sdk.Context, request types.RemovalRequest) {
	ctx.KVStore(k.storeKey).Delete(types.GetRemovalRequestKey(request))
	queue := k.GetRemovalQueue(ctx)
	queue.Count -= 1
	k.SetRemovalQueue(ctx, queue)
}
