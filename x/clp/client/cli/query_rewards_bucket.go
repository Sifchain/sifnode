package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func GetCmdListRewardsBucket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-rewards-bucket",
		Short: "list all rewards-bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.AllRewardsBucketReq{
				Pagination: pageReq,
			}

			res, err := queryClient.GetRewardsBucketAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdShowRewardsBucket() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-rewards-bucket [denom]",
		Short: "shows a rewards-bucket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			argDenom := args[0]

			params := &types.RewardsBucketReq{
				Denom: argDenom,
			}

			res, err := queryClient.GetRewardsBucket(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
