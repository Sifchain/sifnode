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

	rootCmd.AddCommand(_networkCmd)
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
