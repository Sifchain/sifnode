package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) QueueRemoval(ctx sdk.Context, msg *types.MsgRemoveLiquidity, rowanValue sdk.Uint) {
	queue := k.GetRemovalQueue(ctx)
	request := types.RemovalRequest{
		Id:    queue.Id + 1,
		Msg:   msg,
		Value: rowanValue,
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRemovalRequestKey(request), k.cdc.MustMarshal(&request))

	queue.Count += 1
	queue.Id += 1
	queue.TotalValue = queue.TotalValue.Add(rowanValue)
	if queue.Count == 1 {
		queue.StartHeight = ctx.BlockHeight()
	}
	k.SetRemovalQueue(ctx, queue)

	emitQueueRemoval(ctx, &request, &queue)
}

func (k Keeper) GetRemovalQueue(ctx sdk.Context) types.RemovalQueue {
	queue := types.RemovalQueue{
		Count:       0,
		Id:          0,
		StartHeight: 0,
		TotalValue:  sdk.ZeroUint(),
	}
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

		if msg.ExternalAsset.Equals(*request.Msg.ExternalAsset) {
			lp, err := k.GetLiquidityProvider(ctx, request.Msg.ExternalAsset.Symbol, request.Msg.Signer)
			if err != nil {
				continue
			}

			requestUnits := ConvWBasisPointsToUnits(lp.LiquidityProviderUnits, request.Msg.WBasisPoints)
			withdrawUnits := sdk.MinUint(requestUnits, perRequestUnits)
			withdrawWBasisPoints := ConvUnitsToWBasisPoints(lp.LiquidityProviderUnits, withdrawUnits)

			// Reuse removal logic using withdrawWBasisPoints
			_, _, totalRowanValue, err := k.ProcessRemoveLiquidityMsg(ctx, &types.MsgRemoveLiquidity{
				Signer:        request.Msg.Signer,
				ExternalAsset: msg.ExternalAsset,
				WBasisPoints:  withdrawWBasisPoints,
				Asymmetry:     request.Msg.Asymmetry,
			})

			// Update the queued request
			k.SetProcessedRemovalRequest(ctx, request, withdrawWBasisPoints, totalRowanValue)
		}
	}
}

func (k Keeper) SetProcessedRemovalRequest(ctx sdk.Context, request types.RemovalRequest, pointsProcessed sdk.Int, rowanRemoved sdk.Uint) {
	request.Msg.WBasisPoints = request.Msg.WBasisPoints.Sub(pointsProcessed)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRemovalRequestKey(request), k.cdc.MustMarshal(&request))

	queue := k.GetRemovalQueue(ctx)
	queue.TotalValue = queue.TotalValue.Sub(rowanRemoved)
	k.SetRemovalQueue(ctx, queue)

	if request.Msg.WBasisPoints.LT(sdk.ZeroInt()) {
		k.DequeueRemovalRequest(ctx, request)
	}

	emitProcessedRemovalRequest(ctx, &request, pointsProcessed, rowanRemoved)
}

func (k Keeper) DequeueRemovalRequest(ctx sdk.Context, request types.RemovalRequest) {
	ctx.KVStore(k.storeKey).Delete(types.GetRemovalRequestKey(request))
	queue := k.GetRemovalQueue(ctx)
	queue.Count -= 1
	k.SetRemovalQueue(ctx, queue)

	emitDequeueRemoval(ctx, &request, &queue)
}

func (k Keeper) GetRemovalQueueUnitsForLP(ctx sdk.Context, lp types.LiquidityProvider) sdk.Uint {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetRemovalRequestLPPrefix(lp.LiquidityProviderAddress)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	units := sdk.ZeroUint()
	for ; iterator.Valid(); iterator.Next() {
		var request types.RemovalRequest
		k.cdc.MustUnmarshal(iterator.Value(), &request)

		units = units.Add(ConvWBasisPointsToUnits(lp.LiquidityProviderUnits, request.Msg.WBasisPoints))
	}

	return units
}
