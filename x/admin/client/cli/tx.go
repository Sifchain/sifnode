package cli

import (
	"github.com/Sifchain/sifnode/x/admin/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Admin key management transactions sub-commands",
	}
	cmd.AddCommand()
	return cmd
}
