package cli

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/cli"
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
		GetCmdSwap(),
		GetCmdDecommissionPool(),
	)

	return clpTxCmd
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

			asset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			externalAmount := viper.GetString(FlagExternalAssetAmount)
			nativeAmount := viper.GetString(FlagNativeAssetAmount)
			signer := clientCtx.GetFromAddress()

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
	cmd.Flags().StringP(cli.OutputFlag, "o", "text", "Output format (text|json)")

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
