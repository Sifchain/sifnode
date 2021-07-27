package cli

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	whiteListcmd := &cobra.Command{
		Use:   "whitelist",
		Short: "Token whitelist transactions subcommands",
	}

	whiteListcmd.AddCommand(
		GetCmdUpdate(),
	)
	return whiteListcmd
}

func GetCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [Denom] [Decimals]",
		Short: "Add new deon to the whitelist",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(2)(cmd, args)
			if err != nil {
				return err
			}
			decimals, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			msg := types.MsgRegister{
				From:     clientCtx.GetFromAddress().String(),
				Denom:    args[0],
				Decimals: int64(decimals),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}
