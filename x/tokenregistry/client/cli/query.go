package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	whitelistutils "github.com/Sifchain/sifnode/x/tokenregistry/utils"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmdQueryEntries(),
		GetCmdGenerateEntry(),
		GetCmdAddEntry(),
		GetCmdAddAllEntries(),
	)
	return cmd
}

func GetCmdQueryEntries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entries",
		Short: "query the complete token registry",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Entries(context.Background(), &types.QueryEntriesRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(clientCtx.Codec.MustMarshalJSON(res.Registry))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdGenerateEntry() *cobra.Command {
	var flagDenom = "token_denom"
	var flagBaseDenom = "token_base_denom"
	var flagUnitDenom = "token_unit_denom"
	var flagIbcChannelID = "token_ibc_channel_id"
	var flagIbcCounterpartyChannelID = "token_ibc_counterparty_channel_id"
	var flagIbcCounterpartyChainID = "token_ibc_counterparty_chain_id"
	var flagIbcCounterpartyDenom = "token_ibc_counterparty_denom"
	var flagDecimals = "token_decimals"
	var flagDisplayName = "token_display_name"
	var flagDisplaySymbol = "token_display_symbol"
	var flagExternalSymbol = "token_external_symbol"
	var flagTransferLimit = "token_transfer_limit"
	var flagNetwork = "token_network"
	var flagAddress = "token_address"
	var flagsPermission = []string{"token_permission_clp", "token_permission_ibc_export", "token_permission_ibc_import"}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate JSON for a token registration",
		Args:  cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			clientCtx, err = client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			flags := cmd.Flags()
			decimals, err := flags.GetInt64(flagDecimals)
			if err != nil {
				return err
			}
			initialDenom, err := flags.GetString(flagDenom)
			if err != nil {
				return err
			}
			baseDenom, err := flags.GetString(flagBaseDenom)
			if err != nil {
				return err
			}
			unitDenom, err := flags.GetString(flagUnitDenom)
			if err != nil {
				return err
			}
			ibcChannelID, err := flags.GetString(flagIbcChannelID)
			if err != nil {
				return err
			}
			ibcCounterpartyChannelID, err := flags.GetString(flagIbcCounterpartyChannelID)
			if err != nil {
				return err
			}
			ibcCounterpartyChainID, err := flags.GetString(flagIbcCounterpartyChainID)
			if err != nil {
				return err
			}
			ibcCounterpartyDenom, err := flags.GetString(flagIbcCounterpartyDenom)
			if err != nil {
				return err
			}
			displayName, err := flags.GetString(flagDisplayName)
			if err != nil {
				return err
			}
			displaySymbol, err := flags.GetString(flagDisplaySymbol)
			if err != nil {
				return err
			}
			externalSymbol, err := flags.GetString(flagExternalSymbol)
			if err != nil {
				return err
			}
			transferLimit, err := flags.GetString(flagTransferLimit)
			if err != nil {
				return err
			}
			network, err := flags.GetString(flagNetwork)
			if err != nil {
				return err
			}
			address, err := flags.GetString(flagAddress)
			if err != nil {
				return err
			}
			permissions := []types.Permission{}
			permissionCLP, err := flags.GetBool("token_permission_clp")
			if err != nil {
				return err
			}
			if permissionCLP {
				permissions = append(permissions, types.Permission_CLP)
			}
			permissionIBCExport, err := flags.GetBool("token_permission_ibc_export")
			if err != nil {
				return err
			}
			if permissionIBCExport {
				permissions = append(permissions, types.Permission_IBCEXPORT)
			}
			permissionIBCImport, err := flags.GetBool("token_permission_ibc_import")
			if err != nil {
				return err
			}
			if permissionIBCImport {
				permissions = append(permissions, types.Permission_IBCIMPORT)
			}
			var denom string
			var path string
			// base_denom is required.
			// override the IBC generation with --token_denom if specified explicitly.
			// otherwise fallback to base_denom
			if ibcChannelID != "" {
				path = "transfer/" + ibcChannelID
				// generate IBC hash from baseDenom and ibc channel id
				denomTrace := transfertypes.DenomTrace{
					Path:      path,
					BaseDenom: baseDenom,
				}
				denom = denomTrace.IBCDenom()
			} else if initialDenom == "" {
				// either initialDenom or channel id must be specified,
				// to prevent accidentally leaving off IBC details and
				return errors.New("--token_denom must be specified if no IBC channel is provided")
			}
			// --token_denom always takes precedence over IBC generation if specified
			if initialDenom != "" {
				denom = initialDenom
			}

			entry := types.RegistryEntry{
				Decimals:                 decimals,
				Denom:                    denom,
				BaseDenom:                baseDenom,
				Path:                     path,
				IbcChannelId:             ibcChannelID,
				IbcCounterpartyChannelId: ibcCounterpartyChannelID,
				IbcCounterpartyChainId:   ibcCounterpartyChainID,
				IbcCounterpartyDenom:     ibcCounterpartyDenom,
				UnitDenom:                unitDenom,
				DisplayName:              displayName,
				DisplaySymbol:            displaySymbol,
				Network:                  network,
				Address:                  address,
				ExternalSymbol:           externalSymbol,
				TransferLimit:            transferLimit,
				Permissions:              permissions,
			}
			return clientCtx.PrintProto(&types.Registry{Entries: []*types.RegistryEntry{&entry}})
		},
	}
	cmd.Flags().String(flagDenom, "",
		"The IBC hash / denom  stored on sifchain - to generate this hash for IBC token, leave blank and specify base_denom and ibc_channel_id")
	cmd.Flags().String(flagBaseDenom, "",
		"The base denom native to our chain, or native to an original chain (ie not the ibc hash)")
	cmd.Flags().String(flagIbcChannelID, "",
		"The channel id on our chain if this is an IBC token. Specify this to generate a new IBC hash to overwrite the denom field - used by clients when initiating send from this chain")
	cmd.Flags().String(flagIbcCounterpartyChannelID, "",
		"The counterparty channel if this is an IBC token - used by clients when initiating send from a counterparty chain")
	cmd.Flags().String(flagIbcCounterpartyChainID, "",
		"The chain id of ibc counter party chain")
	cmd.Flags().Int64(flagDecimals, -1,
		"The number of decimal points")
	cmd.Flags().String(flagUnitDenom, "",
		"The denom in registry that holds the funds for this denom, ie the most precise denom for a token")
	cmd.Flags().String(flagIbcCounterpartyDenom, "",
		"The denom in registry that funds in this account will get sent as over IBC")
	cmd.Flags().String(flagDisplayName, "",
		"Friendly name for use by UI etc")
	cmd.Flags().String(flagDisplaySymbol, "",
		"Friendly symbol for use by UI etc")
	cmd.Flags().String(flagExternalSymbol, "",
		"The original symbol as seen on external network")
	cmd.Flags().String(flagTransferLimit, "",
		"Used by UI")
	cmd.Flags().String(flagNetwork, "",
		"Original network of token i.e ethereum")
	cmd.Flags().String(flagAddress, "",
		"Contract address i.e in EVM cases")
	// Permission flags, default true.
	for _, flag := range flagsPermission {
		cmd.Flags().Bool(flag, true, fmt.Sprintf("Flag to specify permission for %s", types.GetPermissionFromString(flag)))
	}
	_ = cmd.MarkFlagRequired(flagBaseDenom)
	_ = cmd.MarkFlagRequired(flagDecimals)
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdAddEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [registry.json] [entry.json]",
		Short: "",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			registry, err := whitelistutils.ParseDenoms(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}
			reg, err := whitelistutils.ParseDenoms(clientCtx.Codec, args[1])
			if err != nil {
				return err
			}
			entryToAdd := reg.Entries[0]
			entries := registry.Entries
			entries = append(entries, entryToAdd)
			registry.Entries = entries
			return clientCtx.PrintBytes(clientCtx.Codec.MustMarshalJSON(&registry))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdAddAllEntries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-all [registry.json]",
		Short: "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			registry, err := whitelistutils.ParseDenoms(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}
			finalRegistry := types.Registry{Entries: []*types.RegistryEntry{}}
			for _, entry := range registry.Entries {
				entryForConversion := entry
				finalRegistry.Entries = append(finalRegistry.Entries, entryForConversion)
				if entry.Decimals > 10 {
					conversionDenom := ""
					if strings.HasPrefix(entry.Denom, "c") {
						conversionDenom = "x" + strings.TrimPrefix(entry.Denom, "c")
					} else if types.StringCompare(entry.Denom, "rowan") {
						conversionDenom = "xrowan"
					}
					entryForConversion.IbcCounterpartyDenom = conversionDenom
					entryForConversion.Permissions = []types.Permission{
						types.Permission_CLP,
						types.Permission_IBCEXPORT,
					}
					finalRegistry.Entries = append(finalRegistry.Entries, &types.RegistryEntry{
						Denom:       conversionDenom,
						BaseDenom:   conversionDenom,
						Decimals:    10,
						UnitDenom:   entry.Denom,
						Permissions: []types.Permission{types.Permission_IBCIMPORT},
					})
				} else {
					entryForConversion.Permissions = []types.Permission{
						types.Permission_CLP,
						types.Permission_IBCEXPORT,
						types.Permission_IBCIMPORT,
					}
				}
			}
			return clientCtx.PrintBytes(clientCtx.Codec.MustMarshalJSON(&finalRegistry))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
