package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Sifchain/sifnode/tools/sifgen"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "sifgen"}

	_networkCmd := networkCmd()
	_networkCmd.PersistentFlags().String("bond-amount", "1000000000000000000000000rowan", "bond amount")
	_networkCmd.PersistentFlags().String("mint-amount", "999000000000000000000000000rowan", "mint amount")
	_networkCmd.AddCommand(networkCreateCmd(), networkResetCmd())

	_nodeCmd := nodeCmd()
	_nodeCreateCmd := nodeCreateCmd()
	_nodeCreateCmd.PersistentFlags().Bool("standalone", false, "standalone node")
	_nodeCreateCmd.PersistentFlags().String("admin-clp-addresses", "", "admin clp addresses")
	_nodeCreateCmd.PersistentFlags().String("admin-oracle-address", "", "admin oracle addresses")
	_nodeCreateCmd.PersistentFlags().String("bind-ip-address", "127.0.0.1", "IPv4 address to bind the node to")
	_nodeCreateCmd.PersistentFlags().String("peer-address", "", "peer node to connect to")
	_nodeCreateCmd.PersistentFlags().String("genesis-url", "", "genesis URL")
	_nodeCreateCmd.PersistentFlags().String("bond-amount", "1000000000000000000000000rowan", "bond amount")
	_nodeCreateCmd.PersistentFlags().String("mint-amount", "999000000000000000000000000rowan", "mint amount")
	_nodeCreateCmd.PersistentFlags().Uint64("min-clp-create-pool-threshold", 100, "minimum CLP create pool threshold")
	_nodeCreateCmd.PersistentFlags().Duration("gov-max-deposit-period", time.Duration(900000000000), "governance max deposit period")
	_nodeCreateCmd.PersistentFlags().Duration("gov-voting-period", time.Duration(900000000000), "governance voting period")
	_nodeCreateCmd.PersistentFlags().String("clp-config-url", "", "URL of the JSON file to use to pre-populate CLPs during genesis")
	_nodeCreateCmd.PersistentFlags().Bool("print-details", false, "print the node details")
	_nodeCreateCmd.PersistentFlags().Bool("with-cosmovisor", false, "setup cosmovisor")
	_nodeCmd.AddCommand(_nodeCreateCmd, nodeResetStateCmd())

	_keyCmd := keyCmd()
	_keyCmd.AddCommand(keyGenerateMnemonicCmd(), keyRecoverFromMnemonicCmd())

	rootCmd.AddCommand(_networkCmd, _nodeCmd, _keyCmd)
	_ = rootCmd.Execute()
}

func networkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "network",
		Short: "Network commands.",
		Args:  cobra.MaximumNArgs(1),
	}
}

func networkCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [chain-id] [validator-count] [output-dir] [seed-ip-address] [output-file]",
		Short: "Create a new network.",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			bondAmount, _ := cmd.Flags().GetString("bond-amount")
			mintAmount, _ := cmd.Flags().GetString("mint-amount")

			count, _ := strconv.Atoi(args[1])
			network := sifgen.NewSifgen(&args[0]).NewNetwork()
			network.BondAmount = bondAmount
			network.MintAmount = mintAmount

			summary, err := network.Build(count, args[2], args[3])
			if err != nil {
				log.Fatal(err)
			}

			if err = ioutil.WriteFile(args[4], []byte(*summary), 0600); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func networkResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset [chain-id] [network-directory]",
		Short: "Reset the state of a network.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(&args[0]).NetworkReset(args[1])
		},
	}
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
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			standalone, _ := cmd.Flags().GetBool("standalone")
			adminCLPAddresses, _ := cmd.Flags().GetString("admin-clp-addresses")
			adminOracleAddress, _ := cmd.Flags().GetString("admin-oracle-address")
			bindIPAddress, _ := cmd.Flags().GetString("bind-ip-address")
			peerAddress, _ := cmd.Flags().GetString("peer-address")
			genesisURL, _ := cmd.Flags().GetString("genesis-url")
			bondAmount, _ := cmd.Flags().GetString("bond-amount")
			mintAmount, _ := cmd.Flags().GetString("mint-amount")
			minCLPCreatePoolThreshold, _ := cmd.Flags().GetUint64("min-clp-create-pool-threshold")
			govMaxDepositPeriod, _ := cmd.Flags().GetDuration("gov-max-deposit-period")
			govVotingPeriod, _ := cmd.Flags().GetDuration("gov-voting-period")
			printDetails, _ := cmd.Flags().GetBool("print-details")
			withCosmovisor, _ := cmd.Flags().GetBool("with-cosmovisor")

			node := sifgen.NewSifgen(&args[0]).NewNode()
			node.Moniker = args[1]
			node.Mnemonic = args[2]

			if standalone {
				node.Standalone = true
				if len(adminCLPAddresses) > 0 {
					node.AdminCLPAddresses = strings.Split(adminCLPAddresses, "|")
				}
				node.AdminOracleAddress = adminOracleAddress
				node.BondAmount = bondAmount
				node.MintAmount = mintAmount
				node.MinCLPCreatePoolThreshold = minCLPCreatePoolThreshold
				node.GovMaxDepositPeriod = govMaxDepositPeriod
				node.GovVotingPeriod = govVotingPeriod
			} else {
				node.PeerAddress = peerAddress
				node.GenesisURL = genesisURL
			}

			node.IPAddr = bindIPAddress
			node.WithCosmovisor = withCosmovisor
			summary, err := node.Build()
			if err != nil {
				log.Fatal(err)
			}

			if printDetails && summary != nil {
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
			sifgen.NewSifgen(nil).KeyGenerateMnemonic("", "")
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
