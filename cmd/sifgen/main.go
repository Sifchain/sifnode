package main

import (
	"strconv"

	"github.com/Sifchain/sifnode/tools/sifgen"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "sifgen"}

	_networkCmd := networkCmd()
	_networkCmd.AddCommand(networkCreateCmd(), networkResetCmd())

	_nodeCmd := nodeCmd()
	_nodeCreateCmd := nodeCreateCmd()
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
			count, _ := strconv.Atoi(args[1])
			sifgen.NewSifgen(&args[0]).NetworkCreate(count, args[2], args[3], args[4])
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
		Use:   "create [chain-id] [moniker] [mnemonic] [ip_addr] [peer-address] [genesis-url]",
		Short: "Create a new node.",
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			printDetails, _ := cmd.Flags().GetBool("print-details")
			withCosmovisor, _ := cmd.Flags().GetBool("with-cosmovisor")
			if len(args) == 4 {
				sifgen.NewSifgen(&args[0]).NodeCreate(args[1], args[2], args[3], nil, nil, &printDetails, &withCosmovisor)
			} else {
				sifgen.NewSifgen(&args[0]).NodeCreate(args[1], args[2], args[3], &args[4], &args[5], &printDetails, &withCosmovisor)
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
