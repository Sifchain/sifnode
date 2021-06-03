package dispensation

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

// NewHandler creates an sdk.Handler for all the clp type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgDistribution:
			return handleMsgCreateDistribution(ctx, k, msg)
		case MsgCreateClaim:
			return handleMsgCreateClaim(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

//handleMsgRunDistribution 
func handleMsgRunDistribution(ctx sdk.Context, keeper Keeper, msg MsgRunDistribution) (*sdk.Result, error) {
	_ = k.DistributeDrops(ctx, ctx.BlockHeight())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDistributionRun,
			sdk.NewAttribute(types.AttributeKeyDistributionName, msg.distributionName),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

//handleMsgCreateDistribution is the top level function for calling all executors.
func handleMsgCreateDistribution(ctx sdk.Context, keeper Keeper, msg MsgDistribution) (*sdk.Result, error) {
	// Verify if distribution already exists
	distributionName := fmt.Sprintf("%d_%s", ctx.BlockHeight(), msg.Distributor.String())
	err := keeper.VerifyAndSetDistribution(ctx, distributionName, msg.DistributionType, msg.Runner)
	if err != nil {
		return nil, err
	}
	//Accumulate all Drops into the ModuleAccount
	totalOutput, err := dispensationUtils.TotalOutput(msg.Output)
	if err != nil {
		return nil, err
	}
	err = keeper.AccumulateDrops(ctx, msg.Distributor, totalOutput)
	if err != nil {
		return nil, err
	}
	//Create drops and Store Historical Data
	err = keeper.CreateDrops(ctx, msg.Output, distributionName, msg.DistributionType)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDistributionStarted,
			sdk.NewAttribute(types.AttributeKeyFromModuleAccount, types.GetDistributionModuleAddress().String()),
			sdk.NewAttribute(types.AttributeKeyDistributionName, distributionName),
			sdk.NewAttribute(types.AttributeKeyDistributionType, msg.DistributionType.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateClaim(ctx sdk.Context, keeper Keeper, msg MsgCreateClaim) (*sdk.Result, error) {
	if keeper.ExistsClaim(ctx, msg.UserClaimAddress.String(), msg.UserClaimType) {
		ctx.Logger().Info("Claim already exists for user :", msg.UserClaimAddress.String())
		return nil, errors.Wrap(types.ErrInvalid, "Claim already exists for user")
	}
	newClaim := types.NewUserClaim(msg.UserClaimAddress, msg.UserClaimType, ctx.BlockTime().UTC())
	err := keeper.SetClaim(ctx, newClaim)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimCreated,
			sdk.NewAttribute(types.AttributeKeyClaimUser, newClaim.UserAddress.String()),
			sdk.NewAttribute(types.AttributeKeyClaimType, newClaim.UserClaimType.String()),
			sdk.NewAttribute(types.AttributeKeyClaimTime, newClaim.UserClaimTime.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
