package keeper

import (
	"context"
	"strings"

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

	mtp := types.MTP{
		Address:          msg.Signer,
		CollateralAsset:  msg.CollateralAsset,
		CollateralAmount: msg.CollateralAmount,
	}

	var err error
	var pool clptypes.Pool
	nativeAsset := types.GetSettlementAsset()

	if strings.EqualFold(msg.CollateralAsset, nativeAsset) {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.BorrowAsset)
		if err != nil {
			return nil, err
		}
	} else {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.CollateralAsset)
		if err != nil {
			return nil, err
		}
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
	pool := clptypes.Pool{}

	nativeAsset := types.GetSettlementAsset()
	if msg.CollateralAsset == nativeAsset {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.BorrowAsset)
	} else {
		pool, err = k.ClpKeeper().GetPool(ctx, msg.CollateralAsset)
	}

	err = k.TakeOutCustody(ctx, mtp, pool)

	repayAmount, err := k.CustodySwap(ctx, pool, mtp.CollateralAsset, mtp.CustodyAmount)

	interestRate, err := k.InterestRateComputation(ctx, pool)

	err = k.UpdateMTPInterestLiabilities(ctx, &mtp, interestRate)

	err = k.Repay(ctx, mtp, pool, repayAmount)

	return &types.MsgCloseLongResponse{}, nil
}
