package cli

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	dispensationTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	dispensationTxCmd.AddCommand(
		GetCmdCreate(),
		GetCmdClaim(),
		GetCmdRun(),
	)

	return dispensationTxCmd
}

// GetCmdCreate adds a new command to the main dispensationTxCmd to create a new airdrop
// Airdrop is a type of distribution on the network .
func GetCmdCreate() *cobra.Command {
	// Note ,the command only creates a airdrop for now .
	cmd := &cobra.Command{
		Use:   "create [DistributionType] [Output JSON File Path] [AuthorizedRunner]",
		Short: "Create new distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(3)(cmd, args)
			if err != nil {
				return err
			}
			distributionType, ok := types.GetDistributionTypeFromShortString(args[0])
			if !ok {
				return fmt.Errorf("invalid distribution Type %s: Types supported [Airdrop/LiquidityMining/ValidatorSubsidy]", args[2])
			}
			outputList, err := dispensationUtils.ParseOutput(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgCreateDistribution(clientCtx.GetFromAddress(), distributionType, outputList, args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [ClaimType]",
		Short: "Create new Claim",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			claimType, ok := types.GetClaimType(args[0])
			if !ok {
				return fmt.Errorf("invalid Claim Type %s: Types supported [LiquidityMining/ValidatorSubsidy]", args[0])
			}
			msg := types.NewMsgCreateUserClaim(clientCtx.GetFromAddress(), claimType)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd

}

func GetCmdRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [DistributionName] [DistributionType] [DistributionCount]",
		Short: "run limited records dispensation by specifying the name / should only be called by the authorized runner",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(3)(cmd, args)
			if err != nil {
				return err
			}
			distributionType, ok := types.GetDistributionTypeFromShortString(args[1])
			if !ok {
				return fmt.Errorf("invalid distribution Type %s: Types supported [Airdrop/LiquidityMining/ValidatorSubsidy]", args[1])
			}
			dispensationCount, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid dispensation count :%d", dispensationCount)
			}
			msg := types.NewMsgRunDistribution(clientCtx.GetFromAddress().String(), args[0], distributionType, dispensationCount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
