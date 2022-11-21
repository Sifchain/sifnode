package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
)

func runTest(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	key, err := txf.Keybase().Key(clientCtx.GetFromName())
	if err != nil {
		return err
	}

	accountNumber, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, key.GetAddress())
	if err != nil {
		panic(err)
	}

	txf = txf.WithAccountNumber(accountNumber).WithSequence(seq)
	err = TestOpenPosition(clientCtx, txf, key)
	if err != nil {
		panic(err)
	}

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

	pmtpParams, err := clpQueryClient.GetPmtpParams(context.Background(), &clptypes.PmtpParamsReq{})
	if err != nil {
		return err
	}

	swapFeeParams, err := clpQueryClient.GetSwapFeeParams(context.Background(), &clptypes.SwapFeeParamsReq{})
	if err != nil {
		return err
	}

	// get native swap fee rate
	sellNativeSwapFeeRate := swapFeeParams.DefaultSwapFeeRate
	for _, tokenParam := range swapFeeParams.TokenParams {
		if tokenParam.Asset == clptypes.GetSettlementAsset().Symbol {
			sellNativeSwapFeeRate = tokenParam.SwapFeeRate
			break
		}
	}

	// get external token swap fee rate
	buyNativeSwapFeeRate := swapFeeParams.DefaultSwapFeeRate
	for _, tokenParam := range swapFeeParams.TokenParams {
		if tokenParam.Asset == msg.ExternalAsset.Symbol {
			buyNativeSwapFeeRate = tokenParam.SwapFeeRate
			break
		}
	}

	// calculate expected result
	newPoolUnits, lpUnits, _, _, err := clpkeeper.CalculatePoolUnits(
		poolBefore.Pool.PoolUnits,
		poolBefore.Pool.NativeAssetBalance,
		poolBefore.Pool.ExternalAssetBalance,
		msg.NativeAssetAmount,
		msg.ExternalAssetAmount,
		sellNativeSwapFeeRate,
		buyNativeSwapFeeRate,
		pmtpParams.PmtpRateParams.PmtpCurrentRunningRate)
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}
	rowanAfter, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: key.GetAddress().String(),
		Denom:   "rowan",
	})
	if err != nil {
		return err
	}
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

func TestOpenPosition(clientCtx client.Context, txf tx.Factory, key keyring.Info) error {
	msg := types.MsgOpen{
		Signer:           key.GetAddress().String(),
		CollateralAsset:  "rowan",
		CollateralAmount: sdk.NewUint(100),
		BorrowAsset:      "ceth",
		Position:         types.Position_LONG,
		Leverage:         sdk.NewDec(2),
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
