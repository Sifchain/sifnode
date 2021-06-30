package cli

import (
	"context"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetCmdGetEthBridgeProphecy queries information about a specific prophecy
func GetCmdGetEthBridgeProphecy() *cobra.Command {
	return &cobra.Command{
		Use:   `prophecy [prophecy-id]`,
		Short: "Query prophecy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryEthProphecyRequest{
				ProphecyId: []byte(args[0]),
			}

			res, err := queryClient.EthProphecy(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
