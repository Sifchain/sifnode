package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func CreateEventMsg(signer string) sdk.Event {
	return sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, signer))
}

func CreateEventBlockHeight(ctx sdk.Context, eventType string, attribute sdk.Attribute) sdk.Event {
	return sdk.NewEvent(
		eventType,
		attribute,
		sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	)
}

func emitProcessedRemovalRequest(ctx sdk.Context, request *types.RemovalRequest, points sdk.Int, rowanRemoved sdk.Uint) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeProcessedRemovalQueue,
		sdk.NewAttribute("id", strconv.FormatInt(request.Id, 10)),
		sdk.NewAttribute(types.AttributeKeyLiquidityProvider, request.Msg.Signer),
		sdk.NewAttribute(types.AttributeKeyPool, request.Msg.ExternalAsset.Symbol),
		sdk.NewAttribute("points_requested", request.Msg.WBasisPoints.String()),
		sdk.NewAttribute("points_processed", points.String()),
		sdk.NewAttribute("value_in_rowan_processed", rowanRemoved.String()),
	))
}

func emitQueueRemoval(ctx sdk.Context, request *types.RemovalRequest, queue *types.RemovalQueue) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeQueueRemovalRequest,
		sdk.NewAttribute("id", strconv.FormatInt(request.Id, 10)),
		sdk.NewAttribute("rowan_value", request.Value.String()),
		sdk.NewAttribute(types.AttributeKeyLiquidityProvider, request.Msg.Signer),
		sdk.NewAttribute(types.AttributeKeyPool, request.Msg.ExternalAsset.Symbol),
		sdk.NewAttribute("points_requested", request.Msg.WBasisPoints.String()),
		sdk.NewAttribute("asymmetry", request.Msg.Asymmetry.String()),
	))
}

func emitDequeueRemoval(ctx sdk.Context, request *types.RemovalRequest, queue *types.RemovalQueue) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeDequeueRemovalRequest,
		sdk.NewAttribute("id", strconv.FormatInt(request.Id, 10)),
		sdk.NewAttribute("rowan_value", request.Value.String()),
		sdk.NewAttribute(types.AttributeKeyLiquidityProvider, request.Msg.Signer),
		sdk.NewAttribute(types.AttributeKeyPool, request.Msg.ExternalAsset.Symbol),
	))
}

func emitRemovalQueueError(ctx sdk.Context, request *types.RemovalRequest) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeProcessRemovalError,
		sdk.NewAttribute("id", strconv.FormatInt(request.Id, 10)),
		sdk.NewAttribute(types.AttributeKeyLiquidityProvider, request.Msg.Signer),
		sdk.NewAttribute(types.AttributeKeyPool, request.Msg.ExternalAsset.Symbol),
		sdk.NewAttribute("points_requested", request.Msg.WBasisPoints.String()),
	))
}
