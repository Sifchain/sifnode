package keeper

import (
	"strconv"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EmitForceClose(ctx sdk.Context, mtp *types.MTP, repayAmount sdk.Uint, closer string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventForceClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", mtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", mtp.CustodyAmount.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollateral.String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustody.String()),
		sdk.NewAttribute("interest_unpaid_collateral", mtp.InterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
		sdk.NewAttribute("closer", closer),
	))
}

func (k Keeper) EmitAdminClose(ctx sdk.Context, mtp *types.MTP, repayAmount sdk.Uint, closer string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAdminClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", mtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", mtp.CustodyAmount.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollateral.String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustody.String()),
		sdk.NewAttribute("interest_unpaid_collateral", mtp.InterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
		sdk.NewAttribute("closer", closer),
	))
}

func (k Keeper) EmitAdminCloseAll(ctx sdk.Context, takeMarginFund bool) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAdminCloseAll,
		sdk.NewAttribute("takeMarginFund", strconv.FormatBool(takeMarginFund)),
	))
}

func (k Keeper) EmitFundPayment(ctx sdk.Context, mtp *types.MTP, takeAmount sdk.Uint, takeAsset string, paymentType string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(paymentType,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("payment_amount", takeAmount.String()),
		sdk.NewAttribute("payment_asset", takeAsset),
	))
}

func (k Keeper) EmitBelowRemovalThreshold(ctx sdk.Context, pool *clptypes.Pool) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventBelowRemovalThreshold,
		sdk.NewAttribute("pool", pool.ExternalAsset.Symbol),
		sdk.NewAttribute("height", strconv.FormatInt(ctx.BlockHeight(), 10))))
}

func (k Keeper) EmitAboveRemovalThreshold(ctx sdk.Context, pool *clptypes.Pool) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAboveRemovalThreshold,
		sdk.NewAttribute("pool", pool.ExternalAsset.Symbol),
		sdk.NewAttribute("height", strconv.FormatInt(ctx.BlockHeight(), 10))))
}
