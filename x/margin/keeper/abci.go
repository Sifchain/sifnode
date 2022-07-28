//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"strconv"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	//check if epoch has passed then execute
	currentHeight := ctx.BlockHeight()
	epochLength := k.GetEpochLength(ctx)
	if currentHeight%epochLength == 0 { // if epoch has passed
		events := sdk.EmptyEvents()
		pools := k.ClpKeeper().GetPools(ctx)
		for _, pool := range pools {
			if k.IsPoolEnabled(ctx, pool.ExternalAsset.Symbol) {
				rate, err := k.InterestRateComputation(ctx, *pool)
				if err != nil {
					ctx.Logger().Error(err.Error())
					continue // ?
				}
				pool.InterestRate = rate
				_ = k.UpdatePoolHealth(ctx, pool)
				_ = k.clpKeeper.SetPool(ctx, pool)
				mtps := k.GetMTPsForPool(ctx, pool.ExternalAsset.Symbol)
				for _, mtp := range mtps {
					BeginBlockerProcessMTP(ctx, k, mtp, pool)
				}
				events = events.AppendEvents(sdk.Events{
					sdk.NewEvent(
						types.EventInterestRateComputation,
						sdk.NewAttribute(clptypes.AttributeKeyPool, pool.ExternalAsset.Symbol),
						sdk.NewAttribute(types.AttributeKeyPoolInterestRate, rate.String()),
						sdk.NewAttribute(clptypes.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
					),
				})
			}
		}
		ctx.EventManager().EmitEvents(events)
	}

}

func BeginBlockerProcessMTP(ctx sdk.Context, k Keeper, mtp *types.MTP, pool *clptypes.Pool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.UpdateMTPHealth(ctx, *mtp, *pool)
	if err != nil {
		return
	}
	mtp.MtpHealth = h
	_ = k.UpdateMTPInterestLiabilities(ctx, mtp, pool.InterestRate)
	_ = k.SetMTP(ctx, mtp)
	_, repayAmount, err := k.ForceCloseLong(ctx, &types.MsgForceClose{Id: mtp.Id, MtpAddress: mtp.Address})
	if err == nil {
		// Emit event if position was closed
		k.EmitForceClose(ctx, mtp, repayAmount, "")
	}
}
