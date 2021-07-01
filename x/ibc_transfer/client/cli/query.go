package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	ibc_transfer_types "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	// Group dispensation queries under a subcommand
	queryCmd := &cobra.Command{
		Use:                        ibc_transfer_types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", ibc_transfer_types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand()
	return queryCmd
}
