package keeper_test

import (
	"fmt"
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_CreatePool(t *testing.T) {
	// nativeAssetAmount sdk.Uint, externalAssetAmount
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	_ = app.ClpKeeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.NoError(t, err)
	if err != nil {
		fmt.Println("Error Generating new pool :", err)
	}
	assert.Equal(t, pool.ExternalAssetBalance, externalAssetAmount)
	assert.Equal(t, pool.NativeAssetBalance, nativeAssetAmount)
}

func TestKeeper_CreateLiquidityProvider(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	asset := types.NewAsset("eth")
	lpAddess, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
	if err != nil {
		fmt.Println("Error Creating Liquidity Provider :", err)
	}
	lp := app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), lpAddess)
	assert.NoError(t, err)
	assert.Equal(t, lp.LiquidityProviderAddress, "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")

}

func TestKeeper_AddLiquidity(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	_, err := sdk.AccAddressFromBech32(test.AddressKey1)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(998)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(998)))
	_ = app.ClpKeeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	msg := types.NewMsgAddLiquidity(signer, asset, nativeAssetAmount, externalAssetAmount)
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(998), sdk.NewUint(998))
	if err != nil {
		fmt.Println("Error Creating Liquidity Provider :", err)
	}
}

func TestKeeper_RemoveLiquidityProvider(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	asset := types.NewAsset("eth")
	initialBalance := sdk.NewUintFromString("1")
	newAssetCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	lpAddess, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
	if err != nil {
		fmt.Println("Error Creating Liquidity Provider :", err)
	}
	lp := app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), lpAddess)
	app.ClpKeeper.RemoveLiquidityProvider(ctx, sdk.Coins{newAssetCoin}.Sort(), lp)
	assert.False(t, app.ClpKeeper.GetLiquidityProviderIterator(ctx).Valid())
}

func TestKeeper_DecommissionPool(t *testing.T) {

	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	_ = app.ClpKeeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	if err != nil {
		fmt.Println("Error Generating new pool :", err)
	}
	app.ClpKeeper.DecommissionPool(ctx, *pool)
	require.NoError(t, err)
}

func TestKeeper_RemoveLiquidity(t *testing.T) {
	// TODO:
}

func TestKeeper_InitiateSwap(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	senderAddress, _ := sdk.AccAddressFromBech32(test.AddressKey1)
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	err := app.ClpKeeper.InitiateSwap(ctx, externalCoin, senderAddress)
	if err != nil {
		fmt.Println("Error doing initialSwap :", err)
	}
	require.NoError(t, err)

}

func TestKeeper_FinalizeSwap(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	assetEth := types.NewAsset("eth")
	assetDash := types.NewAsset("dash")
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	if err != nil {
		fmt.Println("Error Generating new pool :", err)
	}
	// Test Parameters for swap
	// initialBalance: Initial account balance for all assets created.
	initialBalance := sdk.NewUintFromString("1000000000000000000000")
	// poolBalance: Amount funded to pool. The same amount is used both for native and external asset.
	externalCoin1 := sdk.NewCoin(assetEth.Symbol, sdk.Int(initialBalance))
	externalCoin2 := sdk.NewCoin(assetDash.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(initialBalance))
	// Signer is given ETH and RWN (Signer will creat pool and become LP)
	_ = app.ClpKeeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin1, nativeCoin))
	_ = app.ClpKeeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin2))
	msg := types.NewMsgSwap(signer, assetEth, assetDash, sdk.NewUint(1), sdk.NewUint(10))
	app.ClpKeeper.FinalizeSwap(ctx, "1", *pool, msg)
	require.NoError(t, err)
}

func TestKeeper_ParseToInt(t *testing.T) {
	_, app := test.CreateTestAppClp(false)
	res, boolean := app.ClpKeeper.ParseToInt("1")
	assert.True(t, boolean)
	assert.Equal(t, res, sdk.NewUint(1))
}
