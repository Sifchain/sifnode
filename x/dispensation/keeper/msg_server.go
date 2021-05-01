package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/dispensation/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the clp MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (srv msgServer) CreateDistribution(ctx context.Context,
	msg *types.MsgCreateDistribution) (*types.MsgCreateDistributionResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Verify if distribution already exists
	err := srv.Keeper.VerifyDistribution(sdkCtx, msg.Distribution.DistributionName, msg.Distribution.DistributionType)
	if err != nil {
		return nil, err
	}

	//Accumulate all Drops into the ModuleAccount
	err = srv.Keeper.AccumulateDrops(sdkCtx, msg.Input)
	if err != nil {
		return nil, err
	}

	//Create drops and Store Historical Data
	err = srv.Keeper.CreateDrops(sdkCtx, msg.Output, msg.Distribution.DistributionName)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDistributionStarted,
			sdk.NewAttribute(types.AttributeKeyFromModuleAccount, types.GetDistributionModuleAddress().String()),
		),
	})

	return &types.MsgCreateDistributionResponse{}, nil
}