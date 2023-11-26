package keeper_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	keeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAfterEpochEnd_DistributeToLPWallets(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper

	// distribute rewards to lp addresses
	rewardsParams := types.GetDefaultRewardParams()
	rewardsParams.RewardsDistribute = true
	clpKeeper.SetRewardParams(ctx, rewardsParams)

	asset := types.NewAsset("cusdc")

	// Create a RewardsBucket with a denom and amount
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  asset.Symbol,
		Amount: initialAmount,
	}
	clpKeeper.SetRewardsBucket(ctx, rewardsBucket)

	// mint coins to the module account
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(asset.Symbol, initialAmount),
	))
	require.NoError(t, err)

	// check module balance
	moduleBalance := app.BankKeeper.GetBalance(ctx, types.GetCLPModuleAddress(), asset.Symbol)
	require.Equal(t, initialAmount, moduleBalance.Amount)

	// Create a liquidity provider
	lpAddress, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
	require.NoError(t, err)
	lp := types.NewLiquidityProvider(&asset, sdk.NewUint(1), lpAddress, ctx.BlockHeight()-int64(rewardsParams.RewardsLockPeriod)-1)
	clpKeeper.SetLiquidityProvider(ctx, &lp)

	// check account balance
	preBalance := app.BankKeeper.GetBalance(ctx, lpAddress, asset.Symbol)
	require.Equal(t, sdk.ZeroInt(), preBalance.Amount)

	clpKeeper.AfterEpochEnd(ctx, rewardsParams.RewardsEpochIdentifier, 1)

	// check account balance
	postBalance := app.BankKeeper.GetBalance(ctx, lpAddress, asset.Symbol)
	require.Equal(t, preBalance.Add(sdk.NewCoin(asset.Symbol, initialAmount)), postBalance)
}

func TestAfterEpochEnd_AddToLiquidityPool(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper

	// add to liquidity pool
	rewardsParams := types.GetDefaultRewardParams()
	rewardsParams.RewardsDistribute = false
	clpKeeper.SetRewardParams(ctx, rewardsParams)

	// define signer
	signer := test.GenerateAddress(test.AddressKey1)

	// define the pool external asset
	asset := types.NewAsset("cusdc")

	// define the amount to add to the pool
	nativeAssetAmount := sdk.NewUintFromString("100")
	externalAssetAmount := sdk.NewUintFromString("100")

	// define the user initial balance
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUintFromString("10000")))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUintFromString("10000")))

	// add balance to signer wallet
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	// create pool message
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)

	// execute message
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	require.NoError(t, err)
	require.Equal(t, types.Pool{
		ExternalAsset: &types.Asset{
			Symbol: "cusdc",
		},
		NativeAssetBalance:            sdk.NewUintFromString("100"),
		ExternalAssetBalance:          sdk.NewUintFromString("100"),
		PoolUnits:                     sdk.NewUintFromString("1"),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
		ExternalLiabilities:           sdk.ZeroUint(),
		ExternalCustody:               sdk.ZeroUint(),
		NativeLiabilities:             sdk.ZeroUint(),
		NativeCustody:                 sdk.ZeroUint(),
		Health:                        sdk.ZeroDec(),
		InterestRate:                  sdk.ZeroDec(),
		UnsettledExternalLiabilities:  sdk.ZeroUint(),
		UnsettledNativeLiabilities:    sdk.ZeroUint(),
		BlockInterestNative:           sdk.ZeroUint(),
		BlockInterestExternal:         sdk.ZeroUint(),
		RewardAmountExternal:          sdk.ZeroUint(),
	}, *pool)

	// create liquidity provider
	lp := clpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), signer)
	require.Equal(t, types.LiquidityProvider{
		Asset:                    &asset,
		LiquidityProviderUnits:   sdk.NewUint(1),
		LiquidityProviderAddress: signer.String(),
		Unlocks:                  nil,
		LastUpdatedBlock:         ctx.BlockHeight(),
		RewardAmount:             nil,
	}, lp)

	// set last updated block to be before the rewards lock period
	lp.LastUpdatedBlock = ctx.BlockHeight() - int64(rewardsParams.RewardsLockPeriod) - 1
	clpKeeper.SetLiquidityProvider(ctx, &lp)

	lps, err := clpKeeper.GetAllLiquidityProviders(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(lps))

	// Create a RewardsBucket with a denom and amount
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  asset.Symbol,
		Amount: initialAmount,
	}
	clpKeeper.SetRewardsBucket(ctx, rewardsBucket)

	// mint coins to the module account
	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(asset.Symbol, initialAmount),
	))
	require.NoError(t, err)

	// check module balance
	moduleBalance := app.BankKeeper.GetBalance(ctx, types.GetCLPModuleAddress(), asset.Symbol)
	require.Equal(t, sdk.NewIntFromBigInt(externalAssetAmount.BigInt()).Add(initialAmount), moduleBalance.Amount)

	// check account balance
	preBalance := app.BankKeeper.GetBalance(ctx, signer, asset.Symbol)
	require.Equal(t, externalCoin.Amount.Sub(sdk.NewIntFromBigInt(externalAssetAmount.BigInt())), preBalance.Amount)

	// check liquidity provider units
	lp, err = clpKeeper.GetLiquidityProvider(ctx, asset.Symbol, signer.String())
	require.NoError(t, err)
	require.Equal(t, sdk.NewUint(1), lp.LiquidityProviderUnits)

	clpKeeper.AfterEpochEnd(ctx, rewardsParams.RewardsEpochIdentifier, 1)

	// check if rewards bucket is empty
	rewardsBucket, found := clpKeeper.GetRewardsBucket(ctx, asset.Symbol)
	require.True(t, found)
	require.Equal(t, sdk.ZeroInt(), rewardsBucket.Amount)

	// because rewards is not distributed to LP wallets, the LP balance should not change
	postBalance := app.BankKeeper.GetBalance(ctx, signer, asset.Symbol)
	require.Equal(t, preBalance, postBalance)

	// check liquidity provider units
	lp, err = clpKeeper.GetLiquidityProvider(ctx, asset.Symbol, signer.String())
	require.NoError(t, err)
	require.Equal(t, sdk.NewUint(1), lp.LiquidityProviderUnits)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(asset.Symbol, initialAmount)), lp.RewardAmount)

	// check pool reward amount external
	poolUpdated, err := clpKeeper.GetPool(ctx, asset.Symbol)
	require.NoError(t, err)
	require.Equal(t, sdk.NewUintFromBigInt(initialAmount.BigInt()), poolUpdated.RewardAmountExternal)
}

func TestAfterEpochEnd_AddToLiquidityPoolWithMultipleLiquidityProviders(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper

	// add to liquidity pool
	rewardsParams := types.GetDefaultRewardParams()
	rewardsParams.RewardsDistribute = false
	clpKeeper.SetRewardParams(ctx, rewardsParams)

	// define signers
	signer1 := test.GenerateAddress(test.AddressKey1)
	signer2 := test.GenerateAddress(test.AddressKey2)

	// define the pool external asset
	asset := types.NewAsset("cusdc")

	// define the amount to add to the pool
	nativeAssetAmount := sdk.NewUintFromString("100")
	externalAssetAmount := sdk.NewUintFromString("100")

	// define the user initial balance
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUintFromString("10000")))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUintFromString("10000")))

	// add balance to signers wallet
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer1, sdk.NewCoins(externalCoin, nativeCoin))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer2, sdk.NewCoins(externalCoin, nativeCoin))

	// create pool message
	msgCreatePool := types.NewMsgCreatePool(signer1, asset, nativeAssetAmount, externalAssetAmount)

	// execute message
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(100), &msgCreatePool)
	require.NoError(t, err)
	require.Equal(t, types.Pool{
		ExternalAsset: &types.Asset{
			Symbol: "cusdc",
		},
		NativeAssetBalance:            sdk.NewUintFromString("100"),
		ExternalAssetBalance:          sdk.NewUintFromString("100"),
		PoolUnits:                     sdk.NewUintFromString("100"),
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
		ExternalLiabilities:           sdk.ZeroUint(),
		ExternalCustody:               sdk.ZeroUint(),
		NativeLiabilities:             sdk.ZeroUint(),
		NativeCustody:                 sdk.ZeroUint(),
		Health:                        sdk.ZeroDec(),
		InterestRate:                  sdk.ZeroDec(),
		UnsettledExternalLiabilities:  sdk.ZeroUint(),
		UnsettledNativeLiabilities:    sdk.ZeroUint(),
		BlockInterestNative:           sdk.ZeroUint(),
		BlockInterestExternal:         sdk.ZeroUint(),
		RewardAmountExternal:          sdk.ZeroUint(),
	}, *pool)

	// create liquidity provider 1
	lp1 := clpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(100), signer1)
	require.Equal(t, types.LiquidityProvider{
		Asset:                    &asset,
		LiquidityProviderUnits:   sdk.NewUint(100),
		LiquidityProviderAddress: signer1.String(),
		Unlocks:                  nil,
		LastUpdatedBlock:         ctx.BlockHeight(),
		RewardAmount:             nil,
	}, lp1)

	// create liquidity provider 2
	nativeAssetDepth, externalAssetDepth := pool.ExtractDebt(pool.NativeAssetBalance, pool.ExternalAssetBalance, false)
	pmtpCurrentRunningRate := clpKeeper.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	sellNativeSwapFeeRate := clpKeeper.GetSwapFeeRate(ctx, types.GetSettlementAsset(), false)
	buyNativeSwapFeeRate := clpKeeper.GetSwapFeeRate(ctx, asset, false)
	newPoolUnits, lpUnits, _, _, err := keeper.CalculatePoolUnits(
		pool.PoolUnits,
		nativeAssetDepth,
		externalAssetDepth,
		nativeAssetAmount,
		externalAssetAmount,
		sellNativeSwapFeeRate,
		buyNativeSwapFeeRate,
		pmtpCurrentRunningRate)
	require.NoError(t, err)
	msgAddLiquidity := types.NewMsgAddLiquidity(signer2, asset, nativeAssetAmount, externalAssetAmount)
	lp2ptr, err := clpKeeper.AddLiquidity(ctx, &msgAddLiquidity, *pool, newPoolUnits, lpUnits)
	require.NoError(t, err)
	lp2 := *lp2ptr

	// set last updated block to be before the rewards lock period
	lp2.LastUpdatedBlock = ctx.BlockHeight() - int64(rewardsParams.RewardsLockPeriod) - 1
	clpKeeper.SetLiquidityProvider(ctx, &lp2)

	lps, err := clpKeeper.GetAllLiquidityProviders(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(lps))

	// Create a RewardsBucket with a denom and amount
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  asset.Symbol,
		Amount: initialAmount,
	}
	clpKeeper.SetRewardsBucket(ctx, rewardsBucket)

	// mint coins to the module account
	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(asset.Symbol, initialAmount),
	))
	require.NoError(t, err)

	// check module balance
	moduleBalance := app.BankKeeper.GetBalance(ctx, types.GetCLPModuleAddress(), asset.Symbol)
	require.Equal(t, sdk.NewIntFromBigInt(externalAssetAmount.BigInt()).Mul(sdk.NewInt(2)).Add(initialAmount), moduleBalance.Amount)

	// check account balance 1
	preBalance1 := app.BankKeeper.GetBalance(ctx, signer1, asset.Symbol)
	require.Equal(t, externalCoin.Amount.Sub(sdk.NewIntFromBigInt(externalAssetAmount.BigInt())), preBalance1.Amount)

	// check account balance 2
	preBalance2 := app.BankKeeper.GetBalance(ctx, signer2, asset.Symbol)
	require.Equal(t, externalCoin.Amount.Sub(sdk.NewIntFromBigInt(externalAssetAmount.BigInt())), preBalance2.Amount)

	// check liquidity provider units 1
	lp1, err = clpKeeper.GetLiquidityProvider(ctx, asset.Symbol, signer1.String())
	require.NoError(t, err)
	require.Equal(t, sdk.NewUint(100), lp1.LiquidityProviderUnits)

	// check liquidity provider units 2
	lp2, err = clpKeeper.GetLiquidityProvider(ctx, asset.Symbol, signer2.String())
	require.NoError(t, err)
	require.Equal(t, sdk.NewUint(100), lp2.LiquidityProviderUnits)

	clpKeeper.AfterEpochEnd(ctx, rewardsParams.RewardsEpochIdentifier, 1)

	// because rewards is not distributed to LP wallets, the LP balance 1 should not change
	postBalance1 := app.BankKeeper.GetBalance(ctx, signer1, asset.Symbol)
	require.Equal(t, preBalance1, postBalance1)

	// because rewards is not distributed to LP wallets, the LP balance 2 should not change
	postBalance2 := app.BankKeeper.GetBalance(ctx, signer1, asset.Symbol)
	require.Equal(t, preBalance1, postBalance2)

	// check liquidity provider units 1
	lp1, err = clpKeeper.GetLiquidityProvider(ctx, asset.Symbol, signer1.String())
	require.NoError(t, err)
	require.Equal(t, sdk.NewUint(100), lp1.LiquidityProviderUnits)
	require.Nil(t, lp1.RewardAmount)

	// check liquidity provider units 2
	lp2, err = clpKeeper.GetLiquidityProvider(ctx, asset.Symbol, signer2.String())
	require.NoError(t, err)
	require.Equal(t, sdk.NewUint(128), lp2.LiquidityProviderUnits)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(asset.Symbol, initialAmount)), lp2.RewardAmount)

	// check if pool units changed
	updatedPool, err := clpKeeper.GetPool(ctx, asset.Symbol)
	require.NoError(t, err)
	require.Equal(t, sdk.NewUint(200), updatedPool.NativeAssetBalance)
	require.Equal(t, sdk.NewUint(300), updatedPool.ExternalAssetBalance)
	require.Equal(t, sdk.NewUint(228), updatedPool.PoolUnits)

	// check pool reward amount external
	poolUpdated, err := clpKeeper.GetPool(ctx, asset.Symbol)
	require.NoError(t, err)
	require.Equal(t, sdk.NewUintFromBigInt(initialAmount.BigInt()), poolUpdated.RewardAmountExternal)
}
