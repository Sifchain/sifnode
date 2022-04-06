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
	FsLiquidityRemovalLockPeriod.String(FlagLiquidityRemovalLockPeriod, "", "Lock Period")
	FsLiquidityRemovalCancelPeriod.String(FlagLiquidityRemovalCancelPeriod, "", "Unlock Period")
	FsDefaultMultiplier.String(FlagDefaultMultiplier, "", "Pool Multiplier")
	FsFlagRewardPeriods.String(FlagRewardPeriods, "", "Path to Json File containing reward periods")
}
