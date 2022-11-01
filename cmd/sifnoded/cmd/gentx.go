package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

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
			cdc := clientCtx.Codec
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			networkDescriptorNum, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to pass network descriptor: %w", err)
			}
			networkDescriptor := oracletypes.NetworkDescriptor(networkDescriptorNum)
			// check if the networkDescriptor is valid
			if !networkDescriptor.IsValid() {
				return fmt.Errorf("network id: %d is invalid", networkDescriptor)
			}

			addr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return fmt.Errorf("failed to get validator address: %w", err)
			}

			power, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to pass network descriptor: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			oracleGenState := oracletypes.GetGenesisStateFromAppState(cdc, appState)

			networkConfigData := oracletypes.NetworkConfigData{
				NetworkDescriptor:  networkDescriptor,
				ValidatorWhitelist: &oracletypes.ValidatorWhiteList{},
			}

			// find and remove according to network descriptor
			for index := 0; index < len(oracleGenState.NetworkConfigData); index++ {
				if oracleGenState.NetworkConfigData[index].NetworkDescriptor == networkDescriptor {
					networkConfigData = *oracleGenState.NetworkConfigData[index]
					oracleGenState.NetworkConfigData = append(oracleGenState.NetworkConfigData[:index],
						oracleGenState.NetworkConfigData[:index]...)
				}
			}
			found := false

			// just update the power if validator exists
			for index := 0; index < len(networkConfigData.ValidatorWhitelist.ValidatorPower); index++ {
				if bytes.Compare(networkConfigData.ValidatorWhitelist.ValidatorPower[index].ValidatorAddress, addr) == 0 {
					networkConfigData.ValidatorWhitelist.ValidatorPower[index].VotingPower = uint32(power)
					found = true
				}
			}

			// add a new validator if not exists
			if !found {
				newPower := oracletypes.ValidatorPower{
					ValidatorAddress: addr,
					VotingPower:      uint32(power),
				}
				networkConfigData.ValidatorWhitelist.ValidatorPower = append(networkConfigData.ValidatorWhitelist.ValidatorPower, &newPower)
			}

			// add the updated or new whitelist
			oracleGenState.NetworkConfigData = append(oracleGenState.NetworkConfigData, &networkConfigData)

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
