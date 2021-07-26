package cli

import (
	"context"
	"fmt"
	"github.com/Sifchain/sifnode/x/whitelist/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
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
		GetCmdQueryDenoms(),
	)
	return cmd
}

func GetCmdQueryDenoms() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "denoms",
		Short: "query the complete whitelist",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Entries(context.Background(), &types.QueryEntriesRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.List)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
