package cli

import (
	"fmt"
	"reflect"
	"strings"

	// types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/Sifchain/sifnode/x/trees/types"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	faucetQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryTree(queryRoute, cdc),
		)...,
	)

	return faucetQueryCmd

}

// GetCmdFaucet Query to get faucet balance with the specified denom
func GetCmdQueryTree(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tree [id]",
		Short: "Query Trees",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for Trees.%s`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]
			params := types.NewQueryReqGetTreeById(id)
			fmt.Println(reflect.TypeOf(params))
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}
			fmt.Println(reflect.TypeOf(bz))
			fmt.Println(bz)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetTreeById)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var tree types.Tree
			cdc.MustUnmarshalJSON(res, &tree)
			return cliCtx.PrintOutput(tree)
		},
	}
}
