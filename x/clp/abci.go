package clp

import (
	"fmt"
	"strconv"
	"time"

	kpr "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, keeper kpr.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if keeper.IsDistributionBlock(ctx) {
		keeper.ProviderDistributionPolicyRun(ctx)
	}

	params := keeper.GetRewardsParams(ctx)
	pools := keeper.GetPools(ctx)
	currentPeriod := keeper.GetCurrentRewardPeriod(ctx, params)
	if currentPeriod != nil && !currentPeriod.RewardPeriodAllocation.IsZero() {

		isDistributionBlock := kpr.IsDistributionBlockPure(ctx.BlockHeight(), currentPeriod.RewardPeriodStartBlock, currentPeriod.RewardPeriodMod)

		currentBlockDistribution := kpr.CalcBlockDistribution(currentPeriod)
		blockDistributionAccu := keeper.GetBlockDistributionAccu(ctx)
		blockDistribution := blockDistributionAccu.Add(currentBlockDistribution)
		if isDistributionBlock {
			err := keeper.DistributeDepthRewards(ctx, blockDistribution, currentPeriod, pools)
			keeper.SetBlockDistributionAccu(ctx, sdk.ZeroUint())
			if err != nil {
				keeper.Logger(ctx).Error(fmt.Sprintf("Rewards policy run error %s", err.Error()))
			}
		} else {
			keeper.SetBlockDistributionAccu(ctx, blockDistribution)
		}
	}

	return []abci.ValidatorUpdate{}
}

func BeginBlocker(ctx sdk.Context, k kpr.Keeper) {
	if ctx.BlockHeight() == 8654226 {
		fixAtomPool(ctx, k)
	}
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	MeasureBlockTime(ctx, k)

	// get current block height
	currentHeight := ctx.BlockHeight()

	/*
		Liquidity protection current rowan liquidity threshold update
	*/
	liquidityProtectionParams := k.GetLiquidityProtectionParams(ctx)
	maxRowanLiquidityThreshold := liquidityProtectionParams.MaxRowanLiquidityThreshold
	maxRowanLiquidityThresholdAsset := liquidityProtectionParams.MaxRowanLiquidityThresholdAsset
	if liquidityProtectionParams.IsActive {
		currentRowanLiquidityThreshold := k.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold
		// Validation check ensures that Epoch length =/= zero
		replenishmentAmount := maxRowanLiquidityThreshold.QuoUint64(liquidityProtectionParams.EpochLength)

		// This is equivalent to:
		//    proposedThreshold := currentRowanLiquidityThreshold.Add(replenishmentAmount)
		//    currentRowanLiquidityThreshold = sdk.MinUint(proposedThreshold, maxRowanLiquidityThreshold)
		// except it prevents any overflows when adding the replenishmentAmount
		if maxRowanLiquidityThreshold.Sub(currentRowanLiquidityThreshold).LT(replenishmentAmount) {
			currentRowanLiquidityThreshold = maxRowanLiquidityThreshold
		} else {
			currentRowanLiquidityThreshold = currentRowanLiquidityThreshold.Add(replenishmentAmount)
		}

		k.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, currentRowanLiquidityThreshold)
		k.Logger(ctx).Info(fmt.Sprintf("Liquidity Protection | maxRowanLiquidityThreshold: %s | asset: %s | currentRowanLiquidityThreshold: %s | maxPerBlock: %s", maxRowanLiquidityThreshold, maxRowanLiquidityThresholdAsset, k.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold, replenishmentAmount))
	}

	// get PMTP period params
	pmtpPeriodStartBlock := k.GetPmtpParams(ctx).PmtpPeriodStartBlock
	pmtpPeriodEndBlock := k.GetPmtpParams(ctx).PmtpPeriodEndBlock
	// Start Policy
	if currentHeight == pmtpPeriodStartBlock &&
		k.GetPmtpEpoch(ctx).EpochCounter == 0 &&
		k.GetPmtpEpoch(ctx).BlockCounter == 0 {
		k.PolicyStart(ctx)
		_ = ctx.EventManager().EmitTypedEvent(&types.EventPolicy{
			EventType:            "policy_start",
			PmtpPeriodStartBlock: strconv.Itoa(int(pmtpPeriodStartBlock)),
			PmtpPeriodEndBlock:   strconv.Itoa(int(pmtpPeriodEndBlock)),
		})
		k.Logger(ctx).Info(fmt.Sprintf("Starting new policy | Start Height : %d | End Height : %d", pmtpPeriodStartBlock, pmtpPeriodEndBlock))
	}
	// default to current pmtp current running rate value
	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate

	/*
		Epoch counters are used to keep track of policy execution
		EpochCounter tracks the number of Epochs left (Decrementing counter )
		BlockCounter tracks the number of Blocks left (Decrementing counter ) in an Epoch
	*/

	// Manage Block Counter
	if currentHeight >= pmtpPeriodStartBlock &&
		currentHeight <= pmtpPeriodEndBlock &&
		// TODO : Check this condition might not be needed
		k.GetPmtpEpoch(ctx).EpochCounter > 0 {
		// Calculate R running for policy params
		pmtpCurrentRunningRate = k.PolicyCalculations(ctx)
		k.DecrementPmtpBlockCounter(ctx)
	}
	// Manage Epoch Counter
	if k.GetPmtpEpoch(ctx).BlockCounter == 0 &&
		currentHeight < pmtpPeriodEndBlock &&
		currentHeight >= pmtpPeriodStartBlock {
		k.DecrementPmtpEpochCounter(ctx)
		k.SetPmtpBlockCounter(ctx, k.GetPmtpParams(ctx).PmtpPeriodEpochLength)
	}

	if currentHeight == pmtpPeriodEndBlock {
		// Setting it to zero to check when we start policy
		k.SetPmtpEpoch(ctx, types.PmtpEpoch{
			EpochCounter: 0,
			BlockCounter: 0,
		})
		// Set inter policy rate to running rate
		k.SetPmtpInterPolicyRate(ctx, pmtpCurrentRunningRate)
		_ = ctx.EventManager().EmitTypedEvent(&types.EventPolicy{
			EventType:            "policy_end",
			PmtpPeriodStartBlock: strconv.Itoa(int(pmtpPeriodStartBlock)),
			PmtpPeriodEndBlock:   strconv.Itoa(int(pmtpPeriodEndBlock)),
		})
		k.Logger(ctx).Info(fmt.Sprintf("Ending Policy | Start Height : %d | End Height : %d", pmtpPeriodStartBlock, pmtpPeriodEndBlock))
	}

	err := k.PolicyRun(ctx, pmtpCurrentRunningRate)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error in running policy | Error Message : %s ", err.Error()))
	}

}

/*
fixAtomPool :
*/
func fixAtomPool(ctx sdk.Context, k kpr.Keeper) {
	atomIbcHash := "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"
	// Get Rowan Balance from CLP Module
	clpModuleTotalNativeBalance := k.GetBankKeeper().GetBalance(ctx, types.GetCLPModuleAddress(), types.GetSettlementAsset().Symbol)
	// Get Atom Balance from CLP Module
	clpModulebalanceAtom := k.GetBankKeeper().GetBalance(ctx, types.GetCLPModuleAddress(), atomIbcHash)

	// Get Uint amount from coin
	clpModuleTotalNativeBalanceUint := sdk.NewUintFromString(clpModuleTotalNativeBalance.Amount.String())
	clpModulebalanceAtomUint := sdk.NewUintFromString(clpModulebalanceAtom.Amount.String())

	// Get Atom Pool
	atomPool, err := k.GetPool(ctx, atomIbcHash)
	if err != nil {
		panic(err)
	}

	// Calculate total native balance of all pools
	pools := k.GetPools(ctx)
	poolTotalNativeBalance := sdk.ZeroUint()
	for _, pool := range pools {
		if pool.ExternalAsset.Symbol != atomIbcHash {
			poolTotalNativeBalance = poolTotalNativeBalance.Add(pool.NativeAssetBalance)
		}
	}
	// Set Atom pool back
	atomPool.ExternalAssetBalance = clpModulebalanceAtomUint
	atomPool.NativeAssetBalance = clpModuleTotalNativeBalanceUint.Sub(poolTotalNativeBalance)
	atomPool.ExternalLiabilities = sdk.ZeroUint()
	atomPool.NativeLiabilities = sdk.ZeroUint()
	atomPool.ExternalCustody = sdk.ZeroUint()
	atomPool.NativeCustody = sdk.ZeroUint()
	err = k.SetPool(ctx, &atomPool)
	if err != nil {
		panic(err)
	}
}

var blockTime *time.Time = nil

func MeasureBlockTime(ctx sdk.Context, k kpr.Keeper) {
	now := time.Now()
	if blockTime == nil {
		blockTime = &now
		return
	}

	elapsed := now.Sub(*blockTime)
	blockTime = &now
	k.Logger(ctx).Info(fmt.Sprint("Block took ", elapsed.Seconds(), "s to execute"))
}
