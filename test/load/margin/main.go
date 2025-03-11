package main

import (
	"log"
	"os"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

func main() {
	encodingConfig := app.MakeTestEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")
	app.SetConfig(false)

	rootCmd := &cobra.Command{
		Use: "marginloadtest",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}
			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
		RunE: run,
	}
	flags.AddTxFlagsToCmd(rootCmd)
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	err := svrcmd.Execute(rootCmd, app.DefaultNodeHome)
	if err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	// get pools
	// pools := []string{"stake"}

	positions := 1000

	txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	key, err := txf.Keybase().Key(clientCtx.GetFromName())
	if err != nil {
		panic(err)
	}

	accountNumber, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, key.GetAddress())
	if err != nil {
		panic(err)
	}

	txf.WithAccountNumber(accountNumber)

	for a := 0; a < positions; a++ {
		txf = txf.WithSequence(seq + uint64(a)) // nolint:gosec
		err := broadcastTrade(clientCtx, txf, key)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func broadcastTrade(clientCtx client.Context, txf tx.Factory, key keyring.Info) error {
	collateralAsset := "rowan"
	collateralAmount := uint64(100)
	borrowAsset := "ceth"

	msg := types.MsgOpen{
		Signer:           key.GetAddress().String(),
		CollateralAsset:  collateralAsset,
		CollateralAmount: sdk.NewUint(collateralAmount),
		BorrowAsset:      borrowAsset,
		Position:         types.Position_LONG,
	}
	txb, err := tx.BuildUnsignedTx(txf, &msg)
	if err != nil {
		panic(err)
	}
	err = tx.Sign(txf, key.GetName(), txb, true)
	if err != nil {
		panic(err)
	}
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		panic(err)
	}
	res, err := clientCtx.WithSimulation(true). /*.WithBroadcastMode("block")*/ BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	log.Print(res)

	return err
}
