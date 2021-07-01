package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	ibc_transfer_types "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        ibc_transfer_types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", ibc_transfer_types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()

	return txCmd
}
