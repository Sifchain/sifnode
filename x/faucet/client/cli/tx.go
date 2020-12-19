package cli

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/"

	"github.com/mossid/sdk-nameservice-example/x/faucet"
)

const (
	flagAmount = "amount"
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
	// this line is used by starport scaffolding # 1
	// TODO: Add tx based commands
	// GetCmd<Action>(cdc)
	)...)

	return faucetTxCmd
}

func GetCmdRequestCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-coins [amount]",
		Short: "request coins from faucet",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := faucet.MsgRequestCoins{
				Coins:     coins,
				Requester: account,
			}

			return completeAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagAmount, "", "Amount of coins to request")
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}
