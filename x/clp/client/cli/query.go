package cli

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"strings"

	//"github.com/Sifchain/sifnode/x/clp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group clp queries under a subcommand
	clpQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	clpQueryCmd.AddCommand(flags.GetCommands(
		GetCmdPool(queryRoute, cdc),
		GetCmdPools(queryRoute, cdc),
		GetCmdAssets(queryRoute, cdc),
		GetCmdLiquidityProvider(queryRoute, cdc),
	)...)
	return clpQueryCmd
}

func GetCmdPool(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool [External Asset ticker]",
		Short: "Get Details for a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a liquidity pool .
Example:
$ %s pool ETH ROWAN`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ticker := args[0]
			params := types.NewQueryReqGetPool(ticker)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPool)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var pool types.Pool
			cdc.MustUnmarshalJSON(res, &pool)
			out := types.NewPoolResponse(pool, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdPools(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pools",
		Short: "Get all pools",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPools)
			res, height, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			var pools types.Pools
			cdc.MustUnmarshalJSON(res, &pools)
			out := types.NewPoolsResponse(pools, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdAssets(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "assets [lpAddress]",
		Short: "Get all assets for a liquidity provider ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			lpAddressString := args[0]
			lpAddress, err := sdk.AccAddressFromBech32(lpAddressString)
			if err != nil {
				return err
			}
			params := types.NewQueryReqGetAssetList(lpAddress)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAssetList)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var assetList types.Assets
			cdc.MustUnmarshalJSON(res, &assetList)
			out := types.NewAssetListResponse(assetList, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdLiquidityProvider(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "lp [External Asset ticker] [lpAddress]",
		Short: "Get Liquidity Provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a liquidity provioder.
Example:
$ %s pool ETH sif1h2zjknvr3xlpk22q4dnv396ahftzqhyeth7egd`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ticker := args[0]
			lpAddressString := args[1]
			lpAddress, err := sdk.AccAddressFromBech32(lpAddressString)
			if err != nil {
				return err
			}
			params := types.NewQueryReqLiquidityProvider(ticker, lpAddress)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryLiquidityProvider)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var lp types.LiquidityProvider
			cdc.MustUnmarshalJSON(res, &lp)
			out := types.NewLiquidityProviderResponse(lp, height)
			return cliCtx.PrintOutput(out)
		},
	}
}
