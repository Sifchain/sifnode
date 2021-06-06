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
		Use: `prophecy [bridge-registry-contract] [nonce] [symbol] [ethereum-sender]
		--ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address]`,
		Short: "Query prophecy",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			flags := cmd.Flags()

			queryClient := types.NewQueryClient(clientCtx)

			networkID, err := flags.GetInt64(types.FlagEthereumChainID)
			if err != nil {
				return err
			}

			tokenContractString, err := flags.GetString(types.FlagTokenContractAddr)
			if err != nil {
				return err
			}
			tokenContract := types.NewEthereumAddress(tokenContractString)

			bridgeContract := types.NewEthereumAddress(args[0])

			nonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			symbol := args[2]
			ethereumSender := types.NewEthereumAddress(args[3])

			req := &types.QueryEthProphecyRequest{
				NetworkId:             oracletypes.NetworkID(networkID),
				BridgeContractAddress: bridgeContract.String(),
				Nonce:                 int64(nonce),
				Symbol:                symbol,
				TokenContractAddress:  tokenContract.String(),
				EthereumSender:        ethereumSender.String(),
			}

			res, err := queryClient.EthProphecy(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
