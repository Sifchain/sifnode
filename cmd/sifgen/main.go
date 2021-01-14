package main

import (
	"fmt"
	"strings"

	"github.com/Sifchain/sifnode/tools/sifgen"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "sifgen"}

	_nodeCmd := nodeCmd()
	_nodeCreateCmd := nodeCreateCmd()
	_nodeCreateCmd.PersistentFlags().Bool("standalone", false, "standalone node")
	_nodeCreateCmd.PersistentFlags().String("admin-clp-addresses", "", "admin clp addresses")
	_nodeCreateCmd.PersistentFlags().String("admin-oracle-address", "", "admin oracle addresses")
	_nodeCreateCmd.PersistentFlags().String("bind-ip-address", "127.0.0.1", "IPv4 address to bind the node to")
	_nodeCreateCmd.PersistentFlags().String("peer-address", "", "peer node to connect to")
	_nodeCreateCmd.PersistentFlags().String("genesis-url", "", "genesis URL")
	_nodeCreateCmd.PersistentFlags().String("bond-amount", "100000000000000000rowan", "bond amount")
	_nodeCreateCmd.PersistentFlags().String("mint-amount", "1000000000000000000000000000rowan", "bond amount")
	_nodeCreateCmd.PersistentFlags().Bool("print-details", false, "print the node details")
	_nodeCreateCmd.PersistentFlags().Bool("with-cosmovisor", false, "setup cosmovisor")
	_nodeCmd.AddCommand(_nodeCreateCmd, nodeResetStateCmd())

	_keyCmd := keyCmd()
	_keyCmd.AddCommand(keyGenerateMnemonicCmd(), keyRecoverFromMnemonicCmd())

	rootCmd.AddCommand(_nodeCmd, _keyCmd)
	_ = rootCmd.Execute()
}

func nodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "node",
		Short: "Node commands.",
		Args:  cobra.MaximumNArgs(1),
	}
}

func nodeCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [chain-id] [moniker] [mnemonic]",
		Short: "Create a new node.",
		Args:  cobra.MinimumNArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			standalone, _ := cmd.Flags().GetBool("standalone")
			adminCLPAddresses, _ := cmd.Flags().GetString("admin-clp-addresses")
			adminOracleAddress, _ := cmd.Flags().GetString("admin-oracle-address")
			bindIPAddress, _ := cmd.Flags().GetString("bind-ip-address")
			peerAddress, _ := cmd.Flags().GetString("peer-address")
			genesisURL, _ := cmd.Flags().GetString("genesis-url")
			bondAmount, _ := cmd.Flags().GetString("bond-amount")
			mintAmount, _ := cmd.Flags().GetString("mint-amount")
			printDetails, _ := cmd.Flags().GetBool("print-details")
			withCosmovisor, _ := cmd.Flags().GetBool("with-cosmovisor")

			node := sifgen.NewSifgen(&args[0]).NewNode()
			node.Moniker = args[2]
			node.Mnemonic = args[3]

			if standalone {
				node.Standalone = true
				node.AdminCLPAddresses = strings.Split(adminCLPAddresses, ",")
				node.AdminOracleAddress = adminOracleAddress
				node.IPAddr = bindIPAddress
				node.BondAmount = bondAmount
				node.MintAmount = mintAmount
			} else {
				node.PeerAddress = peerAddress
				node.GenesisURL = genesisURL
			}

			node.WithCosmovisor = withCosmovisor
			summary, err := node.Build()
			if err != nil {
				panic(err)
			}

			if printDetails {
				fmt.Println(*summary)
			}
		},
	}
}

func nodeResetStateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset [chain-id] [node-directory]",
		Short: "Reset the state of a node.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				sifgen.NewSifgen(&args[0]).NodeReset(nil)
			} else {
				sifgen.NewSifgen(&args[0]).NodeReset(&args[1])
			}
		},
	}
}

func keyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "key",
		Short: "Key commands.",
		Args:  cobra.MaximumNArgs(0),
	}
}

func keyGenerateMnemonicCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate a mnemonic phrase.",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(nil).KeyGenerateMnemonic(nil, nil)
		},
	}
}

func keyRecoverFromMnemonicCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recover [mnemonic]",
		Short: "Recover your key details from your mnemonic phrase.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(nil).KeyRecoverFromMnemonic(args[0])
		},
	}
}
