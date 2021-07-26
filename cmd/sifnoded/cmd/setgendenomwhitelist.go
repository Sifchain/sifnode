package cmd

import (
	"encoding/json"
	"fmt"
	whitelisttypes "github.com/Sifchain/sifnode/x/whitelist/types"
	whitelistutils "github.com/Sifchain/sifnode/x/whitelist/utils"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

func SetGenesisDenomWhitelist(defaultNodeHome string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set-gen-denom-whitelist [path to json file]",
		Short: "Add a list of denoms to the whitelist",
		Long:  `Add a list of denoms to the whitelist , this list can only be edited by the admin in future`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)
			genFile := config.GenesisFile()
			// Get input list
			whitelist, err := whitelistutils.ParseDenoms(args[0])
			if err != nil {
				return err
			}
			// Get Existing List
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			whitelistGenState := whitelisttypes.GetGenesisStateFromAppState(cdc, appState)
			// TODO :Append New Entries to existing list
			//whitelistGenState.Registry.Entries = append(whitelistGenState.Registry.Entries, whitelist.Entries...)
			whitelistGenState.Registry = &whitelist
			whitelistGenStateBz, err := json.Marshal(whitelistGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[whitelisttypes.ModuleName] = whitelistGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	return cmd
}
