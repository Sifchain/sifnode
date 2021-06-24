package cli

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gogo/protobuf/codec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
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
	)

	return dispensationTxCmd
}

// GetCmdCreate adds a new command to the main dispensationTxCmd to create a new airdrop
// Airdrop is a type of distribution on the network .
func GetCmdCreate() *cobra.Command {
	// Note ,the command only creates a airdrop for now .
	cmd := &cobra.Command{
		Use:   "distribute [MultiSigKeyName] [DistributionName] [DistributionType] [Input JSON File Path] [Output JSON File Path]",
		Short: "Create new distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(viper.GetString(cli.HomeFlag))

			inBuf := bufio.NewReader(cmd.InOrStdin())
			keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
			if err != nil {
				return err
			}

			// attempt to lookup address from Keybase if no address was provided
			kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
			if err != nil {
				return err
			}

			multisigInfo, err := kb.Key(args[0])
			if err != nil {
				return fmt.Errorf("failed to get address from Keybase: %w", err)
			}

			if multisigInfo.GetType() != keyring.TypeMulti {
				return fmt.Errorf("%q must be of type %s: %s", args[0], keyring.TypeMulti, multisigInfo.GetType())
			}
			name := args[1]
			distributionType, ok := types.IsValidDistribution(args[2])
			if !ok {
				return fmt.Errorf("invalid distribution Type %s: Types supported [Airdrop/LiquidityMining/ValidatorSubsidy]", args[0])
			}

			inputList, err := dispensationUtils.ParseInput(args[3])
			if err != nil {
				return err
			}

			multisigPub := multisigInfo.GetPubKey().(*multisig.LegacyAminoPubKey)
			err = dispensationUtils.VerifyInputList(inputList, multisigPub.PubKeys)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("Invalid Address for authorised distributor : %s", args[2]))
			}

			outputlist, err := dispensationUtils.ParseOutput(args[4])
			if err != nil {
				return err
			}
			msg := types.NewMsgCreateDistribution(clientCtx.GetFromAddress(), name, distributionType, inputList, outputlist)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			claimType, ok := types.IsValidClaim(args[0])
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
	return cmd
}

func GetCmdRun(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [DistributionName] [DistributionType]",
		Short: "run a dispensation by specifying the name / should only be called by the authorized runner",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			distributionType, ok := types.IsValidDistributionType(args[1])
			if !ok {
				return fmt.Errorf("invalid distribution Type %s: Types supported [Airdrop/LiquidityMining/ValidatorSubsidy]", args[0])
			}
			msg := types.NewMsgRunDistribution(cliCtx.GetFromAddress(), args[0], distributionType)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
