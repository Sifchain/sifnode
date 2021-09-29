package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
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
		GetCmdGenerateLowPrecisionEntries(),
		GetCmdGenerateHighPrecisionEntries(),
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

			return clientCtx.PrintBytes(clientCtx.JSONMarshaler.MustMarshalJSON(res.Registry))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdGenerateEntry() *cobra.Command {
	var flagWhitelist = "token_whitelist"
	var flagDenom = "token_denom"
	var flagBaseDenom = "token_base_denom"
	var flagIbcChannelId = "token_ibc_channel_id"
	var flagIbcCounterpartyChannelId = "token_ibc_counterparty_channel_id"
	var flagIbcCounterpartyChainId = "token_ibc_counterparty_chain_id"
	var flagIbcCounterpartyDenom = "token_ibc_counterparty_denom"
	var flagUnitDenom = "token_unit_denom"
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
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			flags := cmd.Flags()

			whitelist, err := flags.GetBool(flagWhitelist)
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

			decimals, err := flags.GetInt(flagDecimals)
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

			ibcChannelId, err := flags.GetString(flagIbcChannelId)
			if err != nil {
				return err
			}

			ibcCounterpartyChannelId, err := flags.GetString(flagIbcCounterpartyChannelId)
			if err != nil {
				return err
			}

			ibcCounterpartyChainId, err := flags.GetString(flagIbcCounterpartyChainId)
			if err != nil {
				return err
			}

			ibcCounterpartyDenom, err := flags.GetString(flagIbcCounterpartyDenom)
			if err != nil {
				return err
			}

			unitDenom, err := flags.GetString(flagUnitDenom)
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

			var path string
			var denom string
			// base_denom is required.
			// generate denom if path is also provided.
			// override the IBC generation with --denom if specified explicitly.
			// otherwise fallback to base_denom

			if ibcChannelId != "" {
				// normalise path slashes before generating hash (do this in MsgRegister.ValidateBasic as well)
				path = "transfer/" + ibcChannelId

				// generate IBC hash from baseDenom and ibc channel id
				denomTrace := transfertypes.DenomTrace{
					Path:      path,
					BaseDenom: baseDenom,
				}

				denom = denomTrace.IBCDenom()
			}

			if initialDenom != "" {
				denom = initialDenom
			} else if denom == "" {
				denom = baseDenom
			}

			entry := types.RegistryEntry{
				IsWhitelisted:            whitelist,
				Decimals:                 int64(decimals),
				Denom:                    denom,
				BaseDenom:                baseDenom,
				Path:                     path,
				IbcChannelId:             ibcChannelId,
				IbcCounterpartyChannelId: ibcCounterpartyChannelId,
				IbcCounterpartyChainId:   ibcCounterpartyChainId,
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

	cmd.Flags().Bool(flagWhitelist, true,
		"Whether this token should be whitelisted i.e disable all permissions.")
	cmd.Flags().String(flagDenom, "",
		"The IBC hash / denom  stored on sifchain - to generate this hash for IBC token, leave blank and specify base_denom and ibc_channel_id.")
	cmd.Flags().String(flagBaseDenom, "",
		"The base denom native to our chain, or native to an original chain (ie not the ibc hash).")
	cmd.Flags().String(flagIbcChannelId, "",
		"The channel id on our chain if this is an IBC token. Specify this to generate a new IBC hash to overwrite the denom field - used by clients when initiating send from this chain")
	cmd.Flags().String(flagIbcCounterpartyChannelId, "",
		"The counterparty channel if this is an IBC token - used by clients when initiating send from a counterparty chain")
	cmd.Flags().String(flagIbcCounterpartyChainId, "",
		"The chain id of ibc counter party chain")
	cmd.Flags().Int(flagDecimals, -1,
		"The number of decimal points")
	cmd.Flags().String(flagUnitDenom, "",
		"The denom in registry that holds the funds for this denom, ie the most precise denom for a token.")
	cmd.Flags().String(flagIbcCounterpartyDenom, "",
		"The denom in registry that funds in this account will get sent as over IBC.")
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

	return cmd
}

func GetCmdGenerateLowPrecisionEntries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-low-precision-entries [registry.json]",
		Short: "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			registry, err := whitelistutils.ParseDenoms(clientCtx.JSONMarshaler, args[0])
			if err != nil {
				return err
			}

			lowPrecisionTokenRegistry := types.Registry{Entries: []*types.RegistryEntry{}}

			for _, entry := range registry.Entries {
				if entry.Decimals > 10 && strings.HasPrefix(entry.Denom, "c") {
					conversionDenom := "x" + strings.TrimPrefix(entry.Denom, "c")

					lowPrecisionTokenRegistry.Entries = append(lowPrecisionTokenRegistry.Entries, &types.RegistryEntry{
						IsWhitelisted: true,
						Denom:         conversionDenom,
						BaseDenom:     conversionDenom,
						Decimals:      10,
						UnitDenom:     entry.Denom,
						Permissions:   []types.Permission{types.Permission_IBCIMPORT},
					})
				}
			}

			return clientCtx.PrintBytes(clientCtx.JSONMarshaler.MustMarshalJSON(&lowPrecisionTokenRegistry))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdGenerateHighPrecisionEntries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-high-precision-entries [registry.json]",
		Short: "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			registry, err := whitelistutils.ParseDenoms(clientCtx.JSONMarshaler, args[0])
			if err != nil {
				return err
			}

			highPrecisionTokenRegistry := types.Registry{Entries: []*types.RegistryEntry{}}

			for _, entry := range registry.Entries {
				if entry.Decimals > 10 && strings.HasPrefix(entry.Denom, "c") {
					entryForConversion := entry

					conversionDenom := "x" + strings.TrimPrefix(entry.Denom, "c")

					highPrecisionTokenRegistry.Entries = append(highPrecisionTokenRegistry.Entries, entryForConversion)

					entryForConversion.IbcCounterpartyDenom = conversionDenom

					entryForConversion.Permissions = []types.Permission{
						types.Permission_CLP,
						types.Permission_IBCEXPORT,
						types.Permission_IBCIMPORT,
					}
				}
			}

			return clientCtx.PrintBytes(clientCtx.JSONMarshaler.MustMarshalJSON(&highPrecisionTokenRegistry))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
