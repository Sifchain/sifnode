//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	//check if epoch has passed then execute
	currentHeight := ctx.BlockHeight()
	epochLength := k.GetEpochLength(ctx)
	if currentHeight%epochLength == 0 { // if epoch has passed
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
			}
		}

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
	_, err = k.ForceCloseLong(ctx, &types.MsgForceClose{Id: mtp.Id, MtpAddress: mtp.Address})
	if err == nil {
		// Emit event if position was closed
		k.EmitForceClose(ctx, mtp, "")
	}
}
