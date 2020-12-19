package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	faucetTxCmd.AddCommand(flags.PostCommands(GetCmdRequestCoins(cdc))...)

	return faucetTxCmd
}

func GetCmdRequestCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-coins [amount]",
		Short: "request coins from faucet",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			signer := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := types.MsgRequestCoins{
				Coins:     coins,
				Requester: signer,
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx,txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagAmount, "", "Amount of coins to request")
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}
