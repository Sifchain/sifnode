package cli

import (
	"bufio"
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	dispensationTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	dispensationTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreate(cdc),
		GetCmdClaim(cdc),
	)...)

	return dispensationTxCmd
}

// GetCmdCreate adds a new command to the main dispensationTxCmd to create a new airdrop
// Airdrop is a type of distribution on the network .
func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [DistributionType] [Output JSON File Path]",
		Short: "Create new distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			distributionType, ok := types.IsValidDistributionType(args[0])
			if !ok {
				return fmt.Errorf("invalid distribution Type %s: Types supported [Airdrop/LiquidityMining/ValidatorSubsidy]", args[1])
			}
			outputList, err := dispensationUtils.ParseOutput(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgDistribution(cliCtx.GetFromAddress(), distributionType, outputList)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdClaim(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [ClaimType]",
		Short: "Create new Claim",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			claimType, ok := types.IsValidClaim(args[0])
			if !ok {
				return fmt.Errorf("invalid Claim Type %s: Types supported [LiquidityMining/ValidatorSubsidy]", args[0])
			}
			msg := types.NewMsgCreateClaim(cliCtx.GetFromAddress(), claimType)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
