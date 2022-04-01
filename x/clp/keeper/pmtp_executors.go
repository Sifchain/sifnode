package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PolicyStart(ctx sdk.Context) {
	pmtpPeriodStartBlock := k.GetPmtpStartBlock(ctx)
	pmtpPeriodEndBlock := k.GetPmtpEndBlock(ctx)
	pmtpPeriodEpochLength := k.GetPmtpEpochLength(ctx)
	// get governance rate
	pmtpPeriodGovernanceRate := k.GetPmtpGovernaceRate(ctx)
	// compute length of policy period in blocks
	numBlocksInPolicyPeriod := pmtpPeriodEndBlock - pmtpPeriodStartBlock + 1
	// compute number of epochs in policy period
	numEpochsInPolicyPeriod := numBlocksInPolicyPeriod / pmtpPeriodEpochLength
	// compute pmtp period block rate
	pmtpPeriodBlockRate := (sdk.NewDec(1).Add(pmtpPeriodGovernanceRate)).Power(uint64((numEpochsInPolicyPeriod / numBlocksInPolicyPeriod))).Sub(sdk.NewDec(1))
	// set block rate
	// Block and Epoch calculations are done only on policy start
	k.SetPmtpBlockRate(ctx, pmtpPeriodBlockRate)
	k.SetPmtpEpoch(ctx, types.PmtpEpoch{
		EpochCounter: numEpochsInPolicyPeriod,
		BlockCounter: pmtpPeriodEpochLength,
	})
}

func (k Keeper) PolicyCalculations(ctx sdk.Context) sdk.Dec {
	currentHeight := ctx.BlockHeight()
	pmtpPeriodStartBlock := k.GetPmtpStartBlock(ctx)
	pmtpPeriodBlockRate := k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate
	// compute running rate
	pmtpCurrentRunningRate := (sdk.NewDec(1).Add(pmtpPeriodBlockRate)).Power(uint64(currentHeight - pmtpPeriodStartBlock)).Sub(sdk.NewDec(1))
	// set running rate
	k.SetPmtpCurrentRunningRate(ctx, pmtpCurrentRunningRate)
	k.DecrementBlockCounter(ctx)
	return pmtpCurrentRunningRate
}

func (k Keeper) PolicyRun(ctx sdk.Context, pmtpCurrentRunningRate sdk.Dec) error {
	pools := k.GetPools(ctx)
	// compute swap prices for each pool
	for _, pool := range pools {

		fmt.Printf("before pool.NativeAssetBalance: %v\n", pool.NativeAssetBalance)
		fmt.Printf("before pool.ExternalAssetBalance: %v\n", pool.ExternalAssetBalance)

		normalizationFactor, adjustExternalToken := k.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
		// compute swap_price_native
		swapPriceNative := CalcSwapPrice(types.GetSettlementAsset(), sdk.OneUint(), *pool.ExternalAsset, *pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)
		// compute swap_price_external
		swapPriceExternal := CalcSwapPrice(*pool.ExternalAsset, sdk.OneUint(), types.GetSettlementAsset(), *pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)

		fmt.Printf("swapPriceNative: %v\n", swapPriceNative)
		fmt.Printf("swapPriceExternal: %v\n", swapPriceExternal)

		pn := sdk.MustNewDecFromStr(swapPriceNative.String())
		pe := sdk.MustNewDecFromStr(swapPriceExternal.String())
		pool.SwapPriceNative = &pn
		pool.SwapPriceExternal = &pe
		// set pool
		err := k.SetPool(ctx, pool)
		if err != nil {
			return err
		}

		fmt.Printf("pool.NativeAssetBalance: %v\n", pool.NativeAssetBalance)
		fmt.Printf("pool.ExternalAssetBalance: %v\n", pool.ExternalAssetBalance)
		fmt.Printf("pool.SwapPriceNative: %v\n", *pool.SwapPriceNative)
		fmt.Printf("pool.SwapPriceExternal: %v\n", *pool.SwapPriceExternal)
	}
	return nil
}
