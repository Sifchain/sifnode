package cli

import (
	"bufio"
	"fmt"

	"github.com/Sifchain/sifnode/x/trees/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	faucetTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateTree(cdc),GetCmdCreateOrder(cdc))...)

	return faucetTxCmd
}

// TX to request coins from faucet module account to the requesters account
func GetCmdCreateTree(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name] [coins] [property]",
		Short: "create a new tree",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if cliCtx.ChainID != "sifchain" {
				name := args[0]
				amount := args[1]
				property := args[2]
				coins, err := sdk.ParseCoins(amount)
				if err != nil {
					return err
				}
				// TODO verify the type the tokens that the user can request , Limit it to rowan ?
				signer := cliCtx.GetFromAddress()
				msg := types.NewMsgCreateTree(name, signer, coins, property)
				fmt.Println(msg)
				return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			}
			return nil
		},
	}
	return cmd
}

// TX to add coins from an account to the faucet module account
func GetCmdCreateOrder(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order [price] [tree_id]",
		Short: "create a new limit order",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if cliCtx.ChainID != "sifchain" {
				amount := args[0]
				id := args[1]
				coins, err := sdk.ParseCoins(amount)
				if err != nil {
					return err
				}
				// TODO verify the type the tokens that the user can request , Limit it to rowan ?
				signer := cliCtx.GetFromAddress()
				msg := types.NewMsgBuyTree(signer, id, coins)
				fmt.Println(msg)
				return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			}
			return nil
		},
	}
	return cmd
}
