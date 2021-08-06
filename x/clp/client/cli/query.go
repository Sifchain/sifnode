package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	//"github.com/Sifchain/sifnode/x/clp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group clp queries under a subcommand
	clpQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	clpQueryCmd.AddCommand(
		GetCmdPool(queryRoute),
		GetCmdPools(queryRoute),
		GetCmdAssets(queryRoute),
		GetCmdLiquidityProvider(queryRoute),
		GetCmdLpList(queryRoute),
		GetCmdAllLps(queryRoute),
	)
	return clpQueryCmd
}

func GetCmdPool(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [External Asset symbol]",
		Short: "Get Details for a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a liquidity pool.
Example:
$ %s pool ETH ROWAN`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			ticker := args[0]
			params := types.NewQueryReqGetPool(ticker)

			result, err := queryClient.GetPool(context.Background(), &params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdPools(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "Get all pools",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			result, err := queryClient.GetPools(context.Background(), &types.PoolsReq{
				Pagination: pageReq,
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pools")

	return cmd
}

func GetCmdAssets(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets [lpAddress]",
		Short: "Get all assets for a liquidity provider ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			lpAddress := args[1]

			assetReq := types.AssetListReq{
				LpAddress: lpAddress,
			}

			res, err := queryClient.GetAssetList(context.Background(), &assetReq)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "assets")

	return cmd
}

func GetCmdLiquidityProvider(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lp [External Asset symbol] [lpAddress]",
		Short: "Get Liquidity Provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a liquidity provioder.
Example:
$ %s pool ETH sif1h2zjknvr3xlpk22q4dnv396ahftzqhyeth7egd`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			symbol := args[0]
			lpAddress := args[1]

			lpReq := types.LiquidityProviderReq{
				Symbol:    symbol,
				LpAddress: lpAddress,
			}

			res, err := queryClient.GetLiquidityProvider(context.Background(), &lpReq)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdLpList(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lplist [symbol]",
		Short: "Get all liquidity providers for the asset ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			assetSymbol := args[0]
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			result, err := queryClient.GetLiquidityProviderList(context.Background(), &types.LiquidityProviderListReq{
				Symbol:     assetSymbol,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "lplist")

	return cmd
}

func GetCmdAllLps(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-lp",
		Short: "Get all liquidity providers on sifnode ",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			result, err := queryClient.GetLiquidityProviders(context.Background(), &types.LiquidityProvidersReq{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "liquidityProviders")

	return cmd
}
