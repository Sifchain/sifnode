package keeper

import (
	"context"
	"fmt"

	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AdminClose(goCtx context.Context, msg *types.MsgAdminClose) (*types.MsgAdminCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.AdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_MARGIN, signer) {
		return nil, sdkerrors.Wrap(admintypes.ErrPermissionDenied, fmt.Sprintf("signer not authorised: %s", msg.Signer))
	}

	mtpToClose, err := k.GetMTP(ctx, msg.Signer, msg.Id)
	if err != nil {
		return nil, err
	}

	var mtp *types.MTP
	var repayAmount sdk.Uint
	switch mtpToClose.Position {
	case types.Position_LONG:
		mtp, repayAmount, err = k.Keeper.ForceCloseLong(ctx, msg.Id, msg.MtpAddress, true, msg.TakeMarginFund)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, mtpToClose.Position.String())
	}

	k.EmitAdminClose(ctx, mtp, repayAmount, msg.Signer)

	return &types.MsgAdminCloseResponse{}, nil
}

func (k msgServer) AdminCloseAll(goCtx context.Context, msg *types.MsgAdminCloseAll) (*types.MsgAdminCloseAllResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.AdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_MARGIN, signer) {
		return nil, sdkerrors.Wrap(admintypes.ErrPermissionDenied, fmt.Sprintf("signer not authorised: %s", msg.Signer))
	}

	params := k.GetParams(ctx)
	params.SafetyFactor = sdk.NewDec(100)
	if !msg.TakeMarginFund {
		params.ForceCloseFundPercentage = sdk.ZeroDec()
	}

	k.SetParams(ctx, &params)
	k.EmitAdminCloseAll(ctx, msg.TakeMarginFund)

	return &types.MsgAdminCloseAllResponse{}, nil
}
