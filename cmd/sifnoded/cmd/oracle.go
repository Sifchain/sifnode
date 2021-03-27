package cmd

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/gogo/protobuf/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/Sifchain/sifnode/x/oracle"
)

// SetGenesisOracleAdminCmd set the admin address can update the whitelist validators
func SetGenesisOracleAdminCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set-genesis-oracle-admin [address_or_key_name]",
		Short: "Set a genesis account to genesis.json",
		Long: `Set a genesis account to genesis.json. The account has the privilege to update whitelist validators.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			inBuf := bufio.NewReader(cmd.InOrStdin())
			if err != nil {
				// attempt to lookup address from Keybase if no address was provided
				kb, err := keys.NewKeyring(
					sdk.KeyringServiceName(),
					viper.GetString(flags.FlagKeyringBackend),
					viper.GetString(flagClientHome),
					inBuf,
				)
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			oracleGenState := oracle.GetGenesisStateFromAppState(cdc, appState)
			oracleGenState.AdminAddress = addr

			oracleGenStateBz, err := cdc.MarshalJSON(oracleGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[oracle.ModuleName] = oracleGenStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")

	return cmd
}
