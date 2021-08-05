package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// GetCmdGetEthBridgeProphecy queries information about a specific prophecy
func GetCmdGetTokenMetadata() *cobra.Command {
	return &cobra.Command{
		Use:   `metadata [denom-hash]`,
		Short: "Query token metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			metadataClient := types.NewTokenMetadataServiceClient(clientCtx)

			req := &types.TokenMetadataRequest{
				Denom: args[0],
			}

			res, err := metadataClient.Search(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Metadata)
		},
	}
}
