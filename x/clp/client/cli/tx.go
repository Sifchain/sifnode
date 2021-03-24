package cli

import (
	"bufio"
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	clpTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	clpTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreatePool(cdc),
		GetCmdAddLiquidity(cdc),
		GetCmdRemoveLiquidity(cdc),
		GetCmdSwap(cdc),
		GetCmdDecommissionPool(cdc),
		MultiSendTxCmd(cdc),
	)...)

	return clpTxCmd
}

func GetCmdCreatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool --from [key] --symbol [asset-symbol] --nativeAmount [amount] --externalAmount [amount]",
		Short: "Create new liquidity pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			asset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			externalAmount := viper.GetString(FlagExternalAssetAmount)
			nativeAmount := viper.GetString(FlagNativeAssetAmount)
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgCreatePool(signer, asset, sdk.NewUintFromString(nativeAmount), sdk.NewUintFromString(externalAmount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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
	return cmd
}

func GetCmdDecommissionPool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decommission-pool",
		Short: "decommission liquidity pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			symbol := viper.GetString(FlagAssetSymbol)
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgDecommissionPool(signer, symbol)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	if err := cmd.MarkFlagRequired(FlagAssetSymbol); err != nil {
		log.Println("MarkFlagRequired failed: ", err.Error())
	}

	return cmd
}

func GetCmdAddLiquidity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-liquidity",
		Short: "Add liquidity to a pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			externalAsset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			externalAmount := viper.GetString(FlagExternalAssetAmount)
			nativeAmount := viper.GetString(FlagNativeAssetAmount)
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgAddLiquidity(signer, externalAsset, sdk.NewUintFromString(nativeAmount), sdk.NewUintFromString(externalAmount))

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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

	return cmd
}

func GetCmdRemoveLiquidity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-liquidity",
		Short: "Remove liquidity from a pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			externalAsset := types.NewAsset(viper.GetString(FlagAssetSymbol))
			wb := viper.GetString(FlagWBasisPoints)
			as := viper.GetString(FlagAsymmetry)
			signer := cliCtx.GetFromAddress()
			wBasis, ok := sdk.NewIntFromString(wb)
			if !ok {
				return types.ErrOverFlow
			}
			asymmetry, ok := sdk.NewIntFromString(as)
			if !ok {
				return types.ErrOverFlow
			}
			msg := types.NewMsgRemoveLiquidity(signer, externalAsset, wBasis, asymmetry)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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

	return cmd
}

func GetCmdSwap(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap",
		Short: "Swap tokens using liquidity pools",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			sentAsset := types.NewAsset(viper.GetString(FlagSentAssetSymbol))
			receivedAsset := types.NewAsset(viper.GetString(FlagReceivedAssetSymbol))
			sentAmount := viper.GetString(FlagAmount)
			minReceivingAmount := viper.GetString(FlagMinimumReceivingAmount)
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgSwap(signer, sentAsset, receivedAsset, sdk.NewUintFromString(sentAmount), sdk.NewUintFromString(minReceivingAmount))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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
	return cmd
}

func MultiSendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-send [from1] [from2] [to_address] [amount1] [amount2]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins1, err := sdk.ParseCoins(args[3])
			if err != nil {
				return err
			}
			coins2, err := sdk.ParseCoins(args[4])
			if err != nil {
				return err
			}
			from1, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			from2, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			// build and sign the transaction, then broadcast to Tendermint
			inputs := []bank.Input{bank.NewInput(from1, coins1), bank.NewInput(from2, coins2)}
			outputs := []bank.Output{bank.NewOutput(to, coins1.Add(coins2...))}
			msg := bank.NewMsgMultiSend(inputs, outputs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
