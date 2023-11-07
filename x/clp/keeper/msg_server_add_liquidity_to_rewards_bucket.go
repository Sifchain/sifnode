package keeper

import (
	"context"
	"strconv"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddLiquidityToRewardsBucket(goCtx context.Context, msg *types.MsgAddLiquidityToRewardsBucketRequest) (*types.MsgAddLiquidityToRewardsBucketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addedCoins, err := k.Keeper.AddLiquidityToRewardsBucket(ctx, msg.Signer, msg.Amount)
	if err != nil {
		return nil, err
	}

	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	}

	// emit as many events as addedCoins exist
	for _, coin := range addedCoins {
		// emit event for each coin added to rewards bucket
		events.AppendEvent(
			sdk.NewEvent(
				types.EventTypeAddLiquidityToRewardsBucket,
				sdk.NewAttribute(types.AttributeKeyAmount, coin.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, coin.Denom),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
		)
	}

	ctx.EventManager().EmitEvents(events)

	return &types.MsgAddLiquidityToRewardsBucketResponse{}, nil
}
