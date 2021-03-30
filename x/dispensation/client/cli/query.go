package cli

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group dispensation queries under a subcommand
	dispensationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	dispensationQueryCmd.AddCommand(flags.GetCommands()...)
	return dispensationQueryCmd
}
