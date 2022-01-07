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

	mtp := types.NewMTP(msg.Signer, msg.CollateralAsset, msg.CollateralAmount, msg.BorrowAsset)

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

	err = k.Borrow(ctx, msg.CollateralAsset, collateralAmount, borrowAmount, mtp, pool, leverage)
	if err != nil {
		return nil, err
	}

	err = k.UpdatePoolHealth(ctx, &pool)
	if err != nil {
		return nil, err
	}

	err = k.TakeInCustody(ctx, mtp, pool)
	if err != nil {
		return nil, err
	}

	return &types.MsgOpenLongResponse{}, nil
}

func (k msgServer) CloseLong(goCtx context.Context, msg *types.MsgCloseLong) (*types.MsgCloseLongResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mtp, err := k.GetMTP(ctx, msg.CollateralAsset, msg.BorrowAsset, msg.Signer)
	if err != nil {
		return nil, err
	}

	var pool clptypes.Pool

	nativeAsset := types.GetSettlementAsset()
	if strings.EqualFold(msg.CollateralAsset, nativeAsset) {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.BorrowAsset)
		if err != nil {
			return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, msg.BorrowAsset)
		}
	} else {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.CollateralAsset)
		if err != nil {
			return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, msg.CollateralAsset)
		}
	}

	err = k.TakeOutCustody(ctx, mtp, pool)
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

	err = k.Repay(ctx, mtp, pool, repayAmount)
	if err != nil {
		return nil, err
	}

	return &types.MsgCloseLongResponse{}, nil
}

func (k msgServer) ForceCloseLong(goCtx context.Context, msg *types.MsgCloseLong) (*types.MsgCloseLongResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mtp, err := k.GetMTP(ctx, msg.CollateralAsset, msg.BorrowAsset, msg.MtpAddress)
	if err != nil {
		return nil, err
	}

	var pool clptypes.Pool

	nativeAsset := types.GetSettlementAsset()
	if strings.EqualFold(msg.CollateralAsset, nativeAsset) {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.BorrowAsset)
		if err != nil {
			return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, msg.BorrowAsset)
		}
	} else {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.CollateralAsset)
		if err != nil {
			return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, msg.CollateralAsset)
		}
	}

	// check MTP health against threshold
	forceCloseThreshold := k.GetForceCloseThreshold(ctx)

	interestRate, err := k.InterestRateComputation(ctx, pool)
	if err != nil {
		return nil, err
	}

	err = k.UpdateMTPInterestLiabilities(ctx, &mtp, interestRate)
	if err != nil {
		return nil, err
	}

	mtpHealth, err := k.UpdateMTPHealth(ctx, mtp, pool)
	if err != nil {
		return nil, err
	}

	if mtpHealth.GT(forceCloseThreshold) {
		return nil, sdkerrors.Wrap(types.ErrMTPHealthy, msg.MtpAddress)
	}

	err = k.TakeOutCustody(ctx, mtp, pool)
	if err != nil {
		return nil, err
	}

	repayAmount, err := k.CustodySwap(ctx, pool, mtp.CollateralAsset, mtp.CustodyAmount)
	if err != nil {
		return nil, err
	}

	err = k.Repay(ctx, mtp, pool, repayAmount)
	if err != nil {
		return nil, err
	}

	return &types.MsgForceCloseLongResponse{}, nil
}
