package keeper

import (
	"context"
	"github.com/pkg/errors"

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
	err = srv.Keeper.CreateDrops(sdkCtx, msg.Output, msg.Distribution.DistributionName, msg.Distribution.DistributionType)
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

func (srv msgServer) CreateUserClaim(ctx context.Context,
	claim *types.MsgCreateUserClaim) (*types.MsgCreateClaimResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if srv.Keeper.ExistsClaim(sdkCtx, claim.Signer, claim.UserClaimType) {
		sdkCtx.Logger().Info("Claim already exists for user :", claim.Signer)
		return nil, errors.Wrap(types.ErrInvalid, "Claim already exists for user")
	}
	newClaim := types.NewUserClaim(claim.Signer, claim.UserClaimType, sdkCtx.BlockTime().UTC().String())
	err := srv.Keeper.SetClaim(sdkCtx, newClaim)
	if err != nil {
		return nil, err
	}
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimCreated,
			sdk.NewAttribute(types.AttributeKeyClaimUser, newClaim.UserAddress),
			sdk.NewAttribute(types.AttributeKeyClaimType, newClaim.UserClaimType.String()),
			sdk.NewAttribute(types.AttributeKeyClaimTime, newClaim.UserClaimTime),
		),
	})
	return &types.MsgCreateClaimResponse{}, nil
}
