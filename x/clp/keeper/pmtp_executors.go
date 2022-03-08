package keeper

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PolicyStart(ctx sdk.Context) {
	pmtpPeriodStartBlock := k.GetPmtpStartBlock(ctx)
	pmtpPeriodEndBlock := k.GetPmtpEndBlock(ctx)
	k.Logger(ctx).Info(fmt.Sprintf("Starting new policy | Start Height : %d | End Height : %d", pmtpPeriodStartBlock, pmtpPeriodEndBlock))
	pmtpPeriodEpochLength := k.GetPmtpEpochLength(ctx)
	// get governance rate
	pmtpPeriodGovernanceRate := k.GetPmtpGovernaceRate(ctx)
	// compute length of policy period in blocks
	numBlocksInPolicyPeriod := pmtpPeriodStartBlock - pmtpPeriodEndBlock
	// compute number of epochs in policy period
	numEpochsInPolicyPeriod := numBlocksInPolicyPeriod / pmtpPeriodEpochLength
	// compute pmtp period block rate
	pmtpPeriodBlockRate := (sdk.NewDec(1).Add(pmtpPeriodGovernanceRate)).Power(uint64((numEpochsInPolicyPeriod / numBlocksInPolicyPeriod))).Sub(sdk.NewDec(1))
	// set block rate
	k.SetPmtpBlockRate(ctx, pmtpPeriodBlockRate)
	k.SetPmtpEpoch(ctx, types.PmtpEpoch{
		EpochCounter: numEpochsInPolicyPeriod,
		BlockCounter: pmtpPeriodEpochLength,
	})
}

func (k Keeper) PolicyCalculations(ctx sdk.Context) {
	currentHeight := ctx.BlockHeight()
	pmtpPeriodStartBlock := k.GetPmtpStartBlock(ctx)
	pmtpPeriodBlockRate := k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate
	// compute running rate
	pmtpCurrentRunningRate := (sdk.NewDec(1).Add(pmtpPeriodBlockRate)).Power(uint64(currentHeight - pmtpPeriodStartBlock)).Sub(sdk.NewDec(1))
	// set running rate
	k.SetPmtpCurrentRunningRate(ctx, pmtpCurrentRunningRate)
	k.DecrementBlockCounter(ctx)
}

func (k Keeper) PolicyRun(ctx sdk.Context) error {
	pools := k.GetPools(ctx)
	// compute swap prices for each pool
	for _, pool := range pools {

		normalizationFactor, adjustExternalToken := k.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
		// compute swap_price_native
		swapPriceNative, _, _, _, err := SwapOne(types.GetSettlementAsset(), sdk.OneUint(), *pool.ExternalAsset, *pool, normalizationFactor, adjustExternalToken)
		if err != nil {
			return err
		}
		// compute swap_price_external
		swapPriceNative, _, _, _, err = SwapOne(*pool.ExternalAsset, sdk.OneUint(), types.GetSettlementAsset(), *pool, normalizationFactor, adjustExternalToken)
		if err != nil {
			return err
		}
		pn := sdk.MustNewDecFromStr(swapPriceNative.String())
		pe := sdk.MustNewDecFromStr(swapPriceNative.String())
		pool.SwapPriceNative = &pn
		pool.SwapPriceExternal = &pe
		// set pool
		err = k.SetPool(ctx, pool)
		if err != nil {
			return err
		}
	}
	return nil
}
