package keeper

import (
	"fmt"
	"math"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) PolicyStart(ctx sdk.Context) {
	pmtpParams := k.GetPmtpParams(ctx)
	pmtpPeriodStartBlock := pmtpParams.PmtpPeriodStartBlock
	pmtpPeriodEndBlock := pmtpParams.PmtpPeriodEndBlock
	pmtpPeriodEpochLength := pmtpParams.PmtpPeriodEpochLength
	// get governance rate
	pmtpPeriodGovernanceRate := pmtpParams.PmtpPeriodGovernanceRate
	// compute length of policy period in blocks
	numBlocksInPolicyPeriod := pmtpPeriodEndBlock - pmtpPeriodStartBlock + 1
	// compute number of epochs in policy period
	numEpochsInPolicyPeriod := numBlocksInPolicyPeriod / pmtpPeriodEpochLength
	// compute pmtp period block rate
	//pmtpPeriodBlockRate = (1 + pmtpPeriodGovernanceRate).Pow(numEpochsInPolicyPeriod / numBlocksInPolicyPeriod) - 1
	// set block rate
	base := sdk.NewDec(1).Add(pmtpPeriodGovernanceRate).MustFloat64()
	pow := float64(numEpochsInPolicyPeriod) / float64(numBlocksInPolicyPeriod)
	firstSection := math.Pow(base, pow)
	pmtpPeriodBlockRate := firstSection - 1
	decBlockrate, err := sdk.NewDecFromStr(fmt.Sprintf("%.18f", pmtpPeriodBlockRate))
	if err != nil {
		panic(err)
	}
	// set block rate
	// Block and Epoch calculations are done only on policy start
	k.SetPmtpBlockRate(ctx, decBlockrate)
	k.SetPmtpEpoch(ctx, types.PmtpEpoch{
		EpochCounter: numEpochsInPolicyPeriod,
		BlockCounter: pmtpPeriodEpochLength,
	})
}

func (k Keeper) PolicyCalculations(ctx sdk.Context) sdk.Dec {
	currentHeight := ctx.BlockHeight()
	pmtpPeriodStartBlock := k.GetPmtpParams(ctx).PmtpPeriodStartBlock
	rateParams := k.GetPmtpRateParams(ctx)
	pmtpPeriodBlockRate := rateParams.PmtpPeriodBlockRate
	pmtpInterPolicyRate := rateParams.PmtpInterPolicyRate
	// compute running rate
	pmtpCurrentRunningRate := (sdk.NewDec(1).Add(pmtpPeriodBlockRate)).Power(uint64(currentHeight - pmtpPeriodStartBlock + 1)).Sub(sdk.NewDec(1))
	pmtpCurrentRunningRate = pmtpCurrentRunningRate.Add(pmtpInterPolicyRate)
	// set running rate
	k.SetPmtpCurrentRunningRate(ctx, pmtpCurrentRunningRate)
	return pmtpCurrentRunningRate
}

// NOTE: the code in this method must not panic otherwise the chain will halt
// see: https://bytemeta.vip/repo/osmosis-labs/osmosis/issues/1305
func (k Keeper) PolicyRun(ctx sdk.Context, pmtpCurrentRunningRate sdk.Dec) error {
	pools := k.GetPools(ctx)

	// NOTE: if an error occurs in this loop we must continue to update the remaining pools
	for _, pool := range pools {
		decimalsExternal, err := k.GetAssetDecimals(ctx, *pool.ExternalAsset)
		if err != nil {
			// Conditions for an error:
			// 1. Asset has decimals value outside of 0-255 range
			// 2. Asset cannot be found in the token registry
			// In either case we should log the error then continue updating the remaining pools
			ctx.Logger().Error(errors.Wrap(err, "error calculating pool spot price").Error())
			continue
		}

		spotPriceNative, err := CalcSpotPriceNative(pool, decimalsExternal, pmtpCurrentRunningRate)
		if err != nil {
			// Error occurs if native asset pool depth is zero or result overflows
			spotPriceNative = sdk.ZeroDec()
		}
		spotPriceExternal, err := CalcSpotPriceExternal(pool, decimalsExternal, pmtpCurrentRunningRate)
		if err != nil {
			// Error occurs if external asset pool depth is zero or result overflows
			spotPriceExternal = sdk.ZeroDec()
		}

		pool.SwapPriceNative = &spotPriceNative
		pool.SwapPriceExternal = &spotPriceExternal

		// ignore error since it will always be nil
		_ = k.SetPool(ctx, pool)
	}
	return nil
}

func (k Keeper) IsInsidePmtpWindow(ctx sdk.Context) bool {
	params := k.GetPmtpParams(ctx)
	return ctx.BlockHeight() <= params.PmtpPeriodEndBlock && ctx.BlockHeight() >= params.PmtpPeriodStartBlock
}
