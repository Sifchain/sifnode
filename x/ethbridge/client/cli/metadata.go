package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
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

			req := &types.TokenMetadataSearchRequest{
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

// GetCmdSetCrossChainFee is the CLI command to send the message to set crosschain fee for network
func GetCmdAddIBCTokenMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata-add [cosmos-sender-address] [token-name] [token-symbol] [token-address] [token-decimals] [network-descriptor]",
		Short: "Used to manually add Token Metadata for IBC tokens.",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			tokenName := args[1]

			tokenSymbol := args[2]

			tokenAddressRaw := args[3]
			if !common.IsHexAddress(tokenAddressRaw) {
				return errors.New("Error parsing tokenAddress invalid format must be a hex address")
			}

			tokenAddress := types.NewEthereumAddress(tokenAddressRaw)

			tokenDecimalsRaw, err := strconv.ParseInt(args[4], 10, 32)
			if err != nil {
				return errors.New("Error parsing token decimals, must be base 10 number")
			}
			tokenDecimals := int32(tokenDecimalsRaw)

			networkDescriptorRaw, err := strconv.Atoi(args[5])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}
			networkDescriptor := oracletypes.NetworkDescriptor(networkDescriptorRaw)

			msg := types.NewTokenMetadata(cosmosSender, tokenName, tokenSymbol, tokenDecimals, tokenAddress, networkDescriptor)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSetCrossChainFee is the CLI command to send the message to set crosschain fee for network
func GetCmdDeleteIBCTokenMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata-delete [cosmos-sender-address] [denom-hash]",
		Short: "Used to manually delete Token Metadata for IBC tokens.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			denomHash := args[1]

			msg := types.DeleteTokenMetadata(cosmosSender, denomHash)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
