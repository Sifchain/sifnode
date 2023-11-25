package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagAssetSymbol                     = "symbol"
	FlagUnits                           = "units"
	FlagSentAssetSymbol                 = "sentSymbol"
	FlagReceivedAssetSymbol             = "receivedSymbol"
	FlagNativeAssetAmount               = "nativeAmount"
	FlagExternalAssetAmount             = "externalAmount"
	FlagWBasisPoints                    = "wBasis"
	FlagAsymmetry                       = "asymmetry"
	FlagWithdrawUnits                   = "withdrawUnits"
	FlagAmount                          = "sentAmount"
	FlagMinimumReceivingAmount          = "minReceivingAmount"
	FlagLiquidityRemovalLockPeriod      = "lockPeriod"
	FlagLiquidityRemovalCancelPeriod    = "cancelPeriod"
	FlagRewardsLockPeriod               = "rewardsLockPeriod"
	FlagRewardsEpochIdentifier          = "rewardsEpochIdentifier"
	FlagRewardsDistribute               = "rewardsDistribute"
	FlagDefaultMultiplier               = "defaultMultiplier"
	FlagRewardPeriods                   = "path"
	FlagBlockRate                       = "blockRate"
	FlagRunningRate                     = "runningRate"
	FlagEndCurrentPolicy                = "endPolicy"
	FlagPeriodGovernanceRate            = "rGov"
	FlagPmtpPeriodEpochLength           = "epochLength"
	FlagPmtpPeriodStartBlock            = "pmtp_start"
	FlagPmtpPeriodEndBlock              = "pmtp_end"
	FlagNewPolicy                       = "newPolicy"
	FlagMintParams                      = "mint-params"
	FlagMinter                          = "minter"
	FlagSymmetryThreshold               = "threshold"
	FlagSymmetryRatioThreshold          = "ratio"
	FlagMaxRowanLiquidityThreshold      = "maxRowanLiquidityThreshold"
	FlagMaxRowanLiquidityThresholdAsset = "maxRowanLiquidityThresholdAsset"
	FlagLiquidityProtectionEpochLength  = "epochLength"
	FlagCurrentRowanLiquidityThreshold  = "currentRowanLiquidityThreshold"
	FlagLiquidityProtectionIsActive     = "isActive"
	FlagProviderDistributionPeriods     = "path"
	FlagSwapFeeParams                   = "path"
)

// common flagsets to add to various functions
var (
	FsAssetSymbol                     = flag.NewFlagSet("", flag.ContinueOnError)
	FsUnits                           = flag.NewFlagSet("", flag.ContinueOnError)
	FsNativeAssetAmount               = flag.NewFlagSet("", flag.ContinueOnError)
	FsExternalAssetAmount             = flag.NewFlagSet("", flag.ContinueOnError)
	FsWBasisPoints                    = flag.NewFlagSet("", flag.ContinueOnError)
	FsAsymmetry                       = flag.NewFlagSet("", flag.ContinueOnError)
	FsWithdrawUnits                   = flag.NewFlagSet("", flag.ContinueOnError)
	FsSentAssetSymbol                 = flag.NewFlagSet("", flag.ContinueOnError)
	FsReceivedAssetSymbol             = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount                          = flag.NewFlagSet("", flag.ContinueOnError)
	FsMinReceivingAmount              = flag.NewFlagSet("", flag.ContinueOnError)
	FsLiquidityRemovalLockPeriod      = flag.NewFlagSet("", flag.ContinueOnError)
	FsLiquidityRemovalCancelPeriod    = flag.NewFlagSet("", flag.ContinueOnError)
	FsRewardsLockPeriod               = flag.NewFlagSet("", flag.ContinueOnError)
	FsRewardsEpochIdentifier          = flag.NewFlagSet("", flag.ContinueOnError)
	FsRewardsDistribute               = flag.NewFlagSet("", flag.ContinueOnError)
	FsDefaultMultiplier               = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagRewardPeriods               = flag.NewFlagSet("", flag.ContinueOnError)
	FsBlockRate                       = flag.NewFlagSet("", flag.ContinueOnError)
	FsRunningRate                     = flag.NewFlagSet("", flag.ContinueOnError)
	FsEndCurrentPolicy                = flag.NewFlagSet("", flag.ContinueOnError)
	FsPeriodGovernanceRate            = flag.NewFlagSet("", flag.ContinueOnError)
	FsPmtpPeriodEpochLength           = flag.NewFlagSet("", flag.ContinueOnError)
	FsPmtpPeriodStartBlock            = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagPmtpPeriodEndBlock          = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagNewPolicy                   = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagMintParams                  = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagMinter                      = flag.NewFlagSet("", flag.ContinueOnError)
	FsSymmetryThreshold               = flag.NewFlagSet("", flag.ContinueOnError)
	FsSymmetryRatioThreshold          = flag.NewFlagSet("", flag.ContinueOnError)
	FsMaxRowanLiquidityThreshold      = flag.NewFlagSet("", flag.ContinueOnError)
	FsMaxRowanLiquidityThresholdAsset = flag.NewFlagSet("", flag.ContinueOnError)
	FsLiquidityThresholdIsActive      = flag.NewFlagSet("", flag.ContinueOnError)
	FsLiquidityProtectionEpochLength  = flag.NewFlagSet("", flag.ContinueOnError)
	FsCurrentRowanLiquidityThreshold  = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagProviderDistributionPeriods = flag.NewFlagSet("", flag.ContinueOnError)
	FsFlagSwapFeeParams               = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {

	FsAssetSymbol.String(FlagAssetSymbol, "", "Symbol for Asset")
	FsUnits.String(FlagUnits, "", "Liquidity provider units")
	FsNativeAssetAmount.String(FlagNativeAssetAmount, "", "Native Asset Amount")
	FsExternalAssetAmount.String(FlagExternalAssetAmount, "", "External Asset Amount")
	FsWBasisPoints.String(FlagWBasisPoints, "", "WBasis Points ")
	FsAsymmetry.String(FlagAsymmetry, "", "Asymmetry")
	FsWithdrawUnits.String(FlagWithdrawUnits, "", "Withdraw Units ")
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
	FsRewardsLockPeriod.String(FlagRewardsLockPeriod, "", "Rewards Lock Period")
	FsRewardsEpochIdentifier.String(FlagRewardsEpochIdentifier, "", "Rewards Epoch Identifier")
	FsRewardsDistribute.String(FlagRewardsDistribute, "", "Rewards Distribute")
	FsDefaultMultiplier.String(FlagDefaultMultiplier, "", "Pool Multiplier")
	FsFlagRewardPeriods.String(FlagRewardPeriods, "", "Path to Json File containing reward periods")
	FsFlagMintParams.String(FlagMintParams, "", "Inflation")
	FsFlagMinter.String(FlagMinter, "", "Inflation Max")
	FsSymmetryThreshold.String(FlagSymmetryThreshold, "", "Set slippage adjustement threshold for symmetric liquitidy add")
	FsSymmetryRatioThreshold.String(FlagSymmetryRatioThreshold, "", "Set ratio threshold for symmetric liquitidy add")
	FsMaxRowanLiquidityThreshold.String(FlagMaxRowanLiquidityThreshold, "", "Set max rowan liquidity threshold value according to asset used in maxRowanLiquidityThresholdAsset")
	FsMaxRowanLiquidityThresholdAsset.String(FlagMaxRowanLiquidityThresholdAsset, "", "Set max rowan liquidity threshold asset to use to determine native price (if set to 'rowan' the threshold can be define as rowan)")
	FsLiquidityProtectionEpochLength.String(FlagLiquidityProtectionEpochLength, "", "Set liquidity protection epoch length")
	FsLiquidityThresholdIsActive.String(FlagLiquidityProtectionIsActive, "", "Set liquidity protection isActive")
	FsCurrentRowanLiquidityThreshold.String(FlagCurrentRowanLiquidityThreshold, "", "Set current rowan liquidity threshold value")
	FsFlagProviderDistributionPeriods.String(FlagProviderDistributionPeriods, "", "Path to Json File containing LP provider distribution periods")
	FsFlagSwapFeeParams.String(FlagProviderDistributionPeriods, "", "Path to Json File containing swap fee params")
}
