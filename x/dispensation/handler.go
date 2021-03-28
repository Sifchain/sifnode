package dispensation

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the clp type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgDistribution:
			return handleMsgDistribution(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDistribution(ctx sdk.Context, keeper Keeper, msg MsgDistribution) (*sdk.Result, error) {
	// Verify if airdrop already exists
	err := keeper.VerifyDistribution(ctx, msg.DistributionName)
	if err != nil {
		return nil, err
	}
	//Accumulate all Drops into the ModuleAccount
	err = keeper.AccumulateDrops(ctx, msg.Input)
	if err != nil {
		return nil, err
	}
	//Distribute rewards and Store Historical Data
	// TODO Break create and Distribute in two different steps so that distribute can be moved to the Block beginner
	err = keeper.CreateAndDistributeDrops(ctx, msg.Output, msg.DistributionName)
	if err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
