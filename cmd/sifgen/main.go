package main

import (
	//"flag"
	//"os"
	//
	//"github.com/Sifchain/sifnode/tools/sifgen"

	"github.com/Sifchain/sifnode/tools/sifgen"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "sifgen"}

	networkLocalnetNodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Node commands.",
		Args:  cobra.MaximumNArgs(1),
	}
	networkLocalnetNodeCreateCmd := &cobra.Command{
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
	networkLocalnetNodePromoteCmd := &cobra.Command{
		Use:   "promote [chain-id] [moniker] [validator-public-key] [key-password] [bond-amount]",
		Short: "Promote the node to full validator.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(args[0]).NodePromote(args[1], args[2], args[3], args[4])
		},
	}
	networkLocalnetNodeCmd.AddCommand(networkLocalnetNodeCreateCmd, networkLocalnetNodePromoteCmd)

	bankCmd := &cobra.Command{
		Use:   "bank",
		Short: "Bank operations.",
		Args:  cobra.MinimumNArgs(1),
	}
	bankTransferCmd := &cobra.Command{
		Use:   "transfer [chain-id] [from-password] [from-address] [to-address] [amount]",
		Short: "Transfer coins from one account to another.",
		Args:  cobra.MinimumNArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			sifgen.NewSifgen(args[0]).Transfer(args[1], args[2], args[3], args[4])
		},
	}
	bankCmd.AddCommand(bankTransferCmd)

	rootCmd.AddCommand(networkLocalnetNodeCmd, bankCmd)
	_ = rootCmd.Execute()
}
