package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

func AddGenesisValidatorCmd(defaultNodeHome string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-genesis-validators [network_descriptor] [address_or_key_name] [power]",
		Short: "add genesis validators to genesis.json",
		Long: `add validator to genesis.json. The provided account must specify
the account address or key name. If a key name is given, the address will be looked up in the local Keybase. 
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			networkDescriptor, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to pass network descriptor: %w", err)
			}
			// check if the networkDescriptor is valid
			if !oracletypes.NetworkDescriptor(networkDescriptor).IsValid() {
				return fmt.Errorf("network id: %d is invalid", networkDescriptor)
			}

			addr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return fmt.Errorf("failed to get validator address: %w", err)
			}

			power, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to pass network descriptor: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			oracleGenState := oracletypes.GetGenesisStateFromAppState(cdc, appState)
			if oracleGenState.AddressWhitelist == nil {
				oracleGenState.AddressWhitelist = make(map[uint32]*oracletypes.ValidatorWhiteList)
			}

			_, ok := oracleGenState.AddressWhitelist[uint32(networkDescriptor)]

			if !ok {
				oracleGenState.AddressWhitelist[uint32(networkDescriptor)] = &oracletypes.ValidatorWhiteList{WhiteList: make(map[string]uint32)}
			}

			whiteList := oracleGenState.AddressWhitelist[uint32(networkDescriptor)].WhiteList
			whiteList[addr.String()] = uint32(power)

			oracleGenStateBz, err := json.Marshal(oracleGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[oracletypes.ModuleName] = oracleGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")

	return cmd
}
