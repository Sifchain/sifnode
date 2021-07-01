package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"

	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	// Group dispensation queries under a subcommand
	queryCmd := &cobra.Command{
		Use:                        ibchost.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", ibchost.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand()
	return queryCmd
}
