package cmd

import (
	"bytes"
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

func AddGenesisValidatorCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-genesis-validators [address_or_key_name]",
		Short: "add genesis validators to genesis.json",
		Long: `add validator to genesis.json. The provided account must specify
the account address or key name. If a key name is given, the address will be looked up in the local Keybase. 
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to get validator address: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			oracleGenState := oracle.GetGenesisStateFromAppState(cdc, appState)

			for _, item := range oracleGenState.AddressWhitelist {
				if bytes.Equal(item, addr) {
					return fmt.Errorf("address %s already in white list", addr)
				}
			}
			oracleGenState.AddressWhitelist = append(oracleGenState.AddressWhitelist, addr)

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
