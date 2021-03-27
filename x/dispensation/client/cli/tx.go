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
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/multisig"
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
		GetCmdAirdrop(cdc),
	)...)

	return dispensationTxCmd
}

func GetCmdAirdrop(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "airdrop [address] [input] [output]",
		Short: "Create new airdrop",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			kb, err := keys.NewKeyring(sdk.KeyringServiceName(),
				viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), inBuf)
			if err != nil {
				return err
			}
			fmt.Println("trying to get keys", args[0])
			fmt.Println("KB ", kb)
			fmt.Println("Flags ", viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome))

			multisigInfo, err := kb.Get(args[0])
			if err != nil {
				return err
			}

			multisigPub := multisigInfo.GetPubKey().(multisig.PubKeyMultisigThreshold)
			pubkeys := multisigPub.PubKeys
			for _, p := range pubkeys {
				fmt.Println(p.Address().String())
			}

			ko, err := keys.Bech32KeyOutput(multisigInfo)
			if err != nil {
				return err
			}
			if multisigInfo.GetType() != keys.TypeMulti {
				return fmt.Errorf("%q must be of type %s: %s", args[0], keys.TypeMulti, multisigInfo.GetType())
			}

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, ko.Address).WithCodec(cdc)

			inputList, err := dispensationUtils.ParseInput(args[1])
			if err != nil {
				return err
			}
			outputlist, err := dispensationUtils.ParseOutput(args[2])
			if err != nil {
				return err
			}
			for _, i := range inputList {
				fmt.Println("Inputs : ", i.Address.String())
			}

			msg := types.NewMsgAirdrop(cliCtx.GetFromAddress(), inputList, outputlist)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
