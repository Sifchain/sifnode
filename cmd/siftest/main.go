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
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Use: "siftest",
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
		//RunE: run,
	}

	flags.AddTxFlagsToCmd(rootCmd)
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify transaction results",
	}

	verifyCmd.AddCommand(GetVerifyRemove())

	rootCmd.AddCommand(verifyCmd)

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

	//txf = txf.WithAccountNumber(accountNumber).WithSequence(seq)
	//err = TestAddLiquidity(clientCtx, txf, key)
	//if err != nil {
	//	panic(err)
	//}

	txf = txf.WithAccountNumber(accountNumber).WithSequence(seq)
	err = TestSwap(clientCtx, txf, key)
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

func TestSwap(clientCtx client.Context, txf tx.Factory, key keyring.Info) error {

	bankQueryClient := banktypes.NewQueryClient(clientCtx)
	cethBefore, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: key.GetAddress().String(),
		Denom:   "ceth",
	})
	if err != nil {
		return err
	}
	rowanBefore, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: key.GetAddress().String(),
		Denom:   "rowan",
	})
	if err != nil {
		return err
	}

	clpQueryClient := clptypes.NewQueryClient(clientCtx)
	poolBefore, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: "ceth"})
	if err != nil {
		return err
	}

	msg := clptypes.MsgSwap{
		Signer:             key.GetAddress().String(),
		SentAsset:          &clptypes.Asset{Symbol: "ceth"},
		ReceivedAsset:      &clptypes.Asset{Symbol: "rowan"},
		SentAmount:         sdk.NewUint(10000),
		MinReceivingAmount: sdk.NewUint(0),
	}

	if _, err := buildAndBroadcast(clientCtx, txf, key, &msg); err != nil {
		return err
	}

	cethAfter, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: key.GetAddress().String(),
		Denom:   "ceth",
	})
	rowanAfter, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: key.GetAddress().String(),
		Denom:   "rowan",
	})
	poolAfter, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: "ceth"})
	if err != nil {
		return err
	}

	rowanDiff := rowanAfter.Balance.Amount.Sub(rowanBefore.Balance.Amount)
	// negative
	cethDiff := cethAfter.Balance.Amount.Sub(cethBefore.Balance.Amount)
	// negative
	poolNativeDiff := poolBefore.Pool.NativeAssetBalance.Sub(poolAfter.Pool.NativeAssetBalance)
	poolExternalDiff := poolAfter.Pool.ExternalAssetBalance.Sub(poolBefore.Pool.ExternalAssetBalance)

	fmt.Printf("Pool sent diff: %s\n", poolNativeDiff.String())
	fmt.Printf("Pool received diff: %s\n", poolExternalDiff.String())
	fmt.Printf("Address received diff: %s\n", rowanDiff.String())
	fmt.Printf("Address sent diff: %s\n", cethDiff.String())

	return nil
}

/* VerifySwap verifies amounts sent and received from wallet address.
 */
func VerifySwap(clientCtx client.Context, key keyring.Info) {

}

func GetVerifyRemove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove --height --from --units --external-asset",
		Short: "Verify a removal",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("verifying removal...\n")
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			unitsRemoved := sdk.NewUintFromString(viper.GetString("units"))

			err = VerifyRemove(clientCtx,
				viper.GetString("from"),
				viper.GetUint64("height"),
				unitsRemoved,
				viper.GetString("external-asset"))
			if err != nil {
				panic(err)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	//cmd.Flags().Uint64("height", 0, "height of transaction")
	cmd.Flags().String("from", "", "address of transactor")
	cmd.Flags().String("units", "0", "number of units removed")
	cmd.Flags().String("external-asset", "", "external asset of pool")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("units")
	cmd.MarkFlagRequired("external-asset")
	cmd.MarkFlagRequired("height")
	return cmd
}

/*
	 VerifyRemove verifies amounts received after remove.
		--height --from --units --external-asset
*/
func VerifyRemove(clientCtx client.Context, from string, height uint64, units sdk.Uint, externalAsset string) error {
	// Lookup wallet balances before remove
	// Lookup wallet balances after remove
	bankQueryClient := banktypes.NewQueryClient(clientCtx.WithHeight(int64(height - 1)))
	extBefore, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   externalAsset,
	})
	if err != nil {
		return err
	}
	rowanBefore, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   "rowan",
	})
	if err != nil {
		return err
	}

	// Lookup LP units before remove
	// Lookup LP units after remove
	clpQueryClient := clptypes.NewQueryClient(clientCtx.WithHeight(int64(height - 1)))
	lpBefore, err := clpQueryClient.GetLiquidityProvider(context.Background(), &clptypes.LiquidityProviderReq{
		Symbol:    externalAsset,
		LpAddress: from,
	})
	if err != nil {
		return err
	}

	// Lookup pool balances before remove
	poolBefore, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: externalAsset})
	if err != nil {
		return err
	}

	// Calculate expected values
	nativeAssetDepth := poolBefore.Pool.NativeAssetBalance.Add(poolBefore.Pool.NativeLiabilities)
	externalAssetDepth := poolBefore.Pool.ExternalAssetBalance.Add(poolBefore.Pool.ExternalLiabilities)
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft := clpkeeper.CalculateWithdrawalFromUnits(poolBefore.Pool.PoolUnits,
		nativeAssetDepth.String(), externalAssetDepth.String(), lpBefore.LiquidityProvider.LiquidityProviderUnits.String(),
		units)

	// Lookup wallet balances after
	bankQueryClient = banktypes.NewQueryClient(clientCtx.WithHeight(int64(height)))
	extAfter, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   externalAsset,
	})
	if err != nil {
		return err
	}
	rowanAfter, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   "rowan",
	})
	if err != nil {
		return err
	}

	// Lookup LP after
	clpQueryClient = clptypes.NewQueryClient(clientCtx.WithHeight(int64(height)))
	lpAfter, err := clpQueryClient.GetLiquidityProvider(context.Background(), &clptypes.LiquidityProviderReq{
		Symbol:    externalAsset,
		LpAddress: from,
	})
	if err != nil {
		lpAfter = &clptypes.LiquidityProviderRes{
			LiquidityProvider: &clptypes.LiquidityProvider{
				LiquidityProviderUnits: sdk.ZeroUint(),
			},
		}
	}

	// Verify LP units are reduced by --units
	// Verify native received amount
	// Verify external received amount
	//fee, _ := sdk.NewIntFromString("1000000000000000000")
	externalDiff := extAfter.Balance.Amount.Sub(extBefore.Balance.Amount)
	nativeDiff := rowanAfter.Balance.Amount.Sub(rowanBefore.Balance.Amount)

	fmt.Printf("External received %s \n", externalDiff.String())
	fmt.Printf("External expected %s \n\n", withdrawExternalAssetAmount.String())

	fmt.Printf("Native received %s \n", nativeDiff.String())
	//fmt.Printf("Native received excluding fee deduction %s \n", nativeDiff.Add(fee).String())
	fmt.Printf("Native expected %s \n", withdrawNativeAssetAmount.String())
	fmt.Printf("Native expected - received %s \n\n", sdk.NewIntFromBigInt(withdrawNativeAssetAmount.BigInt()).Sub(nativeDiff).String())

	//fmt.Printf("LP units before %s \n", lpBefore.LiquidityProvider.LiquidityProviderUnits.String())
	fmt.Printf("LP units after %s \n", lpAfter.LiquidityProvider.LiquidityProviderUnits.String())
	fmt.Printf("LP units expected after %s \n", lpUnitsLeft.String())

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
