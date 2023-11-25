package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"log"

	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	clpTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	clpTxCmd.AddCommand(
		GetCmdCreatePool(),
		GetCmdAddLiquidity(),
		GetCmdRemoveLiquidity(),
		GetCmdRemoveLiquidityUnits(),
		GetCmdSwap(),
		GetCmdDecommissionPool(),
		GetCmdUnlockLiquidity(),
		GetCmdCancelUnlockLiquidity(),
		GetCmdUpdateRewardParams(),
		GetCmdAddRewardPeriod(),
		GetCmdModifyPmtpRates(),
		GetCmdUpdatePmtpParams(),
		GetCmdUpdateStakingRewards(),
		GetCmdSetSymmetryThreshold(),
		GetCmdUpdateLiquidityProtectionParams(),
		GetCmdModifyLiquidityProtectionRates(),
		GetCmdSetProviderDistributionPeriods(),
		GetCmdSetSwapFeeParams(),
		CmdAddLiquidityToRewardsBucket(),
	)

	return clpTxCmd
}
func GetCmdAddRewardPeriod() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-period",
		Short: "Update reward params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var rewardPeriods []*types.RewardPeriod
			signer := clientCtx.GetFromAddress()
			filePath := viper.GetString(FlagRewardPeriods)
			file, err := filepath.Abs(filePath)
			if err != nil {
				return err
			}
			input, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			err = json.Unmarshal(input, &rewardPeriods)
			if err != nil {
				return err
			}
			msg := types.MsgAddRewardPeriodRequest{
				Signer:        signer.String(),
				RewardPeriods: rewardPeriods,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsFlagRewardPeriods)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdUpdateRewardParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-params",
		Short: "Update reward params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if err != nil {
				return err
			}
			msg := types.MsgUpdateRewardsParamsRequest{
				Signer:                       signer.String(),
				LiquidityRemovalCancelPeriod: viper.GetUint64(FlagLiquidityRemovalCancelPeriod),
				LiquidityRemovalLockPeriod:   viper.GetUint64(FlagLiquidityRemovalLockPeriod),
				RewardsLockPeriod:            viper.GetUint64(FlagRewardsLockPeriod),
				RewardsEpochIdentifier:       viper.GetString(FlagRewardsEpochIdentifier),
				RewardsDistribute:            viper.GetBool(FlagRewardsDistribute),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsLiquidityRemovalCancelPeriod)
	cmd.Flags().AddFlagSet(FsLiquidityRemovalLockPeriod)
	cmd.Flags().AddFlagSet(FsRewardsLockPeriod)
	cmd.Flags().AddFlagSet(FsRewardsEpochIdentifier)
	cmd.Flags().AddFlagSet(FsRewardsDistribute)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdSetSymmetryThreshold() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-symmetry-threshold",
		Short: "Set symmetry threshold",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if err != nil {
				return err
			}
			threshold, err := sdk.NewDecFromStr(viper.GetString(FlagSymmetryThreshold))
			if err != nil {
				return err
			}
			ratioThreshold, err := sdk.NewDecFromStr(viper.GetString(FlagSymmetryRatioThreshold))
			if err != nil {
				return err
			}
			msg := types.MsgSetSymmetryThreshold{
				Signer:    signer.String(),
				Threshold: threshold,
				Ratio:     ratioThreshold,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsSymmetryThreshold)
	cmd.Flags().AddFlagSet(FsSymmetryRatioThreshold)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdSetSwapFeeParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-swap-fee-params",
		Short: "Set swap fee params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			if err != nil {
				return err
			}
			filePath := viper.GetString(FlagSwapFeeParams)
			file, err := filepath.Abs(filePath)
			if err != nil {
				return err
			}
			input, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			var swapFeeParams types.SwapFeeParams
			err = json.Unmarshal(input, &swapFeeParams)
			if err != nil {
				return err
			}
			msg := types.MsgUpdateSwapFeeParamsRequest{
				Signer:             signer.String(),
				DefaultSwapFeeRate: swapFeeParams.DefaultSwapFeeRate,
				TokenParams:        swapFeeParams.TokenParams,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsFlagSwapFeeParams)
	if err := cmd.MarkFlagRequired(FlagSwapFeeParams); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool --from [key] --symbol [asset-symbol] --nativeAmount [amount] --externalAmount [amount]",
		Short: "Create new liquidity pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			f := cmd.Flags()

			assetSymbol, err := f.GetString(FlagAssetSymbol)
			if err != nil {
				return err
			}

			externalAmount, err := f.GetString(FlagExternalAssetAmount)
			if err != nil {
				return err
			}

			nativeAmount, err := f.GetString(FlagNativeAssetAmount)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()

			asset := types.NewAsset(assetSymbol)
			msg := types.NewMsgCreatePool(signer, asset, sdk.NewUintFromString(nativeAmount), sdk.NewUintFromString(externalAmount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsExternalAssetAmount)
	cmd.Flags().AddFlagSet(FsNativeAssetAmount)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}

	if err := cmd.MarkFlagRequired(FlagExternalAssetAmount); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagNativeAssetAmount); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdDecommissionPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decommission-pool",
		Short: "decommission liquidity pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			symbol := viper.GetString(FlagAssetSymbol)
			signer := clientCtx.GetFromAddress()
			msg := types.NewMsgDecommissionPool(signer, symbol)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdAddLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity",
		Short: "Add liquidity to a pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			externalAsset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			externalAmount := viper.GetString(FlagExternalAssetAmount)
			nativeAmount := viper.GetString(FlagNativeAssetAmount)
			signer := clientCtx.GetFromAddress()

			msg := types.NewMsgAddLiquidity(signer, externalAsset, sdk.NewUintFromString(nativeAmount), sdk.NewUintFromString(externalAmount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsExternalAssetAmount)
	cmd.Flags().AddFlagSet(FsNativeAssetAmount)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagExternalAssetAmount); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	if err := cmd.MarkFlagRequired(FlagNativeAssetAmount); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdRemoveLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity",
		Short: "Remove liquidity from a pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			externalAsset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			wb := viper.GetString(FlagWBasisPoints)
			as := viper.GetString(FlagAsymmetry)
			signer := clientCtx.GetFromAddress()
			wBasis, ok := sdk.NewIntFromString(wb)
			if !ok {
				return types.ErrOverFlow
			}
			asymmetry, ok := sdk.NewIntFromString(as)
			if !ok {
				return types.ErrOverFlow
			}

			msg := types.NewMsgRemoveLiquidity(signer, externalAsset, wBasis, asymmetry)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsWBasisPoints)
	cmd.Flags().AddFlagSet(FsAsymmetry)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagWBasisPoints); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagAsymmetry); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdRemoveLiquidityUnits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity-units",
		Short: "Remove liquidity from a pool by number of units",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			externalAsset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			wU := viper.GetString(FlagWithdrawUnits)

			signer := clientCtx.GetFromAddress()
			withdrawUnits := sdk.NewUintFromString(wU)

			msg := types.NewMsgRemoveLiquidityUnits(signer, externalAsset, withdrawUnits)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsWithdrawUnits)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagWithdrawUnits); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdModifyPmtpRates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pmtp-rates",
		Short: "Modify pmtp block rate and running rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			isEndPolicy := viper.GetBool(FlagEndCurrentPolicy)
			signer := clientCtx.GetFromAddress()
			msg := types.MsgModifyPmtpRates{
				Signer:      signer.String(),
				BlockRate:   viper.GetString(FlagBlockRate),
				RunningRate: viper.GetString(FlagRunningRate),
				EndPolicy:   isEndPolicy,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsBlockRate)
	cmd.Flags().AddFlagSet(FsRunningRate)
	cmd.Flags().AddFlagSet(FsEndCurrentPolicy)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdUpdateStakingRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking-rewards",
		Short: "Update params to modify staking rewards",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			params := minttypes.Params{}
			minter := minttypes.Minter{}
			filePathParams := viper.GetString(FlagMintParams)
			file, err := filepath.Abs(filePathParams)
			if err != nil {
				return err
			}
			input, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			err = json.Unmarshal(input, &params)
			if err != nil {
				return err
			}
			// Minter is an optional flag
			filePathMinter := viper.GetString(FlagMinter)
			if filePathMinter != "" {
				file, err = filepath.Abs(filePathMinter)
				if err != nil {
					return err
				}
				input, err = ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				err = json.Unmarshal(input, &minter)
				if err != nil {
					return err
				}
			}
			signer := clientCtx.GetFromAddress()
			msg := types.MsgUpdateStakingRewardParams{
				Signer: signer.String(),
				Params: params,
				Minter: minter,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsFlagMintParams)
	cmd.Flags().AddFlagSet(FsFlagMinter)
	if err := cmd.MarkFlagRequired(FlagMintParams); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdUpdatePmtpParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pmtp-params",
		Short: "Update pmtp params to set a new policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			msg := types.MsgUpdatePmtpParams{
				Signer:                   signer.String(),
				PmtpPeriodGovernanceRate: viper.GetString(FlagPeriodGovernanceRate),
				PmtpPeriodEpochLength:    viper.GetInt64(FlagPmtpPeriodEpochLength),
				PmtpPeriodStartBlock:     viper.GetInt64(FlagPmtpPeriodStartBlock),
				PmtpPeriodEndBlock:       viper.GetInt64(FlagPmtpPeriodEndBlock),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsPeriodGovernanceRate)
	cmd.Flags().AddFlagSet(FsPmtpPeriodEpochLength)
	cmd.Flags().AddFlagSet(FsPmtpPeriodStartBlock)
	cmd.Flags().AddFlagSet(FsFlagPmtpPeriodEndBlock)
	if err := cmd.MarkFlagRequired(FlagPmtpPeriodEpochLength); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagPmtpPeriodStartBlock); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagPmtpPeriodEndBlock); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap",
		Short: "Swap tokens using liquidity pools",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sentAsset := types.NewAsset(viper.GetString(FlagSentAssetSymbol))
			receivedAsset := types.NewAsset(viper.GetString(FlagReceivedAssetSymbol))

			sentAmount := viper.GetString(FlagAmount)
			minReceivingAmount := viper.GetString(FlagMinimumReceivingAmount)

			signer := clientCtx.GetFromAddress()

			msg := types.NewMsgSwap(signer, sentAsset, receivedAsset, sdk.NewUintFromString(sentAmount), sdk.NewUintFromString(minReceivingAmount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSentAssetSymbol)
	cmd.Flags().AddFlagSet(FsReceivedAssetSymbol)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(FsMinReceivingAmount)

	if err := cmd.MarkFlagRequired(FlagSentAssetSymbol); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagReceivedAssetSymbol); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagAmount); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagMinimumReceivingAmount); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdUnlockLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbond-liquidity",
		Short: "Unbond liquidity from a pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			externalAsset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			signer := clientCtx.GetFromAddress()
			units := viper.GetString(FlagUnits)
			unitsInt := sdk.NewUintFromString(units)
			msg := types.MsgUnlockLiquidityRequest{
				Signer:        signer.String(),
				ExternalAsset: &externalAsset,
				Units:         unitsInt,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsUnits)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagUnits); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdCancelUnlockLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-unbond",
		Short: "Cancel unbonding of liquidity from a pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			externalAsset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			signer := clientCtx.GetFromAddress()
			units := viper.GetString(FlagUnits)
			unitsInt := sdk.NewUintFromString(units)
			msg := types.MsgCancelUnlock{
				Signer:        signer.String(),
				ExternalAsset: &externalAsset,
				Units:         unitsInt,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsUnits)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagUnits); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdUpdateLiquidityProtectionParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidity-protection-params",
		Short: "Update liquidity protection params to set new threshold",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			msg := types.MsgUpdateLiquidityProtectionParams{
				Signer:                          signer.String(),
				MaxRowanLiquidityThreshold:      sdk.NewUintFromString(viper.GetString(FlagMaxRowanLiquidityThreshold)),
				MaxRowanLiquidityThresholdAsset: viper.GetString(FlagMaxRowanLiquidityThresholdAsset),
				EpochLength:                     viper.GetUint64(FlagLiquidityProtectionEpochLength),
				IsActive:                        viper.GetBool(FlagLiquidityProtectionIsActive),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsMaxRowanLiquidityThreshold)
	cmd.Flags().AddFlagSet(FsMaxRowanLiquidityThresholdAsset)
	cmd.Flags().AddFlagSet(FsPmtpPeriodEpochLength)
	cmd.Flags().AddFlagSet(FsLiquidityThresholdIsActive)
	if err := cmd.MarkFlagRequired(FlagPmtpPeriodEpochLength); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagMaxRowanLiquidityThreshold); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagMaxRowanLiquidityThresholdAsset); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	if err := cmd.MarkFlagRequired(FlagLiquidityProtectionIsActive); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdModifyLiquidityProtectionRates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidity-protection-rates",
		Short: "Modify liquidity protection current rowan liquidity threshold",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			msg := types.MsgModifyLiquidityProtectionRates{
				Signer:                         signer.String(),
				CurrentRowanLiquidityThreshold: sdk.NewUintFromString(viper.GetString(FlagCurrentRowanLiquidityThreshold)),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsCurrentRowanLiquidityThreshold)
	if err := cmd.MarkFlagRequired(FlagCurrentRowanLiquidityThreshold); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdSetProviderDistributionPeriods() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-lppd-params",
		Short: "Set LP provider distribution params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var distributionPeriods []*types.ProviderDistributionPeriod
			signer := clientCtx.GetFromAddress()
			filePath := viper.GetString(FlagProviderDistributionPeriods)
			file, err := filepath.Abs(filePath)
			if err != nil {
				return err
			}
			input, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			err = json.Unmarshal(input, &distributionPeriods)
			if err != nil {
				return err
			}
			msg := types.MsgAddProviderDistributionPeriodRequest{
				Signer:              signer.String(),
				DistributionPeriods: distributionPeriods,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().AddFlagSet(FsFlagProviderDistributionPeriods)
	if err := cmd.MarkFlagRequired(FlagProviderDistributionPeriods); err != nil {
		log.Println("MarkFlagRequired  failed: ", err.Error())
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
