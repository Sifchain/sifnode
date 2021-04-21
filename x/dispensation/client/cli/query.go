package cli

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group dispensation queries under a subcommand
	dispensationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	dispensationQueryCmd.AddCommand(
		GetCmdDistributions(queryRoute),
		GetCmdDistributionRecordForRecipient(queryRoute),
		GetCmdDistributionRecordForDistNameAll(queryRoute),
		GetCmdDistributionRecordForDistNamePending(queryRoute),
		GetCmdDistributionRecordForDistNameCompleted(queryRoute),
	)
	return dispensationQueryCmd
}

//GetCmdDistributions returns a list of all distributions ever created
func GetCmdDistributions(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "distributions-all",
		Short: "get a list of all distributions ",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			return clientCtx.PrintObjectLegacy("temp")
			if err != nil {
				return err
			}
			return clientCtx.PrintObjectLegacy("temp")
			//route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAllDistributions)
			//res, height, err := clientCtx.QueryWithData(route, nil)
			//if err != nil {
			//	return err
			//}
			//var dr types.Distributions
			//types.ModuleCdc.MustUnmarshalJSON(res, &dr)
			//out := types.NewDistributionsResponse(dr, height)
			//return clientCtx.PrintObjectLegacy(out)
		},
	}
}

// GetCmdDistributionRecordForRecipient returns the completed and pending records for the recipient address
func GetCmdDistributionRecordForRecipient(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-addr [recipient address]",
		Short: "get a list of all distribution records ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			return clientCtx.PrintObjectLegacy("temp")
			//address := args[0]
			//recipientAddress, err := sdk.AccAddressFromBech32(address)
			//if err != nil {
			//	return err
			//}
			//params := types.NewQueryRecordsByRecipientAddr(recipientAddress)
			//bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			//if err != nil {
			//	return err
			//}
			//route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRecordsByRecipient)
			//res, height, err := clientCtx.QueryWithData(route, bz)
			//if err != nil {
			//	return err
			//}
			//var drs types.DistributionRecords
			//types.ModuleCdc.MustUnmarshalJSON(res, &drs)
			//out := types.NewDistributionRecordsResponse(drs, height)
			//return clientCtx.PrintObjectLegacy(out)
		},
	}
}

//GetCmdDistributionRecordForDistNameAll returns all records for a given distribution name
func GetCmdDistributionRecordForDistNameAll(queryRoute string) *cobra.Command {
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
			params := types.NewQueryRecordsByDistributionName(name, sdk.NewUint(2))
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
			types.ModuleCdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, sdk.NewInt(height))
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}

//GetCmdDistributionRecordForDistNamePending returns all pending records for a given distribution name
func GetCmdDistributionRecordForDistNamePending(queryRoute string) *cobra.Command {
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
			params := types.NewQueryRecordsByDistributionName(name, sdk.NewUint(1))
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
			types.ModuleCdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, sdk.NewInt(height))
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}

//GetCmdDistributionRecordForDistNamePending returns all completed records for a given distribution name
func GetCmdDistributionRecordForDistNameCompleted(queryRoute string) *cobra.Command {
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
			params := types.NewQueryRecordsByDistributionName(name, sdk.NewUint(2))
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
			types.ModuleCdc.MustUnmarshalJSON(res, &drs)
			out := types.NewDistributionRecordsResponse(drs, sdk.NewInt(height))
			return clientCtx.PrintObjectLegacy(out)
		},
	}
}
