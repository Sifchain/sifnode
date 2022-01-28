package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

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
	cmd.AddCommand(GetCmdQueryPositionsForAddress())
	return cmd
}

func GetCmdQueryPositionsForAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "positions-for-address [address]",
		Short: "query positions for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetPositionsForAddress(context.Background(), &types.PositionsForAddressRequest{
				Address: addr.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintBytes(clientCtx.Codec.MustMarshalJSON(res))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
