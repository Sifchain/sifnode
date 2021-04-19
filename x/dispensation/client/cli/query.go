package cli

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	// Group dispensation queries under a subcommand
	dispensationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	dispensationQueryCmd.AddCommand(
		GetCmdDistributions(queryRoute, cdc),
		GetCmdDistributionRecordForRecipient(queryRoute, cdc),
		GetCmdDistributionRecordForDistNameAll(queryRoute, cdc),
		GetCmdDistributionRecordForDistNamePending(queryRoute, cdc),
		GetCmdDistributionRecordForDistNameCompleted(queryRoute, cdc),
	)
	return dispensationQueryCmd
}

//GetCmdDistributions returns a list of all distributions ever created
func GetCmdDistributions(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "distributions-all",
		Short: "get a list of all distributions ",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAllDistributions)
			res, height, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			var dr types.Distributions
			cdc.MustUnmarshalJSON(res, &dr)
			out := types.NewDistributionsResponse(dr, height)
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}

// GetCmdDistributionRecordForRecipient returns the completed and pending records for the recipient address
func GetCmdDistributionRecordForRecipient(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-addr [recipient address]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			address := args[0]
			recipientAddress, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				return err
			}
			params := types.NewQueryRecordsByRecipientAddr(recipientAddress)
			bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByRecipient)
			res, height, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}

//GetCmdDistributionRecordForDistNameAll returns all records for a given distribution name
func GetCmdDistributionRecordForDistNameAll(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-name-all [distribution name]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			name := args[0]
			params := types.NewQueryRecordsByDistributionName(name, 3)
			bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByDistrName)
			res, height, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}

//GetCmdDistributionRecordForDistNamePending returns all pending records for a given distribution name
func GetCmdDistributionRecordForDistNamePending(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-name-pending [distribution name]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			name := args[0]
			params := types.NewQueryRecordsByDistributionName(name, types.Pending)
			bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByDistrName)
			res, height, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}

//GetCmdDistributionRecordForDistNamePending returns all completed records for a given distribution name
func GetCmdDistributionRecordForDistNameCompleted(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-name-completed [distribution name]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			name := args[0]
			params := types.NewQueryRecordsByDistributionName(name, types.Completed)
			bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByDistrName)
			res, height, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}
