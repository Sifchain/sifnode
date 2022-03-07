package clp

import (
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	// get current block height
	currentHeight := ctx.BlockHeight()
	// get PMTP period params
	pmtpPeriodStartBlock := k.GetPmtpStartBlock(ctx)
	pmtpPeriodEndBlock := k.GetPmtpEndBlock(ctx)
	// Logic for start of PMTP period only
	if currentHeight == pmtpPeriodStartBlock {
		// get epoch length
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

		// Todo Set EpochCounter to Number to numEpochsInPolicyPeriod
		// Todo Set BlockPerEpochCounter Epoch Length
	}

	// Logic for every block in current Policy period
	if currentHeight >= pmtpPeriodStartBlock && currentHeight <= pmtpPeriodEndBlock {
		// get block rate
		pmtpPeriodBlockRate := k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate
		// compute running rate
		pmtpCurrentRunningRate := (sdk.NewDec(1).Add(pmtpPeriodBlockRate)).Power(uint64(currentHeight - pmtpPeriodStartBlock)).Sub(sdk.NewDec(1))
		// set running rate
		k.SetPmtpCurrentRunningRate(ctx, pmtpCurrentRunningRate)

		// Todo Decrement BlockPerEpochCounter
	}
	// Todo If BlockPerEpochCounter == 0 , Decrement EpochCounter
	// Todo If EpochCounter == 0 , Mark policy as Ended
	pools := k.GetPools(ctx)
	// compute swap prices for each pool
	for _, pool := range pools {

		normalizationFactor, adjustExternalToken := k.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
		// compute swap_price_native
		swapPriceNative, _, _, _, err := keeper.SwapOne(types.GetSettlementAsset(), sdk.OneUint(), *pool.ExternalAsset, *pool, normalizationFactor, adjustExternalToken)
		if err != nil {
			panic(err)
		}
		// compute swap_price_external
		swapPriceNative, _, _, _, err = keeper.SwapOne(*pool.ExternalAsset, sdk.OneUint(), types.GetSettlementAsset(), *pool, normalizationFactor, adjustExternalToken)
		if err != nil {
			panic(err)
		}
		pn := sdk.MustNewDecFromStr(swapPriceNative.String())
		pe := sdk.MustNewDecFromStr(swapPriceNative.String())
		pool.SwapPriceNative = &pn
		pool.SwapPriceExternal = &pe
		// set pool
		err = k.SetPool(ctx, pool)
		if err != nil {
			panic(err)
		}
	}
}
