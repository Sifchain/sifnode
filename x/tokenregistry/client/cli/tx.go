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
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Token registry transactions sub-commands",
	}

	cmd.AddCommand(
		GetCmdRegister(),
	)
	return cmd
}

func GetCmdRegister() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [Denom] [Decimals]",
		Short: "Add / update token on the registry",
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
