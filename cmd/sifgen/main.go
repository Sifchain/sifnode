package main

import (
	"github.com/Sifchain/sifnode/tools/sifgen"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "sifgen"}

	_nodeCmd := nodeCmd()
	_nodeCmd.AddCommand(nodeCreateCmd(), nodePromoteCmd())

	_faucetCmd := faucetCmd()
	_faucetCmd.AddCommand(faucetTransferCmd())

	rootCmd.AddCommand(_nodeCmd, _faucetCmd)
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
		Use:   "create [chain-id] [peer-address] [genesis-url]",
		Short: "Create a new node.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				sifgen.NewSifgen(args[0]).NodeCreate(nil, nil)
			} else if len(args) == 3 {
				sifgen.NewSifgen(args[0]).NodeCreate(&args[1], &args[2])
			}
		},
	}
}

func nodePromoteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promote [chain-id] [moniker] [validator-public-key] [key-password] [bond-amount]",
		Short: "Promote the node to full validator.",
		Args:  cobra.MaximumNArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(args[0]).NodePromote(args[1], args[2], args[3], args[4])
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
		Use:   "transfer [chain-id] [faucet-password] [faucet-address] [to-address] [amount]",
		Short: "Transfer coins from the faucet to an account.",
		Args:  cobra.MinimumNArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(args[0]).Transfer(args[1], args[2], args[3], args[4])
		},
	}
}
