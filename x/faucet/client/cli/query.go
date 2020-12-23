package cli

import (
	"fmt"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/Sifchain/sifnode/x/faucet/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group faucet queries under a subcommand
	faucetQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetQueryCmd.AddCommand(
		flags.GetCommands(
			// this line is used by starport scaffolding # 1
			GetCmdFaucet(queryRoute, cdc),
			GetCmdFaucetAddress(queryRoute, cdc),
		)...,
	)

	return faucetQueryCmd
}

// Query to get faucet balance with the specified denom
func GetCmdFaucet(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "balance",
		Short: "Get Faucet Balances",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for faucet balance.`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBalance)
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}
			var coins types2.Coins
			cdc.MustUnmarshalJSON(res, &coins)
			return cliCtx.PrintOutput(coins)

		},
	}
}

// Query to get faucet module address
func GetCmdFaucetAddress(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "address",
		Short: "Get Faucet Address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query address for faucet.`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			return cliCtx.PrintOutput(types.GetFaucetModuleAddress())
		},
	}
}
