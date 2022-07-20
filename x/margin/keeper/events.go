//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"strconv"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EmitForceClose(ctx sdk.Context, mtp *types.MTP, closer string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventForceClose,
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
		sdk.NewAttribute("closer", closer),
	))
}

func (k Keeper) EmitRepayInsuranceFund(ctx sdk.Context, mtp *types.MTP, takeAmount sdk.Uint) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventRepayInsuranceFund,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("takeAmount", takeAmount.String()),
	))
}
