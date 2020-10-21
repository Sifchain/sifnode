package main

import (
	"github.com/Sifchain/sifnode/tools/sifgen"
	"github.com/spf13/cobra"
	"strconv"
)

func main() {
	rootCmd := &cobra.Command{Use: "sifgen"}

	_nodeCmd := nodeCmd()
	_nodeCmd.AddCommand(nodeCreateCmd(), nodePromoteCmd(), nodePeerCmd())

	_faucetCmd := faucetCmd()
	_faucetCmd.AddCommand(faucetTransferCmd())

	_networkCmd := networkCmd()
	_networkCmd.AddCommand(networkCreateCmd())

	rootCmd.AddCommand(_networkCmd, _nodeCmd, _faucetCmd)
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
		Use:   "create [chain-id] [node-count] [output-dir] [seed-ip-address]",
		Short: "Create a new network.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			count, _ := strconv.Atoi(args[1])
			sifgen.NewSifgen(args[0]).NetworkCreate(count, args[2], args[3])
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
		Use:   "create [chain-id] [peer-address] [genesis-url]",
		Short: "Create a new node.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				sifgen.NewSifgen(args[0]).NodeCreate( nil)
			} else if len(args) == 2 {
				sifgen.NewSifgen(args[0]).NodeCreate(&args[1])
			}
		},
	}
}

func nodePromoteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promote [chain-id] [moniker] [validator-public-key-address] [key-password] [bond-amount]",
		Short: "Promote the node to full validator.",
		Args:  cobra.MinimumNArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(args[0]).NodePromote(args[1], args[2], args[3], args[4])
		},
	}
}

func nodePeerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-peers [chain-id] [moniker] [[peer-address],...]",
		Short: "Update peers.",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(args[0]).NodePeers(args[1], args[2:])
		},
	}
}

func faucetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "faucet",
		Short: "Faucet operations.",
		Args:  cobra.MinimumNArgs(1),
	}
}

func faucetTransferCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [chain-id] [faucet-password] [faucet-address] [validator-address] [amount]",
		Short: "Transfer coins from the faucet to an account.",
		Args:  cobra.MinimumNArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(args[0]).Transfer(args[1], args[2], args[3], args[4])
		},
	}
}
