package keeper

import (
	"fmt"
	"math"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCurrentPmtpPolicy(ctx sdk.Context, params *types.PmtpParams) *types.PmtpPolicy {
	height := uint64(ctx.BlockHeight())
	for _, policy := range params.PmtpPolicies {
		if height >= uint64(policy.PmtpPeriodStartBlock) && height <= uint64(policy.PmtpPeriodEndBlock) {
			return policy
		}
	}
	return nil
}

func (k Keeper) PolicyStart(ctx sdk.Context, policy *types.PmtpPolicy) {
	pmtpPeriodStartBlock := policy.PmtpPeriodStartBlock
	pmtpPeriodEndBlock := policy.PmtpPeriodEndBlock
	pmtpPeriodEpochLength := policy.PmtpPeriodEpochLength
	// get governance rate
	pmtpPeriodGovernanceRate := policy.PmtpPeriodGovernanceRate
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

func (k Keeper) PolicyCalculations(ctx sdk.Context, policy *types.PmtpPolicy) sdk.Dec {
	currentHeight := ctx.BlockHeight()
	pmtpPeriodStartBlock := policy.PmtpPeriodStartBlock
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

func (k Keeper) PolicyRun(ctx sdk.Context, pmtpCurrentRunningRate sdk.Dec) error {
	pools := k.GetPools(ctx)
	// compute swap prices for each pool
	for _, pool := range pools {
		normalizationFactor, adjustExternalToken := k.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
		// compute swap_price_native
		swapPriceNative := CalcSwapPrice(types.GetSettlementAsset(), sdk.OneUint(), *pool.ExternalAsset, *pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)
		// compute swap_price_external
		swapPriceExternal := CalcSwapPrice(*pool.ExternalAsset, sdk.OneUint(), types.GetSettlementAsset(), *pool, normalizationFactor, adjustExternalToken, pmtpCurrentRunningRate)

		pn := sdk.MustNewDecFromStr(swapPriceNative.String())
		pe := sdk.MustNewDecFromStr(swapPriceExternal.String())
		pool.SwapPriceNative = &pn
		pool.SwapPriceExternal = &pe
		// set pool
		err := k.SetPool(ctx, pool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) IsInsidePmtpWindow(ctx sdk.Context) bool {
	params := k.GetPmtpParams(ctx)
	for _, policy := range params.PmtpPolicies {
		if ctx.BlockHeight() <= policy.PmtpPeriodEndBlock && ctx.BlockHeight() >= policy.PmtpPeriodStartBlock {
			return true
		}
	}
	return false
}
