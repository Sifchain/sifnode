package cmd

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	v039clp "github.com/Sifchain/sifnode/x/clp/legacy/v39"
	v042clp "github.com/Sifchain/sifnode/x/clp/legacy/v42"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	v039dispensation "github.com/Sifchain/sifnode/x/dispensation/legacy/v39"
	v042dispensation "github.com/Sifchain/sifnode/x/dispensation/legacy/v42"
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	v039ethbridge "github.com/Sifchain/sifnode/x/ethbridge/legacy/v39"
	v042ethbridge "github.com/Sifchain/sifnode/x/ethbridge/legacy/v42"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	v039oracle "github.com/Sifchain/sifnode/x/oracle/legacy/v39"
	v042oracle "github.com/Sifchain/sifnode/x/oracle/legacy/v42"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

var migrationMap = types.MigrationMap{
	"v0.9": Migrate,
}

// GetMigrationCallback returns a MigrationCallback for a given version.
func GetMigrationCallback(version string) types.MigrationCallback {
	return migrationMap[version]
}

// GetMigrationVersions get all migration version in a sorted slice.
func GetMigrationVersions() []string {
	versions := make([]string, len(migrationMap))

	var i int

	for version := range migrationMap {
		versions[i] = version
		i++
	}

	sort.Strings(versions)

	return versions
}

// MigrateGenesisDataCmd returns a command to execute genesis state migration.
func MigrateGenesisDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-data [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate v0.9 /path/to/genesis.json
`, version.AppName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var err error

			target := args[0]
			importGenesis := args[1]

			genDoc, err := tmtypes.GenesisDocFromFile(importGenesis)
			if err != nil {
				return err
			}

			var initialState types.AppMap
			if err := json.Unmarshal(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			migrationFunc := GetMigrationCallback(target)
			if migrationFunc == nil {
				return fmt.Errorf("unknown migration function for version: %s", target)
			}

			// TODO: handler error from migrationFunc call
			newGenState := migrationFunc(initialState, clientCtx)

			genDoc.AppState, err = json.Marshal(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			bz, err := tmjson.Marshal(genDoc)
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			_, err = cmd.OutOrStdout().Write(sortedBz)
			if err != nil {
				return errors.Wrap(err, "failed to write output")
			}

			return nil
		},
	}

	return cmd
}

// Migrate migrates exported state from v0.39 to a v0.40 genesis state.
func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	v039Codec := codec.NewLegacyAmino()
	v039ethbridge.RegisterLegacyAminoCodec(v039Codec)

	v040Codec := clientCtx.JSONMarshaler

	// CLP
	if appState[v039clp.ModuleName] != nil {
		var genesis v039clp.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039clp.ModuleName], &genesis)

		newGenesis := v042clp.Migrate(genesis)
		appState[clptypes.ModuleName] = v040Codec.MustMarshalJSON(&newGenesis)
	}
	// Ethbridge
	if appState[v039ethbridge.ModuleName] != nil {
		var ethbridgeGenesis v039ethbridge.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039ethbridge.ModuleName], &ethbridgeGenesis)

		newGenesis := v042ethbridge.Migrate(ethbridgeGenesis)
		appState[ethbridgetypes.ModuleName] = v040Codec.MustMarshalJSON(newGenesis)
	}
	// Oracle
	if appState[v039oracle.ModuleName] != nil {
		var genesis v039oracle.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039oracle.ModuleName], &genesis)

		newGenesis := v042oracle.Migrate(genesis)
		appState[oracletypes.ModuleName] = v040Codec.MustMarshalJSON(newGenesis)
	}
	// Dispensation
	if appState[v039dispensation.ModuleName] != nil {
		var genesis v039dispensation.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v039dispensation.ModuleName], &genesis)

		newGenesis := v042dispensation.Migrate(genesis)
		appState[dispensationtypes.ModuleName] = v040Codec.MustMarshalJSON(newGenesis)
	}
	// Install evidence module genesis
	if appState[evidencetypes.ModuleName] == nil {
		appState[evidencetypes.ModuleName] = v040Codec.MustMarshalJSON(evidencetypes.DefaultGenesisState())
	}

	return appState
}
