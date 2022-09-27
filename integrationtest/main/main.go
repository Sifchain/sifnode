package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Sifchain/sifnode/app"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
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
		Use: "integrationtest",
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

	txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	key, err := txf.Keybase().Key(clientCtx.GetFromName())

	accountNumber, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, key.GetAddress())
	if err != nil {
		panic(err)
	}

	txf = txf.WithAccountNumber(accountNumber).WithSequence(seq)
	err = TestAddLiquidity(clientCtx, txf, key)
	if err != nil {
		panic(err)
	}

	return nil
}

func TestAddLiquidity(clientCtx client.Context, txf tx.Factory, key keyring.Info) error {
	clpQueryClient := clptypes.NewQueryClient(clientCtx)

	poolBefore, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: "ceth"})
	if err != nil {
		return err
	}

	nativeAdd := poolBefore.Pool.NativeAssetBalance.Quo(sdk.NewUint(10))
	externalAdd := poolBefore.Pool.ExternalAssetBalance.Quo(sdk.NewUint(10))

	msg := clptypes.MsgAddLiquidity{
		Signer:              key.GetAddress().String(),
		ExternalAsset:       &clptypes.Asset{Symbol: "ceth"},
		NativeAssetAmount:   nativeAdd,
		ExternalAssetAmount: externalAdd,
	}

	if err := buildAndBroadcast(clientCtx, txf, key, &msg); err != nil {
		return err
	}

	poolAfter, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: "ceth"})
	if err != nil {
		return err
	}

	if !poolBefore.Pool.NativeAssetBalance.Equal(poolAfter.Pool.NativeAssetBalance.Sub(nativeAdd)) {
		return errors.New(fmt.Sprintf("native balance mismatch afer add (before: %s after: %s)",
			poolBefore.Pool.NativeAssetBalance.String(),
			poolAfter.Pool.NativeAssetBalance.String()))
	}

	if !poolAfter.Pool.ExternalAssetBalance.Sub(externalAdd).Equal(poolBefore.Pool.ExternalAssetBalance) {
		return errors.New(fmt.Sprintf("external balance mismatch afer add (added: %s diff: %s)",
			externalAdd,
			poolAfter.Pool.ExternalAssetBalance.Sub(poolBefore.Pool.ExternalAssetBalance).String()))
	}

	return nil
}

func broadcastOpenPosition(clientCtx client.Context, txf tx.Factory, key keyring.Info) error {
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

func buildAndBroadcast(clientCtx client.Context, txf tx.Factory, key keyring.Info, msg sdk.Msg) error {
	txb, err := tx.BuildUnsignedTx(txf, msg)
	if err != nil {
		return err
	}

	err = tx.Sign(txf, key.GetName(), txb, true)
	if err != nil {
		return err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return err
	}

	res, err := clientCtx.
		WithSimulation(true).
		WithBroadcastMode("block").
		BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	log.Print(res)

	return err
}
