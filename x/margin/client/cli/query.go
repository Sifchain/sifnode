package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/x/margin/chain"
	"github.com/Sifchain/sifnode/x/margin/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetMarginParamsCmd(),
		GetPoolsCmd(),
	)
	return cmd
}

func GetMarginParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			marginParams, err := chain.QueryMarginParams(clientCtx)
			if err != nil {
				return err
			}
			return marginParams
		},
	}
	return cmd
}

func GetPoolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pools, err := chain.QueryPools(clientCtx)
			if err != nil {
				return err
			}
			return pools
		},
	}
	return cmd
}

func GetEventCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			flags := cmd.Flags()

			baseUrl, err := flags.GetString(FlagBaseUrl)
			if err != nil {
				return err
			}

			blockEvents, err := chain.BlockEvents(clientCtx)
			if err != nil {
				return err
			}
			return blockEvents
		},
	}
	return cmd
}
