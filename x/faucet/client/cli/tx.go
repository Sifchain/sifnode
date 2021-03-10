package cli

import (
	"bufio"
	"fmt"

	"github.com/Sifchain/sifnode/x/faucet/types"
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
		GetCmdRequestCoins(cdc),
		GetCmdAddCoins(cdc))...)

	return faucetTxCmd
}

// TX to request coins from faucet module account to the requesters account
func GetCmdRequestCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-coins [amount]",
		Short: "request coins from faucet ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if cliCtx.ChainID != "sifchain" {
				amount := args[0]
				coins, err := sdk.ParseCoins(amount)
				if err != nil {
					return err
				}
				// TODO verify the type the tokens that the user can request , Limit it to rowan ?
				signer := cliCtx.GetFromAddress()
				msg := types.NewMsgRequestCoins(signer, coins)
				return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			}
			return nil
		},
	}
	return cmd
}

// TX to add coins from an account to the faucet module account
func GetCmdAddCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-coins [amount]",
		Short: "add coins to faucet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc)).WithGas(0)
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			if cliCtx.ChainID != "sifchain" {
				amount := args[0]
				coins, err := sdk.ParseCoins(amount)
				if err != nil {
					return err
				}
				signer := cliCtx.GetFromAddress()
				msg := types.NewMsgAddCoins(signer, coins)
				return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			}
			return nil
		},
	}
	return cmd
}

//sif17s95c5jpc6x2l3edwh4dm8yhac68yru7a7kr3x
//sif1d0q6q6afvzk3whu6cy338yth64vau78r43fsuh

//sif1d0q6q6afvzk3whu6cy338yth64vau78r43fsuh
//sif1r8x7m7u8ttwn6ue4z2krstrp62qdvcw2myg0xz
//sif1rlt8svkn4h4vyqmklgzfgcvfeamdce2cwxntsq
//sif18xt24a5psswl6a838ngtr2x29qp8626rkd32uf
