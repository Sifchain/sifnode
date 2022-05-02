package clp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	params := keeper.GetRewardsParams(ctx)
	pools := keeper.GetPools(ctx)
	currentPeriod := keeper.GetCurrentRewardPeriod(ctx, params)
	if currentPeriod != nil && !currentPeriod.RewardPeriodAllocation.IsZero() {
		err := keeper.DistributeDepthRewards(ctx, currentPeriod, pools)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	// get current block height
	currentHeight := ctx.BlockHeight()
	// default to current pmtp current running rate value
	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	// get PMTP period params
	params := k.GetPmtpParams(ctx)
	currentPolicy := k.GetCurrentPmtpPolicy(ctx, params)
	if currentPolicy != nil && currentPolicy.PmtpPeriodStartBlock > 0 && currentPolicy.PmtpPeriodEndBlock > 0 {
		pmtpPeriodStartBlock := currentPolicy.PmtpPeriodStartBlock
		pmtpPeriodEndBlock := currentPolicy.PmtpPeriodEndBlock
		// Start Policy
		if currentHeight == pmtpPeriodStartBlock &&
			k.GetPmtpEpoch(ctx).EpochCounter == 0 &&
			k.GetPmtpEpoch(ctx).BlockCounter == 0 {
			k.PolicyStart(ctx, currentPolicy)
			_ = ctx.EventManager().EmitTypedEvent(&types.EventPolicy{
				EventType:            "policy_start",
				PmtpPeriodStartBlock: strconv.Itoa(int(pmtpPeriodStartBlock)),
				PmtpPeriodEndBlock:   strconv.Itoa(int(pmtpPeriodEndBlock)),
			})
			k.Logger(ctx).Info(fmt.Sprintf("Starting new policy | Start Height : %d | End Height : %d", pmtpPeriodStartBlock, pmtpPeriodEndBlock))
		}

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
			pmtpCurrentRunningRate = k.PolicyCalculations(ctx, currentPolicy)
			k.DecrementBlockCounter(ctx)
		}
		// Manage Epoch Counter
		if k.GetPmtpEpoch(ctx).BlockCounter == 0 &&
			currentHeight < pmtpPeriodEndBlock &&
			currentHeight >= pmtpPeriodStartBlock {
			k.DecrementEpochCounter(ctx)
			k.SetBlockCounter(ctx, currentPolicy.PmtpPeriodEpochLength)
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
	}

	err := k.PolicyRun(ctx, pmtpCurrentRunningRate)
	if err != nil {
		panic(err)
	}
}
