package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

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
				//mtps := k.GetMTPsForPool(ctx, pool) // TODO define
				mtps := k.GetMTPsForAsset(ctx, pool.ExternalAsset.Symbol)
				for _, mtp := range mtps {
					h, err := k.UpdateMTPHealth(ctx, *mtp, *pool)
					if err != nil {
						continue // ?
					}
					mtp.MtpHealth = h
					_ = k.UpdateMTPInterestLiabilities(ctx, mtp, pool.InterestRate)
					_ = k.SetMTP(ctx, mtp)
					_, _ = k.ForceCloseLong(ctx, &types.MsgForceClose{Id: mtp.Id, MtpAddress: mtp.Address})
				}

				_ = k.UpdatePoolHealth(ctx, pool)
				_ = k.clpKeeper.SetPool(ctx, pool)
			}
		}

	}

}
