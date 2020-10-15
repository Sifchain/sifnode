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
			asset := types.NewAsset(viper.GetString(FlagAssetSourceChain),
				viper.GetString(FlagAssetSymbol),
				viper.GetString(FlagAssetTicker))
			externalAmount := viper.GetUint(FlagExternalAssetAmount)
			nativeAmount := viper.GetUint(FlagNativeAssetAmount)
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgCreatePool(signer, asset, nativeAmount, externalAmount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSourceChain)
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsAssetTicker)
	cmd.Flags().AddFlagSet(FsExternalAssetAmount)
	cmd.Flags().AddFlagSet(FsNativeAssetAmount)
	cmd.MarkFlagRequired(FlagAssetSourceChain)
	cmd.MarkFlagRequired(FlagAssetSymbol)
	cmd.MarkFlagRequired(FlagAssetTicker)
	cmd.MarkFlagRequired(FlagExternalAssetAmount)
	cmd.MarkFlagRequired(FlagNativeAssetAmount)

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
			externalAsset := types.NewAsset(viper.GetString(FlagAssetSourceChain),
				viper.GetString(FlagAssetSymbol),
				viper.GetString(FlagAssetTicker))
			externalAmount := viper.GetUint(FlagExternalAssetAmount)
			nativeAmount := viper.GetUint(FlagNativeAssetAmount)
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgAddLiquidity(signer, externalAsset, externalAmount, nativeAmount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSourceChain)
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsAssetTicker)
	cmd.Flags().AddFlagSet(FsExternalAssetAmount)
	cmd.Flags().AddFlagSet(FsNativeAssetAmount)
	cmd.MarkFlagRequired(FlagAssetSourceChain)
	cmd.MarkFlagRequired(FlagAssetSymbol)
	cmd.MarkFlagRequired(FlagAssetTicker)
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
			externalAsset := types.NewAsset(viper.GetString(FlagAssetSourceChain),
				viper.GetString(FlagAssetSymbol),
				viper.GetString(FlagAssetTicker))
			wBasis := viper.GetUint(FlagWBasisPoints)
			asymmetry := viper.GetUint(FlagAsymmetry)
			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgRemoveLiquidity(signer, externalAsset, wBasis, asymmetry)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsAssetSourceChain)
	cmd.Flags().AddFlagSet(FsAssetSymbol)
	cmd.Flags().AddFlagSet(FsAssetTicker)
	cmd.Flags().AddFlagSet(FsWBasisPoints)
	cmd.Flags().AddFlagSet(FsAsymmetry)
	cmd.MarkFlagRequired(FlagAssetSourceChain)
	cmd.MarkFlagRequired(FlagAssetSymbol)
	cmd.MarkFlagRequired(FlagAssetTicker)
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
			sentAsset := types.NewAsset(viper.GetString(FlagSentAssetSourceChain),
				viper.GetString(FlagSentAssetSymbol),
				viper.GetString(FlagSentAssetTicker))
			receivedAsset := types.NewAsset(viper.GetString(FlagReceivedAssetSourceChain),
				viper.GetString(FlagReceivedAssetSymbol),
				viper.GetString(FlagReceivedAssetTicker))
			sentAmount := viper.GetUint(FlagAmount)

			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgSwap(signer, sentAsset, receivedAsset, sentAmount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(FsSentAssetSourceChain)
	cmd.Flags().AddFlagSet(FsSentAssetSymbol)
	cmd.Flags().AddFlagSet(FsSentAssetTicker)
	cmd.Flags().AddFlagSet(FsReceivedAssetSourceChain)
	cmd.Flags().AddFlagSet(FsReceivedAssetSymbol)
	cmd.Flags().AddFlagSet(FsReceivedAssetTicker)
	cmd.Flags().AddFlagSet(FsAmount)

	cmd.MarkFlagRequired(FlagSentAssetSourceChain)
	cmd.MarkFlagRequired(FlagSentAssetSymbol)
	cmd.MarkFlagRequired(FlagSentAssetTicker)
	cmd.MarkFlagRequired(FlagReceivedAssetSourceChain)
	cmd.MarkFlagRequired(FlagReceivedAssetSymbol)
	cmd.MarkFlagRequired(FlagReceivedAssetTicker)
	cmd.MarkFlagRequired(FlagAmount)

	return cmd
}
