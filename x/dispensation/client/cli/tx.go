package cli

import (
	"bufio"
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
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
	)

	return dispensationTxCmd
}

// GetCmdCreate adds a new command to the main dispensationTxCmd to create a new airdrop
// Airdrop is a type of distribution on the network .
func GetCmdCreate() *cobra.Command {
	// Note ,the command only creates a airdrop for now .
	cmd := &cobra.Command{
		Use:   "create [MultiSigKeyName] [DistributionName] [Input JSON File Path] [Output JSON File Path]",
		Short: "Create new distribution",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			//depCdc := clientCtx.JSONMarshaler
			//cdc := depCdc.(codec.Marshaler)
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

			inputList, err := dispensationUtils.ParseInput(args[2])
			if err != nil {
				return err
			}

			multisigPub := multisigInfo.GetPubKey().(*multisig.LegacyAminoPubKey)
			err = dispensationUtils.VerifyInputList(inputList, multisigPub.PubKeys)
			if err != nil {
				return err
			}

			outputlist, err := dispensationUtils.ParseOutput(args[3])
			if err != nil {
				return err
			}
			name := args[1]
			msg := types.NewMsgDistribution(clientCtx.GetFromAddress(), name, types.Airdrop, inputList, outputlist)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	return cmd
}
