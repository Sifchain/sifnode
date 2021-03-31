package dispensation

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the clp type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgDistribution:
			return handleMsgCreateDistribution(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateDistribution is the top level function for calling all executors.
func handleMsgCreateDistribution(ctx sdk.Context, keeper Keeper, msg MsgDistribution) (*sdk.Result, error) {
	// Verify if distribution already exists
	err := keeper.VerifyDistribution(ctx, msg.DistributionName, msg.DistributionType)
	if err != nil {
		return nil, err
	}
	//Accumulate all Drops into the ModuleAccount
	err = keeper.AccumulateDrops(ctx, msg.Input)
	if err != nil {
		return nil, err
	}
	//Distribute rewards and Store Historical Data
	// TODO Break create and distribute in two different steps so that distribute can be moved to the Block beginner
	err = keeper.CreateAndDistributeDrops(ctx, msg.Output, msg.DistributionName)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDistributionCompleted,
			sdk.NewAttribute(types.AttributeKeyFromModuleAccount, types.GetDistributionModuleAddress().String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
