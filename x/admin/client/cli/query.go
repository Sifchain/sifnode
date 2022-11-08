package cli

import (
	"context"
	"fmt"

	"github.com/Sifchain/sifnode/x/admin/types"
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
	cmd.AddCommand(GetCmdAccounts(), GetCmdParams())
	return cmd
}

func GetCmdAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "query registered accounts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ListAccounts(context.Background(), &types.ListAccountsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(clientCtx.Codec.MustMarshalJSON(res))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "query params",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetParams(context.Background(), &types.GetParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(clientCtx.Codec.MustMarshalJSON(res))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
