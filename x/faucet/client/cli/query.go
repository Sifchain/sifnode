package cli

import (
	"fmt"
	"strings"

	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/Sifchain/sifnode/x/faucet/types"
)

func GetQueryCmd(queryRoute string) *cobra.Command {
	faucetQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdFaucet(queryRoute, cdc),
		)...,
	)

	return faucetQueryCmd

}

// GetCmdFaucet Query to get faucet balance with the specified denom
func GetCmdFaucet(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "balance",
		Short: "Get Faucet Balances",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for faucet balance.%s`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			if cliCtx.ChainID != "sifchain" {
				route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBalance)
				res, _, err := cliCtx.Query(route)
				if err != nil {
					return err
				}
				var coins types2.Coins
				cdc.MustUnmarshalJSON(res, &coins)
				return cliCtx.PrintOutput(coins)
			}
			return nil
		},
	}
}
