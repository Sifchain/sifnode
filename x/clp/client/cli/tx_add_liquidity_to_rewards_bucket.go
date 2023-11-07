package cli

import (
	"strconv"

	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAddLiquidityToRewardsBucket() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add-liquidity-to-rewards-bucket [amount]",
		Example: "sifnodecli tx clp add-liquidity-to-rewards-bucket 100000000000000000000000000rowan,100000000000000000000000000ceth,100000000cusdc",
		Short:   "To add liquidity to the rewards bucket",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddLiquidityToRewardsBucketRequest(
				clientCtx.GetFromAddress().String(),
				argAmount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
