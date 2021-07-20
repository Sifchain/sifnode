package client

import "github.com/spf13/cobra"

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist",
		Short: "Token whitelist transactions subcommands",
	}

	return cmd
}
