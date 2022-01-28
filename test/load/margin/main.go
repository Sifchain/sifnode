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
			initClientCtx = client.ReadHomeFlag(initClientCtx, cmd)
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

	// create x tx's of y positions
	x := 100
	y := 100
	z := 1

	count := make(chan int, 1)

	go func() {
		txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())

		accountNumber, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, clientCtx.GetFromAddress())
		if err != nil {
			panic(err)
		}

		txf = txf.WithAccountNumber(accountNumber).WithSequence(seq)

		for c := 0; c < z; c++ {
			//go func() {
			for a := 0; a < x; a++ {

				collateralAsset := "rowan"
				collateralAmount := uint64(100)
				borrowAsset := "stake"

				var msgs []sdk.Msg
				for b := 0; b < y; b++ {
					msgs = append(msgs, &types.MsgOpenLong{
						Signer:           clientCtx.GetFromAddress().String(),
						CollateralAsset:  collateralAsset,
						CollateralAmount: sdk.NewUint(collateralAmount),
						BorrowAsset:      borrowAsset,
					})
				}

				err = tx.BroadcastTx(clientCtx, txf.WithSequence(seq).WithSimulateAndExecute(true), msgs...)
				if err != nil {
					log.Printf("ERR %s", err)
				}

				seq++
				count <- y
			}
			//}()
		}
	}()

	var total int
	for {
		select {
		case c := <-count:
			total += c
			log.Printf("%d positions opened", total)

			if c >= x*y*z {
				return nil
			}
		}
	}
}
