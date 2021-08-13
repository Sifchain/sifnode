package client

import (
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/Sifchain/sifnode/x/ethbridge/client/cli"
	"github.com/Sifchain/sifnode/x/ethbridge/client/rest"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group ethbridge queries under a subcommand
	ethBridgeQueryCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "Querying commands for the ethbridge module",
	}

	ethBridgeQueryCmd.PersistentFlags().String(types.FlagEthereumChainID, "", "Ethereum chain ID")
	ethBridgeQueryCmd.PersistentFlags().String(types.FlagTokenContractAddr, "", "Token address representing a unique asset type")

	flags.AddQueryFlagsToCmd(ethBridgeQueryCmd)

	ethBridgeQueryCmd.AddCommand(
		cli.GetCmdGetEthBridgeProphecy(),
		cli.GetCmdGetTokenMetadata(),
	)

	return ethBridgeQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	ethBridgeTxCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "EthBridge transactions subcommands",
	}

	ethBridgeTxCmd.PersistentFlags().String(types.FlagEthereumChainID, "", "Ethereum chain ID")
	ethBridgeTxCmd.PersistentFlags().String(types.FlagTokenContractAddr, "", "Token address representing a unique asset type")

	flags.AddTxFlagsToCmd(ethBridgeTxCmd)

	ethBridgeTxCmd.AddCommand(
		cli.GetCmdCreateEthBridgeClaim(),
		cli.GetCmdBurn(),
		cli.GetCmdLock(),
		cli.GetCmdUpdateWhiteListValidator(),
		cli.GetCmdUpdateCrossChainFeeReceiverAccount(),
		cli.GetCmdRescueCrossChainFee(),
		cli.GetCmdSetCrossChainFee(),
	)

	return ethBridgeTxCmd
}

// RegisterRESTRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx sdkclient.Context, r *mux.Router, storeName string) {
	rest.RegisterRESTRoutes(cliCtx, r, storeName)
}
