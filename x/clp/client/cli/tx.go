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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	)...)

	return clpTxCmd
}

func GetCmdCreatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool",
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
	cmd.MarkFlagRequired(FlagAssetSymbol)
	cmd.MarkFlagRequired(FlagExternalAssetAmount)
	cmd.MarkFlagRequired(FlagNativeAssetAmount)

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
	cmd.MarkFlagRequired(FlagAssetSymbol)

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
	cmd.MarkFlagRequired(FlagAssetSymbol)
	cmd.MarkFlagRequired(FlagExternalAssetAmount)
	cmd.MarkFlagRequired(FlagNativeAssetAmount)

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
	cmd.MarkFlagRequired(FlagAssetSymbol)
	cmd.MarkFlagRequired(FlagWBasisPoints)
	cmd.MarkFlagRequired(FlagAsymmetry)

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
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgSwap(signer, sentAsset, receivedAsset, sdk.NewUintFromString(sentAmount))

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsSentAssetSymbol)
	cmd.Flags().AddFlagSet(FsReceivedAssetSymbol)
	cmd.Flags().AddFlagSet(FsAmount)

	cmd.MarkFlagRequired(FlagSentAssetSymbol)
	cmd.MarkFlagRequired(FlagReceivedAssetSymbol)
	cmd.MarkFlagRequired(FlagAmount)

	return cmd
}
