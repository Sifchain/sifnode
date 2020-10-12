package cli

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"

	//"github.com/Sifchain/sifnode/x/clp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group clp queries under a subcommand
	clpQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	return clpQueryCmd
}

//	clpQueryCmd.AddCommand(
//		flags.GetCommands(
//		// TODO: Add query Cmds
//		)...,
//	)
//
//	return clpQueryCmd
//}

// TODO: Add Query Commands
