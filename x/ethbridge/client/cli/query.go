package cli

import (
	"context"
	"strconv"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

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

// GetCmdGetCrosschainFeeConfig queries crosschain fee config for a network
func GetCmdGetCrosschainFeeConfig() *cobra.Command {
	return &cobra.Command{
		Use:   `crosschain-fee-config [network-descriptor]`,
		Short: "Query crosschain-fee-config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			networkDescriptor, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCrosschainFeeConfigRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
			}

			res, err := queryClient.CrosschainFeeConfig(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetLockBurnNonce queries lock burn nonce for a relayer
func GetLockBurnNonce() *cobra.Command {
	return &cobra.Command{
		Use:   `lock-burn-nonce [network-descriptor] [val-address]`,
		Short: "Query lock-burn-nonce",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			networkDescriptor, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryLockBurnNonceRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				RelayerValAddress: args[1],
			}

			res, err := queryClient.LockBurnNonce(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
