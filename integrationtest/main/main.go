package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Sifchain/sifnode/app"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
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

	lpBefore, err := clpQueryClient.GetLiquidityProvider(context.Background(), &clptypes.LiquidityProviderReq{
		Symbol:    "ceth",
		LpAddress: key.GetAddress().String(),
	})
	if err != nil {
		// if lp doesn't exist
		lpBefore = &clptypes.LiquidityProviderRes{
			LiquidityProvider: &clptypes.LiquidityProvider{
				LiquidityProviderUnits: sdk.ZeroUint(),
			},
		}
	}

	nativeAdd := poolBefore.Pool.NativeAssetBalance.Quo(sdk.NewUint(1000))
	externalAdd := poolBefore.Pool.ExternalAssetBalance.Quo(sdk.NewUint(1000))

	msg := clptypes.MsgAddLiquidity{
		Signer:              key.GetAddress().String(),
		ExternalAsset:       &clptypes.Asset{Symbol: "ceth"},
		NativeAssetAmount:   nativeAdd,
		ExternalAssetAmount: externalAdd,
	}

	if _, err := buildAndBroadcast(clientCtx, txf, key, &msg); err != nil {
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

	// calculate expected result
	newPoolUnits, lpUnits, err := clpkeeper.CalculatePoolUnits(
		poolBefore.Pool.PoolUnits,
		poolBefore.Pool.NativeAssetBalance,
		poolBefore.Pool.ExternalAssetBalance,
		msg.NativeAssetAmount,
		msg.ExternalAssetAmount,
		18,
		sdk.NewDecWithPrec(5, 5),
		sdk.NewDecWithPrec(5, 4))

	if !poolAfter.Pool.PoolUnits.Equal(newPoolUnits) {
		return errors.New(fmt.Sprintf("pool unit mismatch (expected: %s after: %s)", newPoolUnits.String(), poolAfter.Pool.PoolUnits.String()))
	}

	lp, err := clpQueryClient.GetLiquidityProvider(context.Background(), &clptypes.LiquidityProviderReq{
		Symbol:    "ceth",
		LpAddress: key.GetAddress().String(),
	})
	if err != nil {
		return err
	}

	if !lp.LiquidityProvider.LiquidityProviderUnits.Sub(lpBefore.LiquidityProvider.LiquidityProviderUnits).Equal(lpUnits) {
		return errors.New(fmt.Sprintf("liquidity provided unit mismatch (expected: %s received: %s)",
			lpUnits.String(),
			lp.LiquidityProvider.LiquidityProviderUnits.String()),
		)
	}

	return nil
}

func TestOpenPosition(clientCtx client.Context, txf tx.Factory, key keyring.Info) error {
	msg := types.MsgOpen{
		Signer:           key.GetAddress().String(),
		CollateralAsset:  "rowan",
		CollateralAmount: sdk.NewUint(100),
		BorrowAsset:      "ceth",
		Position:         types.Position_LONG,
	}

	res, err := buildAndBroadcast(clientCtx, txf, key, &msg)
	if err != nil {
		panic(err)
	}

	log.Print(res)

	return err
}

func buildAndBroadcast(clientCtx client.Context, txf tx.Factory, key keyring.Info, msg sdk.Msg) (*sdk.TxResponse, error) {
	txb, err := tx.BuildUnsignedTx(txf, msg)
	if err != nil {
		return nil, err
	}

	err = tx.Sign(txf, key.GetName(), txb, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, err
	}

	res, err := clientCtx.
		WithSimulation(true).
		WithBroadcastMode("block").
		BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	return res, err
}
