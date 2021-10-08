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
func GetEthereumLockBurnNonce() *cobra.Command {
	return &cobra.Command{
		Use:   `ethereum-lock-burn-nonce [network-descriptor] [val-address]`,
		Short: "Query ethereum-lock-burn-nonce",
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

			req := &types.QueryEthereumLockBurnNonceRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				RelayerValAddress: args[1],
			}

			res, err := queryClient.EthereumLockBurnNonce(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetWitnessLockBurnNonce queries lock burn nonce for a relayer
func GetWitnessLockBurnNonce() *cobra.Command {
	return &cobra.Command{
		Use:   `witness-lock-burn-nonce [network-descriptor] [val-address]`,
		Short: "Query witness-lock-burn-nonce",
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

			req := &types.QueryWitnessLockBurnNonceRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				RelayerValAddress: args[1],
			}

			res, err := queryClient.WitnessLockBurnNonce(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetGlocalNonceBlockNumber queries block number for global nonce
func GetGlocalNonceBlockNumber() *cobra.Command {
	return &cobra.Command{
		Use:   `global-nonce-block-number [network-descriptor] [global-nonce]`,
		Short: "Query global-nonce-block-number",
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

			globalNonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryGlocalNonceBlockNumberRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				GlobalNonce:       uint64(globalNonce),
			}

			res, err := queryClient.GlocalNonceBlockNumber(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
