package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagAssetSymbol                  = "symbol"
	FlagUnits                        = "units"
	FlagSentAssetSymbol              = "sentSymbol"
	FlagReceivedAssetSymbol          = "receivedSymbol"
	FlagNativeAssetAmount            = "nativeAmount"
	FlagExternalAssetAmount          = "externalAmount"
	FlagWBasisPoints                 = "wBasis"
	FlagAsymmetry                    = "asymmetry"
	FlagAmount                       = "sentAmount"
	FlagMinimumReceivingAmount       = "minReceivingAmount"
	FlagLiquidityRemovalLockPeriod   = "lockPeriod"
	FlagLiquidityRemovalCancelPeriod = "cancelPeriod"
	FlagDefaultMultiplier            = "defaultMultiplier"
	FlagRewardPeriods                = "path"
	FlagBlockRate                    = "blockRate"
	FlagRunningRate                  = "runningRate"
	FlagEndCurrentPolicy             = "endPolicy"
	FlagPeriodGovernanceRate         = "rGov"
	FlagPmtpPeriodEpochLength        = "epochLength"
	FlagPmtpPeriodStartBlock         = "pmtp_start"
	FlagPmtpPeriodEndBlock           = "pmtp_end"
	FlagNewPolicy                    = "newPolicy"
	FlagInflation                    = "inflation"
	FlagInflationMax                 = "inflationMax"
	FlagInflationMin                 = "inflationMin"
	FlagAnnualProvisions             = "annualProvision"
)

// common flagsets to add to various functions
var (
	FsAssetSymbol                  = flag.NewFlagSet("", flag.ContinueOnError)
	FsUnits                        = flag.NewFlagSet("", flag.ContinueOnError)
	FsNativeAssetAmount            = flag.NewFlagSet("", flag.ContinueOnError)
	FsExternalAssetAmount          = flag.NewFlagSet("", flag.ContinueOnError)
	FsWBasisPoints                 = flag.NewFlagSet("", flag.ContinueOnError)
	FsAsymmetry                    = flag.NewFlagSet("", flag.ContinueOnError)
	FsSentAssetSymbol              = flag.NewFlagSet("", flag.ContinueOnError)
	FsReceivedAssetSymbol          = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount                       = flag.NewFlagSet("", flag.ContinueOnError)
	FsMinReceivingAmount           = flag.NewFlagSet("", flag.ContinueOnError)
	FsLiquidityRemovalLockPeriod   = flag.NewFlagSet("", flag.ContinueOnError)
	FsLiquidityRemovalCancelPeriod = flag.NewFlagSet("", flag.ContinueOnError)
	FsDefaultMultiplier            = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagRewardPeriods            = flag.NewFlagSet("", flag.ContinueOnError)
	FsBlockRate                    = flag.NewFlagSet("", flag.ContinueOnError)
	FsRunningRate                  = flag.NewFlagSet("", flag.ContinueOnError)
	FsEndCurrentPolicy             = flag.NewFlagSet("", flag.ContinueOnError)
	FsPeriodGovernanceRate         = flag.NewFlagSet("", flag.ContinueOnError)
	FsPmtpPeriodEpochLength        = flag.NewFlagSet("", flag.ContinueOnError)
	FsPmtpPeriodStartBlock         = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagPmtpPeriodEndBlock       = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagNewPolicy                = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagInflation                = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagInflationMax             = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagInflationMin             = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagAnnualProvisions         = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {

	FsAssetSymbol.String(FlagAssetSymbol, "", "Symbol for Asset")
	FsUnits.String(FlagUnits, "", "Liquidity provider units")
	FsNativeAssetAmount.String(FlagNativeAssetAmount, "", "Native Asset Amount")
	FsExternalAssetAmount.String(FlagExternalAssetAmount, "", "External Asset Amount")
	FsWBasisPoints.String(FlagWBasisPoints, "", "WBasis Points ")
	FsAsymmetry.String(FlagAsymmetry, "", "Asymmetry")
	FsSentAssetSymbol.String(FlagSentAssetSymbol, "", "Symbol for Sent Asset")
	FsReceivedAssetSymbol.String(FlagReceivedAssetSymbol, "", "Symbol for Received Asset")
	FsAmount.String(FlagAmount, "", "Sent amount")
	FsMinReceivingAmount.String(FlagMinimumReceivingAmount, "", "Min threshold for receiving amount")
	FsBlockRate.String(FlagBlockRate, "", "Flag to modify Block rate")
	FsRunningRate.String(FlagRunningRate, "", "Flag to modify Running rate")
	FsEndCurrentPolicy.String(FlagEndCurrentPolicy, "", "Set flag to true to end current policy")
	FsPeriodGovernanceRate.String(FlagPeriodGovernanceRate, "", "Modify rGov")
	FsPmtpPeriodEpochLength.String(FlagPmtpPeriodEpochLength, "", "Modify rGov")
	FsPmtpPeriodStartBlock.String(FlagPmtpPeriodStartBlock, "", "Modify pmtp start block")
	FsFlagPmtpPeriodEndBlock.String(FlagPmtpPeriodEndBlock, "", "Modify pmtp end block")
	FsFlagNewPolicy.String(FlagNewPolicy, "", "Set a new policy / Modify existing policy")
	FsLiquidityRemovalLockPeriod.String(FlagLiquidityRemovalLockPeriod, "", "Lock Period")
	FsLiquidityRemovalCancelPeriod.String(FlagLiquidityRemovalCancelPeriod, "", "Unlock Period")
	FsDefaultMultiplier.String(FlagDefaultMultiplier, "", "Pool Multiplier")
	FsFlagRewardPeriods.String(FlagRewardPeriods, "", "Path to Json File containing reward periods")
	FsFlagInflation.String(FlagInflation, "", "Inflation")
	FsFlagInflationMax.String(FlagInflationMax, "", "Inflation Max")
	FsFlagInflationMin.String(FlagInflationMin, "", "Inflation Min")
	FsFlagAnnualProvisions.String(FlagAnnualProvisions, "", "Annual Provisions")
}
