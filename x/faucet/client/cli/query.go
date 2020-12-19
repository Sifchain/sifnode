package cli

import (
	"fmt"
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
		)...,
	)

	return faucetQueryCmd
}

func GetCmdFaucet(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "faucet-balance [denom]",
		Short: "Get account details for faucet",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for faucet balance.`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBalance)
			denom := args[0]
			faucetBalanceRequest := types.NewQueryReqGetFaucetBalance(denom)
			bz, err := cliCtx.Codec.MarshalJSON(faucetBalanceRequest)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(res)
		},
	}
}
