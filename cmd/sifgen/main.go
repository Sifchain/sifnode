package main

import (
	"github.com/Sifchain/sifnode/tools/sifgen"
	"github.com/spf13/cobra"
	"strconv"
)

func main() {
	rootCmd := &cobra.Command{Use: "sifgen"}

	_networkCmd := networkCmd()
	_networkCmd.AddCommand(networkCreateCmd())

	_nodeCmd := nodeCmd()
	_nodeCmd.AddCommand(nodeCreateCmd())

	rootCmd.AddCommand(_networkCmd, _nodeCmd)
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
			sifgen.NewSifgen(args[0]).NetworkCreate(count, args[2], args[3], args[4])
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
				sifgen.NewSifgen(args[0]).NodeCreate(nil, nil)
			} else {
				sifgen.NewSifgen(args[0]).NodeCreate(&args[1], &args[2])
			}
		},
	}
}
