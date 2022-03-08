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
	if currentHeight == pmtpPeriodStartBlock &&
		k.GetPmtpEpoch(ctx).EpochCounter == 0 &&
		k.GetPmtpEpoch(ctx).BlockCounter == 0 {
		k.PolicyStart(ctx)
	}

	// Manage Block Counter and Calculate R running
	if currentHeight >= pmtpPeriodStartBlock &&
		currentHeight <= pmtpPeriodEndBlock &&
		k.GetPmtpEpoch(ctx).EpochCounter > 0 {
		k.PolicyCalculations(ctx)
	}
	// Manage Epoch Counter
	if k.GetPmtpEpoch(ctx).BlockCounter == 0 {
		k.DecrementEpochCounter(ctx)
	}

	err := k.PolicyRun(ctx)
	if err != nil {
		panic(err)
	}
}
