package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

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
		GetCmdSwap(),
		GetCmdDecommissionPool(),
		GetCmdUnlockLiquidity(),
		GetCmdUpdateRewardParams(),
		GetCmdAddRewardPeriod(),
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
			fmt.Println(rewardPeriods)
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
			defaultMultiplier, err := sdk.NewDecFromStr(viper.GetString(FlagDefaultMultiplier))
			if err != nil {
				return err
			}
			msg := types.MsgUpdateRewardsParamsRequest{
				Signer:                       signer.String(),
				LiquidityRemovalCancelPeriod: viper.GetUint64(FlagLiquidityRemovalCancelPeriod),
				LiquidityRemovalLockPeriod:   viper.GetUint64(FlagLiquidityRemovalLockPeriod),
				DefaultMultiplier:            &defaultMultiplier,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().AddFlagSet(FsLiquidityRemovalCancelPeriod)
	cmd.Flags().AddFlagSet(FsLiquidityRemovalLockPeriod)
	cmd.Flags().AddFlagSet(FsDefaultMultiplier)
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

			flags := cmd.Flags()

			assetSymbol, err := flags.GetString(FlagAssetSymbol)
			if err != nil {
				return err
			}

			externalAmount, err := flags.GetString(FlagExternalAssetAmount)
			if err != nil {
				return err
			}

			nativeAmount, err := flags.GetString(FlagNativeAssetAmount)
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
			units := viper.GetUint64(FlagUnits)
			unitsInt := sdk.NewUint(units)
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
