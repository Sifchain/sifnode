//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

	if k.GetOpenMTPCount(ctx) >= k.GetMaxOpenPositions(ctx) {
		return nil, sdkerrors.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}

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
		sdk.NewAttribute("liabilities_p", mtp.LiabilitiesP.String()),
		sdk.NewAttribute("liabilities_i", mtp.LiabilitiesI.String()),
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
		sdk.NewAttribute("liabilities_p", closedMtp.LiabilitiesP.String()),
		sdk.NewAttribute("liabilities_i", closedMtp.LiabilitiesI.String()),
		sdk.NewAttribute("health", closedMtp.MtpHealth.String()),
	))

	return &types.MsgCloseResponse{}, nil
}

func (k msgServer) ForceClose(goCtx context.Context, msg *types.MsgForceClose) (*types.MsgForceCloseResponse, error) {
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
		mtp, repayAmount, err = k.Keeper.ForceCloseLong(ctx, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, mtpToClose.Position.String())
	}

	k.EmitForceClose(ctx, mtp, repayAmount, msg.Signer)

	return &types.MsgForceCloseResponse{}, nil
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

	leveragedAmountDec := collateralAmountDec.Mul(leverage)

	leveragedAmount := sdk.NewUintFromBigInt(leveragedAmountDec.TruncateInt().BigInt())

	ctx.Logger().Info(fmt.Sprintf("leveragedAmount: %s", leveragedAmount.String()))

	custodyAmount, err := k.CLPSwap(ctx, leveragedAmount, msg.BorrowAsset, pool)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Info(fmt.Sprintf("custodyAmount: %s", custodyAmount.String()))

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

	return mtp, nil
}

func (k msgServer) CloseLong(ctx sdk.Context, msg *types.MsgClose) (*types.MTP, sdk.Uint, error) {
	mtp, err := k.GetMTP(ctx, msg.Signer, msg.Id)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	var pool clptypes.Pool

	nativeAsset := types.GetSettlementAsset()
	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) {
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

	err = k.TakeOutCustody(ctx, mtp, &pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}
	repayAmount, err := k.CLPSwap(ctx, mtp.CustodyAmount, mtp.CollateralAsset, pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	interestRate, err := k.InterestRateComputation(ctx, pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	err = k.UpdateMTPInterestLiabilities(ctx, &mtp, interestRate)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	err = k.Repay(ctx, &mtp, pool, repayAmount, false)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	return &mtp, repayAmount, nil
}

func (k Keeper) ForceCloseLong(ctx sdk.Context, msg *types.MsgForceClose) (*types.MTP, sdk.Uint, error) {
	mtp, err := k.GetMTP(ctx, msg.MtpAddress, msg.Id)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	var pool clptypes.Pool

	nativeAsset := types.GetSettlementAsset()
	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) {
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

	// check MTP health against threshold
	forceCloseThreshold := k.GetForceCloseThreshold(ctx)

	interestRate, err := k.InterestRateComputation(ctx, pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	err = k.UpdateMTPInterestLiabilities(ctx, &mtp, interestRate)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	mtpHealth, err := k.UpdateMTPHealth(ctx, mtp, pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	if mtpHealth.GT(forceCloseThreshold) {
		return nil, sdk.ZeroUint(), sdkerrors.Wrap(types.ErrMTPHealthy, msg.MtpAddress)
	}

	err = k.TakeOutCustody(ctx, mtp, &pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	repayAmount, err := k.CLPSwap(ctx, mtp.CustodyAmount, mtp.CollateralAsset, pool)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

	err = k.Repay(ctx, &mtp, pool, repayAmount, true)
	if err != nil {
		return nil, sdk.ZeroUint(), err
	}

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
	k.SetParams(ctx, &params)

	return &types.MsgUpdatePoolsResponse{}, nil
}
