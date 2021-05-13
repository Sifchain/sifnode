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
		//case MsgDistribution:
		//	return handleMsgCreateDistribution(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateDistribution is the top level function for calling all executors.
//func handleMsgCreateDistribution(ctx sdk.Context, keeper Keeper, msg MsgDistribution) (*sdk.Result, error) {
//	return nil, errors.New("Dispensation module is currently disabled")
//	// Verify if distribution already exists
//	err := keeper.VerifyDistribution(ctx, msg.DistributionName, msg.DistributionType)
//	if err != nil {
//		return nil, err
//	}
//	//Accumulate all Drops into the ModuleAccount
//	err = keeper.AccumulateDrops(ctx, msg.Input)
//	if err != nil {
//		return nil, err
//	}
//	//Create drops and Store Historical Data
//	err = keeper.CreateDrops(ctx, msg.Output, msg.DistributionName, msg.DistributionType)
//	if err != nil {
//		return nil, err
//	}
//	ctx.EventManager().EmitEvents(sdk.Events{
//		sdk.NewEvent(
//			types.EventTypeDistributionStarted,
//			sdk.NewAttribute(types.AttributeKeyFromModuleAccount, types.GetDistributionModuleAddress().String()),
//		),
//	})
//	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
//}

//func handleMsgCreateClaim(ctx sdk.Context, keeper Keeper, msg MsgCreateClaim) (*sdk.Result, error) {
//	if keeper.ExistsClaim(ctx, msg.Signer.String(), msg.UserClaimType) {
//		ctx.Logger().Info("Claim already exists for user :", msg.Signer.String())
//		return nil, errors.Wrap(types.ErrInvalid, "Claim already exists for user")
//	}
//	newClaim := types.NewUserClaim(msg.Signer, msg.UserClaimType, ctx.BlockTime().UTC())
//	err := keeper.SetClaim(ctx, newClaim)
//	if err != nil {
//		return nil, err
//	}
//	ctx.EventManager().EmitEvents(sdk.Events{
//		sdk.NewEvent(
//			types.EventTypeClaimCreated,
//			sdk.NewAttribute(types.AttributeKeyClaimUser, newClaim.UserAddress.String()),
//			sdk.NewAttribute(types.AttributeKeyClaimType, newClaim.UserClaimType.String()),
//			sdk.NewAttribute(types.AttributeKeyClaimTime, newClaim.UserClaimTime.String()),
//		),
//	})
//	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
//}
