package cli

import (
	"bufio"
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Use:   "airdrop [address]", //--input [path] --output [path]
		Short: "Create new airdrop",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			kb, err := keys.NewKeyring(sdk.KeyringServiceName(),
				viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), inBuf)
			if err != nil {
				return err
			}
			multisigInfo, err := kb.Get(args[0])
			if err != nil {
				return err
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

			in, _ := sdk.AccAddressFromBech32("sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd")
			out, _ := sdk.AccAddressFromBech32("sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5")
			coin := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(10))}
			inputList := []bank.Input{bank.NewInput(in, coin)}
			outputlist := []bank.Output{bank.NewOutput(out, coin)}
			msg := types.NewMsgAirdrop(cliCtx.GetFromAddress(), inputList, outputlist)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
