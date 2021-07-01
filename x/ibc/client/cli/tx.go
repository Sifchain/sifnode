package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        ibchost.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", ibchost.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()

	return txCmd
}
