package cli

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group dispensation queries under a subcommand
	dispensationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	dispensationQueryCmd.AddCommand(flags.GetCommands(
		GetCmdDistributions(queryRoute, cdc),
		GetCmdDistributionRecordForRecipient(queryRoute, cdc),
		GetCmdDistributionRecordForDistNameAll(queryRoute, cdc),
		GetCmdDistributionRecordForDistNamePending(queryRoute, cdc),
		GetCmdDistributionRecordForDistNameCompleted(queryRoute, cdc),
	)...)
	return dispensationQueryCmd
}

//GetCmdDistributions returns a list of all distributions ever created
func GetCmdDistributions(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "distributions-all",
		Short: "get a list of all distributions ",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAllDistributions)
			res, height, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			var dr types.Distributions
			cdc.MustUnmarshalJSON(res, &dr)
			out := types.NewDistributionsResponse(dr, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdDistributionRecordForRecipient returns the completed and pending records for the recipient address
func GetCmdDistributionRecordForRecipient(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-addr [recipient address]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]
			recipientAddress, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				return err
			}
			params := types.NewQueryRecordsByRecipientAddr(recipientAddress)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByRecipient)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

//GetCmdDistributionRecordForDistNameAll returns all records for a given distribution name
func GetCmdDistributionRecordForDistNameAll(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-name-all [distribution name]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]
			params := types.NewQueryRecordsByDistributionName(name, 3)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByDistrName)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

//GetCmdDistributionRecordForDistNamePending returns all pending records for a given distribution name
func GetCmdDistributionRecordForDistNamePending(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-name-pending [distribution name]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]
			params := types.NewQueryRecordsByDistributionName(name, types.Pending)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByDistrName)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return cliCtx.PrintOutput(out)
		},
	}
}

//GetCmdDistributionRecordForDistNamePending returns all completed records for a given distribution name
func GetCmdDistributionRecordForDistNameCompleted(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-name-completed [distribution name]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]
			params := types.NewQueryRecordsByDistributionName(name, types.Completed)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByDistrName)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var drs types.DistributionRecords
			cdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, height)
			return cliCtx.PrintOutput(out)
		},
	}
}
