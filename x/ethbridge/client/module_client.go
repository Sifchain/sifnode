package client

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/Sifchain/sifnode/x/ethbridge/client/cli"
	"github.com/Sifchain/sifnode/x/ethbridge/client/rest"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	// Group ethbridge queries under a subcommand
	ethBridgeQueryCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "Querying commands for the ethbridge module",
	}

	ethBridgeQueryCmd.PersistentFlags().String(types.FlagEthereumChainID, "", "Ethereum chain ID")
	ethBridgeQueryCmd.PersistentFlags().String(types.FlagTokenContractAddr, "", "Token address representing a unique asset type")

	ethBridgeQueryCmd.AddCommand(flags.GetCommands(
		cli.GetCmdGetEthBridgeProphecy(storeKey, cdc),
	)...)

	return ethBridgeQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ethBridgeTxCmd := &cobra.Command{
		Use:   "ethbridge",
		Short: "EthBridge transactions subcommands",
	}

	ethBridgeTxCmd.PersistentFlags().String(types.FlagEthereumChainID, "", "Ethereum chain ID")
	ethBridgeTxCmd.PersistentFlags().String(types.FlagTokenContractAddr, "", "Token address representing a unique asset type")

	ethBridgeTxCmd.AddCommand(flags.PostCommands(
		cli.GetCmdCreateEthBridgeClaim(cdc),
		cli.GetCmdBurn(cdc),
		cli.GetCmdLock(cdc),
		cli.GetCmdUpdateWhiteListValidator(cdc),
		cli.GetCmdUpdateCethReceiverAccount(cdc),
		cli.GetCmdRescueCeth(cdc),
		cli.GetCmdUpdateGasPrice(cdc),
	)...)

	return ethBridgeTxCmd
}

// RegisterRESTRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	rest.RegisterRESTRoutes(cliCtx, r, storeName)
}
