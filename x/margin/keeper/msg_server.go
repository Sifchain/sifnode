package keeper

import (
	"context"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	types.Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(k types.Keeper) types.MsgServer {
	return msgServer{
		k,
	}
}

func (k msgServer) OpenLong(goCtx context.Context, msg *types.MsgOpenLong) (*types.MsgOpenLongResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	leverage := k.GetLeverageParam(ctx)

	collateralAmount := msg.CollateralAmount

	mtp := types.NewMTP(msg.Signer, msg.CollateralAsset, msg.BorrowAsset)

	var externalAsset string
	nativeAsset := types.GetSettlementAsset()

	if strings.EqualFold(msg.CollateralAsset, nativeAsset) {
		externalAsset = msg.BorrowAsset
	} else {
		externalAsset = msg.CollateralAsset
	}

	pool, err := k.ClpKeeper().GetPool(ctx, externalAsset)
	if err != nil {
		return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, externalAsset)
	}

	if !k.IsPoolEnabled(ctx, externalAsset) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, externalAsset)
	}

	leveragedAmount := collateralAmount.Mul(sdk.NewUint(1).Add(leverage))

	borrowAmount, err := k.CustodySwap(ctx, pool, msg.BorrowAsset, leveragedAmount)
	if err != nil {
		return nil, err
	}

	err = k.Borrow(ctx, msg.CollateralAsset, collateralAmount, borrowAmount, mtp, &pool, leverage)
	if err != nil {
		return nil, err
	}

	err = k.UpdatePoolHealth(ctx, &pool)
	if err != nil {
		return nil, err
	}

	err = k.TakeInCustody(ctx, *mtp, &pool)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("position", "long"),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", mtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", mtp.CustodyAmount.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities_p", mtp.LiabilitiesP.String()),
		sdk.NewAttribute("liabilities_i", mtp.LiabilitiesI.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	))

	return &types.MsgOpenLongResponse{}, nil
}

func (k msgServer) CloseLong(goCtx context.Context, msg *types.MsgCloseLong) (*types.MsgCloseLongResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mtp, err := k.GetMTP(ctx, msg.Signer, msg.Id)
	if err != nil {
		return nil, err
	}

	var pool clptypes.Pool

	nativeAsset := types.GetSettlementAsset()
	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) {
		pool, err = k.ClpKeeper().GetPool(ctx, mtp.CustodyAsset)
		if err != nil {
			return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, mtp.CustodyAsset)
		}
	} else {
		pool, err = k.ClpKeeper().GetPool(ctx, mtp.CollateralAsset)
		if err != nil {
			return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, mtp.CollateralAsset)
		}
	}

	err = k.TakeOutCustody(ctx, mtp, &pool)
	if err != nil {
		return nil, err
	}

	repayAmount, err := k.CustodySwap(ctx, pool, mtp.CollateralAsset, mtp.CustodyAmount)
	if err != nil {
		return nil, err
	}

	interestRate, err := k.InterestRateComputation(ctx, pool)
	if err != nil {
		return nil, err
	}

	err = k.UpdateMTPInterestLiabilities(ctx, &mtp, interestRate)
	if err != nil {
		return nil, err
	}

	err = k.Repay(ctx, &mtp, pool, repayAmount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("position", "long"),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", mtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", mtp.CustodyAmount.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities_p", mtp.LiabilitiesP.String()),
		sdk.NewAttribute("liabilities_i", mtp.LiabilitiesI.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	))

	return &types.MsgCloseLongResponse{}, nil
}
