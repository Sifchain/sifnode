package keeper

import (
	"context"
	"fmt"

	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
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

	mtpToClose, err := k.GetMTP(ctx, msg.MtpAddress, msg.Id)
	if err != nil {
		return nil, err
	}

	var repayAmount sdk.Uint
	switch mtpToClose.Position {
	case types.Position_LONG:
		var pool clptypes.Pool

		nativeAsset := types.GetSettlementAsset()
		if types.StringCompare(mtpToClose.CollateralAsset, nativeAsset) {
			pool, err = k.ClpKeeper().GetPool(ctx, mtpToClose.CustodyAsset)
			if err != nil {
				return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, mtpToClose.CustodyAsset)
			}
		} else {
			pool, err = k.ClpKeeper().GetPool(ctx, mtpToClose.CollateralAsset)
			if err != nil {
				return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, mtpToClose.CollateralAsset)
			}
		}
		repayAmount, err = k.ForceCloseLong(ctx, &mtpToClose, &pool, true, msg.TakeMarginFund)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, mtpToClose.Position.String())
	}

	k.EmitAdminClose(ctx, &mtpToClose, repayAmount, msg.Signer)

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
