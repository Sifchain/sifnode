package cli

import (
	"context"
	"fmt"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	"github.com/spf13/cobra"
	"strings"
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

			return clientCtx.PrintProto(res.Registry)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdGenerateEntry() *cobra.Command {
	var flagWhitelist = "token_whitelist"
	var flagDenom = "token_denom"
	var flagBaseDenom = "token_base_denom"
	var flagPath = "token_path"
	var flagSrcChannel = "token_src_channel"
	var flagDestChannel = "token_dest_channel"
	var flagDecimals = "token_decimals"
	var flagDisplayName = "token_display_name"
	var flagDisplaySymbol = "token_display_symbol"
	var flagExternalSymbol = "token_external_symbol"
	var flagTransferLimit = "token_transfer_limit"
	var flagNetwork = "token_network"
	var flagAddress = "token_address"

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

			path, err := flags.GetString(flagPath)
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

			srcChannel, err := flags.GetString(flagSrcChannel)
			if err != nil {
				return err
			}

			destChannel, err := flags.GetString(flagDestChannel)
			if err != nil {
				return err
			}

			// normalise path slashes before generating hash (do this in MsgRegister.ValidateBasic as well)
			path = strings.Trim(path, "/")

			var denom string
			// base_denom is required.
			// generate denom if path is also provided.
			// override the IBC generation with --denom if specified explicitly.
			// otherwise fallback to base_denom

			if path != "" {
				// generate IBC hash from baseDenom and path
				denomTrace := transfertypes.DenomTrace{
					Path:      path,
					BaseDenom: baseDenom,
				}

				denom = denomTrace.IBCDenom()
			}

			if initialDenom != "" {
				denom = initialDenom
			} else {
				denom = baseDenom
			}

			entry := types.RegistryEntry{
				IsWhitelisted:  whitelist,
				Decimals:       int64(decimals),
				Denom:          denom,
				BaseDenom:      baseDenom,
				Path:           path,
				SrcChannel:     srcChannel,
				DestChannel:    destChannel,
				DisplayName:    displayName,
				DisplaySymbol:  displaySymbol,
				Network:        network,
				Address:        address,
				ExternalSymbol: externalSymbol,
				TransferLimit:  transferLimit,
			}

			return clientCtx.PrintProto(&types.Registry{Entries: []*types.RegistryEntry{&entry}})
		},
	}

	cmd.Flags().Bool(flagWhitelist, true,
		"Whether this token should be on whitelist")
	cmd.Flags().String(flagDenom, "",
		"The IBC hash / denom  stored on sifchain - to generate this hash for IBC token, leave blank and specify base_denom and path.")
	cmd.Flags().String(flagBaseDenom, "",
		"The denom native to our chain, or native to an original chain (i.e the non-path part, underlying an IBC hash token).")
	cmd.Flags().String(flagPath, "",
		"IBC path using the *SRC* port + channel ID on our chain and other IBC hops receiving this token (leave blank for non-IBC) i.e transfer/channel-0")
	cmd.Flags().String(flagSrcChannel, "",
		"The src channel if this is an IBC token - used by UI when initiating send from this chain")
	cmd.Flags().String(flagDestChannel, "",
		"The dest channel if this is an IBC token - used by UI when initiating import from originating chain")
	cmd.Flags().Int(flagDecimals, -1,
		"The number of decimal points")
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

	cmd.MarkFlagRequired(flagBaseDenom)
	cmd.MarkFlagRequired(flagDecimals)

	return cmd
}
