package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	//check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		currentHeight := ctx.BlockHeight()
		pools := k.ClpKeeper().GetPools(ctx)
		for _, pool := range pools {
			if k.IsPoolEnabled(ctx, pool.ExternalAsset.Symbol) {
				rate, err := k.InterestRateComputation(ctx, *pool)
				if err != nil {
					ctx.Logger().Error(err.Error())
					continue // ?
				}
				pool.InterestRate = rate
				pool.LastHeightInterestRateComputed = currentHeight
				_ = k.UpdatePoolHealth(ctx, pool)
				k.TrackSQBeginBlock(ctx, pool)
				_ = k.clpKeeper.SetPool(ctx, pool)
				mtps, _, _ := k.GetMTPsForPool(ctx, pool.ExternalAsset.Symbol, nil)
				for _, mtp := range mtps {
					BeginBlockerProcessMTP(ctx, k, mtp, pool)
				}
			}
		}
		k.EmitInterestRateComputation(ctx)
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
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error updating mtp health: %s", mtp.String())).Error())
		return
	}
	mtp.MtpHealth = h
	// compute interest
	interestPayment := CalcMTPInterestLiabilities(mtp, pool.InterestRate, 0, 0)

	k.HandleInterestPayment(ctx, interestPayment, mtp, pool)

	_ = k.SetMTP(ctx, mtp)
	repayAmount, err := k.ForceCloseLong(ctx, mtp, pool, false, true)

	if err == nil {
		// Emit event if position was closed
		k.EmitForceClose(ctx, mtp, repayAmount, "")
	} else if err != types.ErrMTPHealthy {
		ctx.Logger().Error(errors.Wrap(err, "error executing force close").Error())
	}

}
