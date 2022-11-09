package main

import (
	"context"
	"fmt"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetVerifyAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add --height --from --nativeAmount --externalAmount --external-asset",
		Short: "Verify a liquidity add",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("verifying add...\n")
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			nativeAmount := sdk.NewUintFromString(viper.GetString("nativeAmount"))
			externalAmount := sdk.NewUintFromString(viper.GetString("externalAmount"))

			err = VerifyAdd(clientCtx,
				viper.GetString("from"),
				viper.GetUint64("height"),
				nativeAmount,
				externalAmount,
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
	cmd.Flags().String("nativeAmount", "0", "native amount added")
	cmd.Flags().String("externalAmount", "0", "external amount added")
	cmd.Flags().String("external-asset", "", "external asset of pool")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("nativeAmount")
	_ = cmd.MarkFlagRequired("externalAmount")
	_ = cmd.MarkFlagRequired("external-asset")
	_ = cmd.MarkFlagRequired("height")
	return cmd
}

func VerifyAdd(clientCtx client.Context, from string, height uint64, nativeAmount, externalAmount sdk.Uint, externalAsset string) error {
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
		// Use empty LP if this is the first add
		lpBefore = &clptypes.LiquidityProviderRes{
			LiquidityProvider: &clptypes.LiquidityProvider{
				LiquidityProviderUnits: sdk.ZeroUint(),
			},
		}
	}

	// Lookup pool balances before remove
	poolBefore, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: externalAsset})
	if err != nil {
		return err
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
		if tokenParam.Asset == externalAsset {
			buyNativeSwapFeeRate = tokenParam.SwapFeeRate
			break
		}
	}

	// Calculate expected values
	nativeAssetDepth := poolBefore.Pool.NativeAssetBalance.Add(poolBefore.Pool.NativeLiabilities)
	externalAssetDepth := poolBefore.Pool.ExternalAssetBalance.Add(poolBefore.Pool.ExternalLiabilities)
	_ /*newPoolUnits*/, lpUnits, _, _, err := clpkeeper.CalculatePoolUnits(
		poolBefore.Pool.PoolUnits,
		nativeAssetDepth,
		externalAssetDepth,
		nativeAmount,
		externalAmount,
		sellNativeSwapFeeRate,
		buyNativeSwapFeeRate,
		pmtpParams.PmtpRateParams.PmtpCurrentRunningRate)
	if err != nil {
		return err
	}

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

	// Verify LP units are increased by lpUnits
	// Verify native balance is deducted by nativeAmount
	// Verify external balance is deducted by externalAmount
	externalDiff := extAfter.Balance.Amount.Sub(extBefore.Balance.Amount)
	nativeDiff := rowanAfter.Balance.Amount.Sub(rowanBefore.Balance.Amount)
	lpUnitsBeforeInt := sdk.NewIntFromBigInt(lpBefore.LiquidityProvider.LiquidityProviderUnits.BigInt())
	lpUnitsAfterInt := sdk.NewIntFromBigInt(lpAfter.LiquidityProvider.LiquidityProviderUnits.BigInt())
	lpUnitsDiff := lpUnitsAfterInt.Sub(lpUnitsBeforeInt)

	fmt.Printf("\nWallet native balance before %s\n", rowanBefore.Balance.Amount.String())
	fmt.Printf("Wallet external balance before %s\n\n", extBefore.Balance.Amount.String())

	fmt.Printf("Wallet native balance after %s \n", rowanAfter.Balance.Amount.String())
	fmt.Printf("Wallet external balance after %s \n", extAfter.Balance.Amount.String())

	fmt.Printf("\nWallet native diff %s (expected: %s unexpected: %s)\n",
		nativeDiff.String(),
		sdk.NewIntFromBigInt(nativeAmount.BigInt()).Neg().String(),
		nativeDiff.Sub(sdk.NewIntFromBigInt(nativeAmount.BigInt()).Neg()).String())
	fmt.Printf("Wallet external diff %s (expected: %s unexpected: %s)\n",
		externalDiff.String(),
		sdk.NewIntFromBigInt(externalAmount.BigInt()).Neg().String(),
		externalDiff.Sub(sdk.NewIntFromBigInt(externalAmount.BigInt()).Neg()))

	fmt.Printf("\nLP units before %s \n", lpBefore.LiquidityProvider.LiquidityProviderUnits.String())
	fmt.Printf("LP units after %s \n", lpAfter.LiquidityProvider.LiquidityProviderUnits.String())
	fmt.Printf("LP units diff %s (expected: %s unexpected: %s)\n", lpUnitsDiff.String(), lpUnits.String(), lpUnitsDiff.Sub(sdk.NewIntFromBigInt(lpUnits.BigInt())))

	clpQueryClient = clptypes.NewQueryClient(clientCtx.WithHeight(int64(height)))
	poolAfter, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: externalAsset})
	if err != nil {
		return err
	}

	fmt.Printf("\nPool units before %s\n", poolBefore.Pool.PoolUnits.String())
	fmt.Printf("Pool units after %s\n", poolAfter.Pool.PoolUnits.String())
	fmt.Printf("Pool units diff %s\n", sdk.NewIntFromBigInt(poolAfter.Pool.PoolUnits.BigInt()).Sub(sdk.NewIntFromBigInt(poolBefore.Pool.PoolUnits.BigInt())))

	lpUnitsBeforeDec := sdk.NewDecFromBigInt(lpBefore.LiquidityProvider.LiquidityProviderUnits.BigInt())
	lpUnitsAfterDec := sdk.NewDecFromBigInt(lpAfter.LiquidityProvider.LiquidityProviderUnits.BigInt())
	poolUnitsBeforeDec := sdk.NewDecFromBigInt(poolBefore.Pool.PoolUnits.BigInt())
	poolUnitsAfterDec := sdk.NewDecFromBigInt(poolAfter.Pool.PoolUnits.BigInt())
	poolShareBefore := lpUnitsBeforeDec.Quo(poolUnitsBeforeDec)
	poolShareAfter := lpUnitsAfterDec.Quo(poolUnitsAfterDec)

	fmt.Printf("\nPool share before %s\n", poolShareBefore.String())
	fmt.Printf("Pool share after %s\n", poolShareAfter.String())

	fmt.Printf("\nPool external balance before %s\n", poolBefore.Pool.ExternalAssetBalance.String())
	fmt.Printf("Pool native balance before %s\n", poolBefore.Pool.NativeAssetBalance.String())

	fmt.Printf("\nPool external balance after %s\n", poolAfter.Pool.ExternalAssetBalance.String())
	fmt.Printf("Pool native balance after %s\n", poolAfter.Pool.NativeAssetBalance.String())

	poolExternalDiff := sdk.NewIntFromBigInt(poolAfter.Pool.ExternalAssetBalance.BigInt()).Sub(sdk.NewIntFromBigInt(poolBefore.Pool.ExternalAssetBalance.BigInt()))
	poolNativeDiff := sdk.NewIntFromBigInt(poolAfter.Pool.NativeAssetBalance.BigInt()).Sub(sdk.NewIntFromBigInt(poolBefore.Pool.NativeAssetBalance.BigInt()))

	fmt.Printf("\nPool external balance diff %s (expected: %s)\n", poolExternalDiff.String(), externalAmount.String())
	fmt.Printf("Pool native balance diff %s (expected: %s)\n", poolNativeDiff.String(), nativeAmount.String())

	return nil
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
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("units")
	_ = cmd.MarkFlagRequired("external-asset")
	_ = cmd.MarkFlagRequired("height")
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
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, _ /*lpUnitsLeft*/ := clpkeeper.CalculateWithdrawalFromUnits(poolBefore.Pool.PoolUnits,
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

	poolAfter, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: externalAsset})
	if err != nil {
		return err
	}

	// Verify LP units are reduced by --units
	// Verify native received amount
	// Verify external received amount
	externalDiff := extAfter.Balance.Amount.Sub(extBefore.Balance.Amount)
	nativeDiff := rowanAfter.Balance.Amount.Sub(rowanBefore.Balance.Amount)
	lpUnitsBeforeInt := sdk.NewIntFromBigInt(lpBefore.LiquidityProvider.LiquidityProviderUnits.BigInt())
	lpUnitsAfterInt := sdk.NewIntFromBigInt(lpAfter.LiquidityProvider.LiquidityProviderUnits.BigInt())
	lpUnitsDiff := lpUnitsAfterInt.Sub(lpUnitsBeforeInt)

	fmt.Printf("\nWallet native balance before %s\n", rowanBefore.Balance.Amount.String())
	fmt.Printf("Wallet external balance before %s\n\n", extBefore.Balance.Amount.String())

	fmt.Printf("Wallet native balance after %s \n", rowanAfter.Balance.Amount.String())
	fmt.Printf("Wallet external balance after %s \n", extAfter.Balance.Amount.String())

	fmt.Printf("\nWallet native diff %s (expected: %s unexpected: %s)\n",
		nativeDiff.String(),
		sdk.NewIntFromBigInt(withdrawNativeAssetAmount.BigInt()).String(),
		nativeDiff.Sub(sdk.NewIntFromBigInt(withdrawNativeAssetAmount.BigInt())).String())
	fmt.Printf("Wallet external diff %s (expected: %s unexpected: %s)\n",
		externalDiff.String(),
		sdk.NewIntFromBigInt(withdrawExternalAssetAmount.BigInt()).String(),
		externalDiff.Sub(sdk.NewIntFromBigInt(withdrawExternalAssetAmount.BigInt())))

	fmt.Printf("\nLP units before %s \n", lpBefore.LiquidityProvider.LiquidityProviderUnits.String())
	fmt.Printf("LP units after %s \n", lpAfter.LiquidityProvider.LiquidityProviderUnits.String())
	fmt.Printf("LP units diff %s (expected: -%s)\n", lpUnitsDiff.String(), units.String())

	lpUnitsBeforeDec := sdk.NewDecFromBigInt(lpBefore.LiquidityProvider.LiquidityProviderUnits.BigInt())
	lpUnitsAfterDec := sdk.NewDecFromBigInt(lpAfter.LiquidityProvider.LiquidityProviderUnits.BigInt())
	poolUnitsBeforeDec := sdk.NewDecFromBigInt(poolBefore.Pool.PoolUnits.BigInt())
	poolUnitsAfterDec := sdk.NewDecFromBigInt(poolAfter.Pool.PoolUnits.BigInt())
	poolShareBefore := lpUnitsBeforeDec.Quo(poolUnitsBeforeDec)
	poolShareAfter := lpUnitsAfterDec.Quo(poolUnitsAfterDec)

	fmt.Printf("\nPool share before %s\n", poolShareBefore.String())
	fmt.Printf("Pool share after %s\n", poolShareAfter.String())

	return nil
}

func GetVerifyClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close --height --from --id",
		Short: "Verify a margin position close",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("verifying close...\n")
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			err = VerifyClose(clientCtx,
				viper.GetString("from"),
				int64(viper.GetUint64("height")),
				viper.GetUint64("id"))
			if err != nil {
				panic(err)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	//cmd.Flags().Uint64("height", 0, "height of transaction")
	cmd.Flags().String("from", "", "address of transactor")
	cmd.Flags().Uint64("id", 0, "id of mtp")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("height")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func VerifyClose(clientCtx client.Context, from string, height int64, id uint64) error {
	// Lookup MTP
	marginQueryClient := types.NewQueryClient(clientCtx.WithHeight(height - 1))
	mtpResponse, err := marginQueryClient.GetMTP(context.Background(), &types.MTPRequest{
		Address: from,
		Id:      id,
	})
	if err != nil {
		return sdkerrors.Wrap(err, fmt.Sprintf("error looking up mtp at height %d", height-1))
	}
	fmt.Printf("\nMTP custody %s (%s)\n", mtpResponse.Mtp.CustodyAmount.String(), mtpResponse.Mtp.CustodyAsset)
	fmt.Printf("MTP collateral %s (%s)\n", mtpResponse.Mtp.CollateralAmount.String(), mtpResponse.Mtp.CollateralAsset)
	fmt.Printf("MTP leverage %s\n", mtpResponse.Mtp.Leverage.String())
	fmt.Printf("MTP liability %s\n", mtpResponse.Mtp.Liabilities.String())
	fmt.Printf("MTP health %s\n", mtpResponse.Mtp.MtpHealth)
	fmt.Printf("MTP interest paid custody %s\n", mtpResponse.Mtp.InterestPaidCustody.String())
	fmt.Printf("MTP interest paid collateral %s\n", mtpResponse.Mtp.InterestPaidCollateral.String())
	fmt.Printf("MTP interest unpaid collateral %s\n", mtpResponse.Mtp.InterestUnpaidCollateral.String())

	// lookup wallet before
	bankQueryClient := banktypes.NewQueryClient(clientCtx.WithHeight(height - 1))
	collateralBefore, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   mtpResponse.Mtp.CollateralAsset,
	})
	if err != nil {
		return err
	}
	custodyBefore, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   mtpResponse.Mtp.CustodyAsset,
	})
	if err != nil {
		return err
	}
	fmt.Printf("\nWallet collateral balance before: %s\n", collateralBefore.Balance.Amount.String())
	fmt.Printf("Wallet custody balance before: %s\n\n", custodyBefore.Balance.Amount.String())
	// Ensure mtp does not exist after close
	marginQueryClient = types.NewQueryClient(clientCtx.WithHeight(height))
	_, err = marginQueryClient.GetMTP(context.Background(), &types.MTPRequest{
		Address: from,
		Id:      id,
	})
	if err != nil {
		fmt.Printf("confirmed MTP does not exist at close height %d\n\n", height)
	} else {
		return sdkerrors.Wrap(err, fmt.Sprintf("error: found mtp at close height %d", height))
	}

	var externalAsset string
	if types.StringCompare(mtpResponse.Mtp.CollateralAsset, "rowan") {
		externalAsset = mtpResponse.Mtp.CustodyAsset
	} else {
		externalAsset = mtpResponse.Mtp.CollateralAsset
	}

	clpQueryClient := clptypes.NewQueryClient(clientCtx.WithHeight(height - 1))
	poolBefore, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: externalAsset})
	if err != nil {
		return err
	}
	fmt.Printf("\nPool health before %s\n", poolBefore.Pool.Health.String())
	fmt.Printf("Pool native custody before %s\n", poolBefore.Pool.NativeCustody.String())
	fmt.Printf("Pool external custody before %s\n", poolBefore.Pool.ExternalCustody.String())
	fmt.Printf("Pool native liabilities before %s\n", poolBefore.Pool.NativeLiabilities.String())
	fmt.Printf("Pool external liabilities before %s\n", poolBefore.Pool.ExternalLiabilities.String())
	fmt.Printf("Pool native depth (including liabilities) before %s\n", poolBefore.Pool.NativeAssetBalance.Add(poolBefore.Pool.NativeLiabilities).String())
	fmt.Printf("Pool external depth (including liabilities) before %s\n", poolBefore.Pool.ExternalAssetBalance.Add(poolBefore.Pool.ExternalLiabilities).String())

	clpQueryClient = clptypes.NewQueryClient(clientCtx.WithHeight(height))
	poolAfter, err := clpQueryClient.GetPool(context.Background(), &clptypes.PoolReq{Symbol: externalAsset})
	if err != nil {
		return err
	}

	expectedPoolNativeCustody := sdk.NewIntFromBigInt(poolBefore.Pool.NativeCustody.BigInt())
	expectedPoolExternalCustody := sdk.NewIntFromBigInt(poolBefore.Pool.ExternalCustody.BigInt())
	expectedPoolNativeLiabilities := sdk.NewIntFromBigInt(poolBefore.Pool.NativeLiabilities.BigInt())
	expectedPoolExternalLiabilities := sdk.NewIntFromBigInt(poolBefore.Pool.ExternalLiabilities.BigInt())
	if types.StringCompare(mtpResponse.Mtp.CustodyAsset, "rowan") {
		expectedPoolNativeCustody = expectedPoolNativeCustody.Sub(
			sdk.NewIntFromBigInt(mtpResponse.Mtp.CustodyAmount.BigInt()),
		)
		expectedPoolExternalLiabilities = expectedPoolExternalLiabilities.Sub(
			sdk.NewIntFromBigInt(mtpResponse.Mtp.Liabilities.BigInt()),
		)
	} else {
		expectedPoolExternalCustody = expectedPoolExternalCustody.Sub(
			sdk.NewIntFromBigInt(mtpResponse.Mtp.CustodyAmount.BigInt()),
		)
		expectedPoolNativeLiabilities = expectedPoolNativeLiabilities.Sub(
			sdk.NewIntFromBigInt(mtpResponse.Mtp.Liabilities.BigInt()),
		)
	}

	fmt.Printf("\nPool health after %s\n", poolAfter.Pool.Health.String())
	fmt.Printf("Pool native custody after %s (expected %s)\n", poolAfter.Pool.NativeCustody.String(), expectedPoolNativeCustody.String())
	fmt.Printf("Pool external custody after %s (expected %s)\n", poolAfter.Pool.ExternalCustody.String(), expectedPoolExternalCustody.String())
	fmt.Printf("Pool native liabilities after %s (expected %s)\n", poolAfter.Pool.NativeLiabilities.String(), expectedPoolNativeLiabilities.String())
	fmt.Printf("Pool external liabilities after %s (expected %s)\n", poolAfter.Pool.ExternalLiabilities.String(), expectedPoolExternalLiabilities.String())

	// Final interest payment
	//finalInterest := marginkeeper.CalcMTPInterestLiabilities(mtpResponse.Mtp, pool.Pool.InterestRate, 0, 1)
	//mtpCustodyAmount := mtpResponse.Mtp.CustodyAmount.Sub(finalInterest)
	// get swap params
	clpQueryClient = clptypes.NewQueryClient(clientCtx.WithHeight(height - 1))
	pmtpParams, err := clpQueryClient.GetPmtpParams(context.Background(), &clptypes.PmtpParamsReq{})
	if err != nil {
		return err
	}
	// Calculate expected return
	// Swap custody
	swapFeeParams, err := clpQueryClient.GetSwapFeeParams(context.Background(), &clptypes.SwapFeeParamsReq{})
	if err != nil {
		return err
	}

	nativeAsset := types.GetSettlementAsset()
	pool := *poolBefore.Pool

	if types.StringCompare(mtpResponse.Mtp.CustodyAsset, nativeAsset) {
		pool.NativeCustody = pool.NativeCustody.Sub(mtpResponse.Mtp.CustodyAmount)
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(mtpResponse.Mtp.CustodyAmount)
	} else {
		pool.ExternalCustody = pool.ExternalCustody.Sub(mtpResponse.Mtp.CustodyAmount)
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(mtpResponse.Mtp.CustodyAmount)
	}
	X, Y, toRowan, _ := pool.ExtractValues(clptypes.Asset{Symbol: mtpResponse.Mtp.CollateralAsset})
	X, Y = pool.ExtractDebt(X, Y, toRowan)
	repayAmount, _ := clpkeeper.CalcSwapResult(toRowan, X, mtpResponse.Mtp.CustodyAmount, Y, pmtpParams.PmtpRateParams.PmtpCurrentRunningRate, swapFeeParams.DefaultSwapFeeRate)

	// Repay()
	mtp := mtpResponse.Mtp
	// nolint:staticcheck,ineffassign
	returnAmount, debtP, debtI := sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint()
	Liabilities := mtp.Liabilities
	InterestUnpaidCollateral := mtp.InterestUnpaidCollateral

	have := repayAmount
	owe := Liabilities.Add(InterestUnpaidCollateral)

	if have.LT(Liabilities) {
		//can't afford principle liability
		returnAmount = sdk.ZeroUint()
		debtP = Liabilities.Sub(have)
		debtI = InterestUnpaidCollateral
	} else if have.LT(owe) {
		// v principle liability; x excess liability
		returnAmount = sdk.ZeroUint()
		debtP = sdk.ZeroUint()
		debtI = Liabilities.Add(InterestUnpaidCollateral).Sub(have)
	} else {
		// can afford both
		returnAmount = have.Sub(Liabilities).Sub(InterestUnpaidCollateral)
		debtP = sdk.ZeroUint()
		debtI = sdk.ZeroUint()
	}

	fmt.Printf("\nReturn amount: %s\n", returnAmount.String())
	fmt.Printf("Loss: %s\n\n", debtP.Add(debtI).String())

	// lookup wallet balances after close
	bankQueryClient = banktypes.NewQueryClient(clientCtx.WithHeight(height))
	collateralAfter, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   mtpResponse.Mtp.CollateralAsset,
	})
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	custodyAfter, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: from,
		Denom:   mtpResponse.Mtp.CustodyAsset,
	})
	if err != nil {
		return err
	}
	collateralDiff := collateralAfter.Balance.Amount.Sub(collateralBefore.Balance.Amount)
	custodyDiff := custodyAfter.Balance.Amount.Sub(custodyBefore.Balance.Amount)
	fmt.Printf("Wallet collateral (%s) balance after: %s (diff: %s expected diff: %s)\n",
		mtpResponse.Mtp.CollateralAsset,
		collateralAfter.Balance.Amount.String(),
		collateralDiff.String(),
		returnAmount.String())
	fmt.Printf("Wallet custody (%s) balance after: %s (diff: %s)\n\n", mtpResponse.Mtp.CustodyAsset, custodyAfter.Balance.Amount.String(), custodyDiff.String())

	return nil
}
