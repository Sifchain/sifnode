package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	"strconv"
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
	var pmtpCurrentRunningRate sdk.Dec
	// Epoch counters are used to keep track of policy execution
	// EpochCounter tracks the number of Epochs completed
	// BlockCounter tracks the number of Blocks completed in an Epoch
	// Manage Block Counter and Calculate R running
	if currentHeight >= pmtpPeriodStartBlock &&
		currentHeight <= pmtpPeriodEndBlock &&
		k.GetPmtpEpoch(ctx).EpochCounter > 0 {
		pmtpCurrentRunningRate = k.PolicyCalculations(ctx)
	}
	// Manage Epoch Counter
	if k.GetPmtpEpoch(ctx).BlockCounter == 0 {
		k.DecrementEpochCounter(ctx)
		k.SetBlockCounter(ctx, k.GetPmtpEpochLength(ctx))
	}
	// Last epoch is not included
	if k.GetPmtpEpoch(ctx).BlockCounter == 0 {
		// Setting it to zero to check when we start policy
		k.SetBlockCounter(ctx, 0)
		_ = ctx.EventManager().EmitTypedEvent(&types.EventPolicy{
			EventType:            "policy_end",
			PmtpPeriodStartBlock: strconv.Itoa(int(pmtpPeriodStartBlock)),
			PmtpPeriodEndBlock:   strconv.Itoa(int(pmtpPeriodEndBlock)),
		})
		k.Logger(ctx).Info(fmt.Sprintf("Ending Policy | Start Height : %d | End Height : %d", pmtpPeriodStartBlock, pmtpPeriodEndBlock))
		return
	}

	err := k.PolicyRun(ctx, pmtpCurrentRunningRate)
	if err != nil {
		panic(err)
	}
}
