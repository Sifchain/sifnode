package keeper

import (
	"context"
	"fmt"
	"strconv"

	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.IsWhitelistingEnabled(ctx) && !k.IsWhitelisted(ctx, msg.Signer) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorised, "unauthorised")
	}

	if k.GetOpenMTPCount(ctx) >= k.GetMaxOpenPositions(ctx) {
		return nil, sdkerrors.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}

	var externalAsset string
	nativeAsset := types.GetSettlementAsset()

	if types.StringCompare(msg.CollateralAsset, nativeAsset) {
		externalAsset = msg.BorrowAsset
	} else {
		externalAsset = msg.CollateralAsset
	}

	pool, err := k.ClpKeeper().GetPool(ctx, externalAsset)
	if err != nil {
		return nil, sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, externalAsset)
	}

	if !k.IsPoolEnabled(ctx, externalAsset) || k.IsPoolClosed(ctx, externalAsset) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, externalAsset)
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, "pool health too low to open new positions")
	}

	var mtp *types.MTP

	switch msg.Position {
	case types.Position_LONG:
		mtp, err = k.OpenLong(ctx, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", mtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", mtp.CustodyAmount.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollateral.String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustody.String()),
		sdk.NewAttribute("interest_unpaid_collateral", mtp.InterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	))

	return &types.MsgOpenResponse{}, nil
}

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	mtp, err := k.GetMTP(ctx, msg.Signer, msg.Id)
	if err != nil {
		return nil, err
	}

	var closedMtp *types.MTP
	var repayAmount sdk.Uint
	switch mtp.Position {
	case types.Position_LONG:
		closedMtp, repayAmount, err = k.CloseLong(ctx, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, mtp.Position.String())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(closedMtp.Id), 10)),
		sdk.NewAttribute("position", closedMtp.Position.String()),
		sdk.NewAttribute("address", closedMtp.Address),
		sdk.NewAttribute("collateral_asset", closedMtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", closedMtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", closedMtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", closedMtp.CustodyAmount.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", closedMtp.Leverage.String()),
		sdk.NewAttribute("liabilities", closedMtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollateral.String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustody.String()),
		sdk.NewAttribute("interest_unpaid_collateral", closedMtp.InterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", closedMtp.MtpHealth.String()),
	))

	return &types.MsgCloseResponse{}, nil
}

func (k msgServer) OpenLong(ctx sdk.Context, msg *types.MsgOpen) (*types.MTP, error) {
	maxLeverage := k.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	eta := leverage.Sub(sdk.OneDec())

	collateralAmount := msg.CollateralAmount

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())

	mtp := types.NewMTP(msg.Signer, msg.CollateralAsset, msg.BorrowAsset, msg.Position, leverage)

	var externalAsset string
	nativeAsset := types.GetSettlementAsset()

	if !k.IsRowanCollateralEnabled(ctx) && types.StringCompare(msg.CollateralAsset, nativeAsset) {
		return nil, sdkerrors.Wrap(types.ErrRowanAsCollateralNotAllowed, nativeAsset)
	}

	if types.StringCompare(msg.CollateralAsset, nativeAsset) {
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

	leveragedAmountDec := collateralAmountDec.Mul(leverage)

	leveragedAmount := sdk.NewUintFromBigInt(leveragedAmountDec.TruncateInt().BigInt())

	ctx.Logger().Info(fmt.Sprintf("leveragedAmount: %s", leveragedAmount.String()))

	if types.StringCompare(msg.CollateralAsset, nativeAsset) {
		if leveragedAmount.GT(pool.NativeAssetBalance) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
		}
	} else {
		if leveragedAmount.GT(pool.ExternalAssetBalance) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
		}
	}

	// check if liabilities large enough for interest payments
	err = k.CheckMinLiabilities(ctx, collateralAmount, eta, pool, msg.BorrowAsset)
	if err != nil {
		return nil, err
	}

	custodyAmount, err := k.CLPSwap(ctx, leveragedAmount, msg.BorrowAsset, pool)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Info(fmt.Sprintf("custodyAmount: %s", custodyAmount.String()))

	if types.StringCompare(msg.CollateralAsset, nativeAsset) {
		if custodyAmount.GT(pool.ExternalAssetBalance) {
			return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
		}
	} else {
		if custodyAmount.GT(pool.NativeAssetBalance) {
			return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
		}
	}

	err = k.Borrow(ctx, msg.CollateralAsset, collateralAmount, custodyAmount, mtp, &pool, eta)
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

	safetyFactor := k.GetSafetyFactor(ctx)

	lr, err := k.UpdateMTPHealth(ctx, *mtp, pool)

	if err != nil {
		return nil, err
	}

	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// res, stop := k.ClpKeeper().SingleExternalBalanceModuleAccountCheck(externalAsset)(ctx)
	// if stop {
	// 	return nil, sdkerrors.Wrap(clptypes.ErrBalanceModuleAccountCheck, res)
	// }

	return mtp, nil
}

func (k msgServer) CloseLong(ctx sdk.Context, msg *types.MsgClose) (*types.MTP, sdk.Uint, error) {
	mtp, err := k.GetMTP(ctx, msg.Signer, msg.Id)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	var pool clptypes.Pool

	nativeAsset := types.GetSettlementAsset()
	if types.StringCompare(mtp.CollateralAsset, nativeAsset) {
		pool, err = k.ClpKeeper().GetPool(ctx, mtp.CustodyAsset)
		if err != nil {
			return nil, sdk.ZeroUint(), sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, mtp.CustodyAsset)
		}
	} else {
		pool, err = k.ClpKeeper().GetPool(ctx, mtp.CollateralAsset)
		if err != nil {
			return nil, sdk.ZeroUint(), sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, mtp.CollateralAsset)
		}
	}

	epochLength := k.GetEpochLength(ctx)
	epochPosition := GetEpochPosition(ctx, epochLength)
	if epochPosition > 0 {
		interestPayment := CalcMTPInterestLiabilities(&mtp, pool.InterestRate, epochPosition, epochLength)

		finalInterestPayment := k.HandleInterestPayment(ctx, interestPayment, &mtp, &pool)

		if types.StringCompare(mtp.CollateralAsset, nativeAsset) { // custody is external, payment is custody
			pool.BlockInterestExternal = pool.BlockInterestExternal.Add(finalInterestPayment)
		} else { // custody is native, payment is custody
			pool.BlockInterestNative = pool.BlockInterestNative.Add(finalInterestPayment)
		}

		mtp.MtpHealth, err = k.UpdateMTPHealth(ctx, mtp, pool)
		if err != nil {
			return nil, sdk.ZeroUint(), err
		}
	}

	err = k.TakeOutCustody(ctx, mtp, &pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}
	repayAmount, err := k.CLPSwap(ctx, mtp.CustodyAmount, mtp.CollateralAsset, pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	err = k.Repay(ctx, &mtp, &pool, repayAmount, false)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	// if types.StringCompare(mtp.CollateralAsset, nativeAsset) {
	// 	res, stop := k.ClpKeeper().SingleExternalBalanceModuleAccountCheck(mtp.CustodyAsset)(ctx)
	// 	if stop {
	// 		return nil, sdk.ZeroUint(), sdkerrors.Wrap(clptypes.ErrBalanceModuleAccountCheck, res)
	// 	}
	// } else {
	// 	res, stop := k.ClpKeeper().SingleExternalBalanceModuleAccountCheck(mtp.CollateralAsset)(ctx)
	// 	if stop {
	// 		return nil, sdk.ZeroUint(), sdkerrors.Wrap(clptypes.ErrBalanceModuleAccountCheck, res)
	// 	}
	// }

	return &mtp, repayAmount, nil
}

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.AdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_MARGIN, signer) {
		return nil, sdkerrors.Wrap(admintypes.ErrPermissionDenied, fmt.Sprintf("signer not authorised: %s", msg.Signer))
	}

	params := k.GetParams(ctx)
	msg.Params.Pools = params.Pools
	k.SetParams(ctx, msg.Params)

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventMarginUpdateParams,
		sdk.NewAttribute(types.AttributeKeyMarginParams, params.String()),
		sdk.NewAttribute(clptypes.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	))

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) UpdatePools(goCtx context.Context, msg *types.MsgUpdatePools) (*types.MsgUpdatePoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.AdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_MARGIN, signer) {
		return nil, sdkerrors.Wrap(admintypes.ErrPermissionDenied, fmt.Sprintf("signer not authorised: %s", msg.Signer))
	}

	params := k.GetParams(ctx)
	params.Pools = msg.Pools
	params.ClosedPools = msg.ClosedPools
	k.SetParams(ctx, &params)

	return &types.MsgUpdatePoolsResponse{}, nil
}

func (k msgServer) UpdateRowanCollateral(goCtx context.Context, msg *types.MsgUpdateRowanCollateral) (*types.MsgUpdateRowanCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.AdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_MARGIN, signer) {
		return nil, sdkerrors.Wrap(admintypes.ErrPermissionDenied, fmt.Sprintf("signer not authorised: %s", msg.Signer))
	}

	params := k.GetParams(ctx)
	params.RowanCollateralEnabled = msg.RowanCollateralEnabled
	k.SetParams(ctx, &params)

	return &types.MsgUpdateRowanCollateralResponse{}, nil
}

func (k msgServer) Whitelist(goCtx context.Context, msg *types.MsgWhitelist) (*types.MsgWhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.AdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_MARGIN, signer) {
		return nil, sdkerrors.Wrap(admintypes.ErrPermissionDenied, fmt.Sprintf("signer not authorised: %s", msg.Signer))
	}

	k.WhitelistAddress(ctx, msg.WhitelistedAddress)

	return &types.MsgWhitelistResponse{}, nil
}

func (k msgServer) Dewhitelist(goCtx context.Context, msg *types.MsgDewhitelist) (*types.MsgDewhitelistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.AdminKeeper().IsAdminAccount(ctx, admintypes.AdminType_MARGIN, signer) {
		return nil, sdkerrors.Wrap(admintypes.ErrPermissionDenied, fmt.Sprintf("signer not authorised: %s", msg.Signer))
	}

	k.DewhitelistAddress(ctx, msg.WhitelistedAddress)

	return &types.MsgDewhitelistResponse{}, nil
}

// ForceClose is deprecated replaced by AdminClose
func (k msgServer) ForceClose(goCtx context.Context, msg *types.MsgForceClose) (*types.MsgForceCloseResponse, error) {
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
		repayAmount, err = k.Keeper.ForceCloseLong(ctx, &mtpToClose, &pool, true, false)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, mtpToClose.Position.String())
	}

	k.EmitAdminClose(ctx, &mtpToClose, repayAmount, msg.Signer)

	// if types.StringCompare(mtpToClose.CollateralAsset, types.GetSettlementAsset()) {
	// 	res, stop := k.ClpKeeper().SingleExternalBalanceModuleAccountCheck(mtpToClose.CustodyAsset)(ctx)
	// 	if stop {
	// 		return nil, sdkerrors.Wrap(clptypes.ErrBalanceModuleAccountCheck, res)
	// 	}
	// } else {
	// 	res, stop := k.ClpKeeper().SingleExternalBalanceModuleAccountCheck(mtpToClose.CollateralAsset)(ctx)
	// 	if stop {
	// 		return nil, sdkerrors.Wrap(clptypes.ErrBalanceModuleAccountCheck, res)
	// 	}
	// }

	return &types.MsgForceCloseResponse{}, nil
}
