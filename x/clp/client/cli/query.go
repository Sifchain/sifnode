package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	return &cobra.Command{
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

			result, err := queryClient.QueryGetPool(context.Background(), &params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}
}

func GetCmdPools(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "pools",
		Short: "Get all pools",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPools)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			var pools types.PoolsResponse
			cdc.MustUnmarshalJSON(res, &pools)
			return cliCtx.PrintOutput(pools)
		},
	}
}

func GetCmdAssets(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "assets [lpAddress]",
		Short: "Get all assets for a liquidity provider ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			lpAddressString := args[0]
			lpAddress, err := sdk.AccAddressFromBech32(lpAddressString)
			if err != nil {
				return err
			}
			params := types.NewQueryReqGetAssetList(lpAddress)
			bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAssetList)
			res, height, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var assetList types.Assets
			cdc.MustUnmarshalJSON(res, &assetList)
			out := types.NewAssetListResponse(assetList, height)

			return clientCtx.PrintProto(out)
		},
	}
}

func GetCmdLiquidityProvider(queryRoute string) *cobra.Command {
	return &cobra.Command{
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

			res, err := queryClient.LiquidityProvider(context.Background(), &lpReq)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}
}

func GetCmdLpList(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "lplist [symbol]",
		Short: "Get all liquidity providers for the asset ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			assetSymbol := args[0]
			params := types.NewQueryReqGetLiquidityProviderList(assetSymbol)
			bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryLPList)
			res, height, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var assetList types.LiquidityProviders
			cdc.MustUnmarshalJSON(res, &assetList)
			out := types.NewLpListResponse(assetList, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdAllLps(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "all-lp",
		Short: "Get all liquidity providers on sifnode ",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAllLP)
			res, height, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			var lps types.LiquidityProviders
			cdc.MustUnmarshalJSON(res, &lps)
			out := types.NewLpListResponse(lps, height)
			return clientCtx.PrintProto(out)
		},
	}
}
