package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
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
	var flagDecimals = "token_decimals"
	var flagDenom = "token_denom"
	var flagBaseDenom = "token_base_denom"
	var flagUnitDenom = "token_unit_denom"
	var flagIbcTransferPort = "token_ibc_transfer_port"
	var flagIbcChannelID = "token_ibc_channel_id"
	var flagIbcCounterpartyChannelID = "token_ibc_counterparty_channel_id"
	var flagIbcCounterpartyChainID = "token_ibc_counterparty_chain_id"
	var flagIbcCounterpartyDenom = "token_ibc_counterparty_denom"
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
			decimals, err := flags.GetInt(flagDecimals)
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
			ibcChannelId, err := flags.GetString(flagIbcChannelID)
			if err != nil {
				return err
			}
			ibcTransferPort, err := flags.GetString(flagIbcTransferPort)
			if err != nil {
				return err
			}
			ibcCounterpartyChannelId, err := flags.GetString(flagIbcCounterpartyChannelID)
			if err != nil {
				return err
			}
			ibcCounterpartyChainId, err := flags.GetString(flagIbcCounterpartyChainID)
			if err != nil {
				return err
			}
			ibcCounterpartyDenom, err := flags.GetString(flagIbcCounterpartyDenom)
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
				path = ibcTransferPort + "/" + ibcChannelId
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
				Decimals:                 int32(decimals),
				Denom:                    denom,
				BaseDenom:                baseDenom,
				IbcTransferPort:          ibcTransferPort,
				IbcChannelId:             ibcChannelId,
				IbcCounterpartyChannelId: ibcCounterpartyChannelId,
				IbcCounterpartyChainId:   ibcCounterpartyChainId,
				IbcCounterpartyDenom:     ibcCounterpartyDenom,
				UnitDenom:                unitDenom,
				Permissions:              permissions,
			}
			return clientCtx.PrintProto(&types.Registry{Entries: []*types.RegistryEntry{&entry}})
		},
	}
	cmd.Flags().Bool(flagIbcTransferPort, true,
		"The ibc transfer port name (default: 'transfer')")
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
	cmd.Flags().Int(flagDecimals, -1,
		"The number of decimal points")
	cmd.Flags().String(flagUnitDenom, "",
		"The denom in registry that holds the funds for this denom, ie the most precise denom for a token")
	cmd.Flags().String(flagIbcCounterpartyDenom, "",
		"The denom in registry that funds in this account will get sent as over IBC")
	// Permission flags, default true.
	for _, flag := range flagsPermission {
		cmd.Flags().Bool(flag, true, fmt.Sprintf("Flag to specify permission for %s", types.GetPermissionFromString(flag)))
	}
	_ = cmd.MarkFlagRequired(flagBaseDenom)
	_ = cmd.MarkFlagRequired(flagDecimals)
	return cmd
}
