package cli

import (
	"context"
	"strconv"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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

			req := &types.QueryEthereumLockBurnSequenceRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				RelayerValAddress: args[1],
			}

			res, err := queryClient.EthereumLockBurnSequence(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetWitnessLockBurnSequence queries lock burn nonce for a relayer
func GetWitnessLockBurnSequence() *cobra.Command {
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

			req := &types.QueryWitnessLockBurnSequenceRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				RelayerValAddress: args[1],
			}

			res, err := queryClient.WitnessLockBurnSequence(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetGlobalSequenceBlockNumber queries block number for global nonce
func GetGlobalSequenceBlockNumber() *cobra.Command {
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

			globalSequence, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryGlobalSequenceBlockNumberRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				GlobalSequence:    uint64(globalSequence),
			}

			res, err := queryClient.GlobalSequenceBlockNumber(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// GetProphecyCompleted queries block number for global nonce
func GetProphecyCompleted() *cobra.Command {
	return &cobra.Command{
		Use:   `prophecy-completed [network-descriptor] [global-nonce]`,
		Short: "Query prophecy-completed",
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

			globalSequence, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryPropheciesCompletedRequest{
				NetworkDescriptor: oracletypes.NetworkDescriptor(networkDescriptor),
				GlobalSequence:    uint64(globalSequence),
			}

			res, err := queryClient.PropheciesCompleted(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func GetCmdGetBlacklist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blacklist",
		Short: "Query full address blacklist",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryBlacklistRequest{}

			res, err := queryClient.GetBlacklist(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
