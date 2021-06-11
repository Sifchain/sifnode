package dispensation

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"strconv"
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
		case MsgRunDistribution:
			return handleMsgRunDistribution(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgRunDistribution(ctx sdk.Context, keeper Keeper, msg MsgRunDistribution) (*sdk.Result, error) {
	// Not checking whether the distribution exists or not .
	// We only need to run and execute distribution records
	// Distribute 10 drops for msg.DistributionName authorized to msg.DistributionRunner
	records, err := keeper.DistributeDrops(ctx, ctx.BlockHeight(), msg.DistributionName, msg.DistributionRunner, msg.DistributionType)
	if err != nil {
		return nil, err
	}

	var recordEvents []sdk.Event
	for i, record := range records {
		ev := sdk.NewEvent(
			types.EventTypeDistributionRecordsList+strconv.Itoa(i),
			sdk.NewAttribute(types.AttributeKeyDistributionRecordAddress, record.RecipientAddress.String()),
			sdk.NewAttribute(types.AttributeKeyDistributionRecordType, record.DistributionType.String()),
			sdk.NewAttribute(types.AttributeKeyDistributionRecordAmount, record.Coins.String()),
		)
		recordEvents = append(recordEvents, ev)
	}
	recordEvents = append(recordEvents, sdk.NewEvent(
		types.EventTypeDistributionRun,
		sdk.NewAttribute(types.AttributeKeyDistributionName, msg.DistributionName),
		sdk.NewAttribute(types.AttributeKeyDistributionRunner, msg.DistributionRunner.String()),
	))

	ctx.EventManager().EmitEvents(recordEvents)
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
	err = keeper.CreateDrops(ctx, msg.Output, distributionName, msg.DistributionType, msg.Runner)
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
