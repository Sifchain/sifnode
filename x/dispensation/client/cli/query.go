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
		GetCmdDistributionRecordForDistName(queryRoute),
		GetCmdClaimsByType(queryRoute),
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
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAllDistributions)
			res, height, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}
			var dr types.Distributions
			types.ModuleCdc.MustUnmarshalJSON(res, &dr)
			out := types.NewQueryAllDistributionsResponse(dr, height)
			return clientCtx.PrintProto(&out)
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
			address := args[0]
			recipientAddress, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				return err
			}
			params := types.QueryRecordsByRecipientAddrRequest{
				Address: recipientAddress.String()}
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
			types.ModuleCdc.MustUnmarshalJSON(res, &drs)
			out := types.NewQueryRecordsByRecipientAddrResponse(drs, height)
			return clientCtx.PrintProto(&out)
		},
	}
}

//GetCmdDistributionRecordForDistName returns all records for a given distribution name
func GetCmdDistributionRecordForDistName(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "records-by-name [distribution name] [status]",
		Short: "get a list of all distribution records Status : [Completed/Pending/All]",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			name := args[0]
			status, ok := types.GetDistributionStatus(args[1])
			if !ok {
				return fmt.Errorf("invalid Status %s: Status supported [Completed/Pending/Failed]", args[0])
			}
			params := types.QueryRecordsByDistributionNameRequest{
				DistributionName: name,
				Status:           status}
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
			out := types.NewQueryRecordsByDistributionNameResponse(drs, height)
			return clientCtx.PrintProto(&out)
		},
	}
}

func GetCmdClaimsByType(queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "claims-by-type [ClaimType]",
		Short: "get a list of all claims for mentioned type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			claimType, ok := types.GetClaimType(args[0])
			if !ok {
				return fmt.Errorf("invalid Claim Type %s: Types supported [LiquidityMining/ValidatorSubsidy]", args[0])
			}
			params := types.QueryClaimsByTypeRequest{
				UserClaimType: claimType,
			}
			bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryClaimsByType)
			res, height, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var claims types.UserClaims
			types.ModuleCdc.MustUnmarshalJSON(res, &claims)
			out := types.QueryClaimsResponse{
				Claims: claims.UserClaims,
				Height: height,
			}
			return clientCtx.PrintProto(&out)
		},
	}
}
