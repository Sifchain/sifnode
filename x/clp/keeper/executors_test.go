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
	lpAddess, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
	if err != nil {
		fmt.Println("Error Creating Liquidity Provider :", err)
	}
	lp := app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), lpAddess)
	app.ClpKeeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.False(t, app.ClpKeeper.GetLiquidityProviderIterator(ctx).Valid())
}

func TestKeeper_DecommissionPool(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	err := clpKeeper.SetPool(ctx, &pool)
	assert.NoError(t, err)
	getpool, err := clpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	err = clpKeeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	_, err = clpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.Error(t, err, "Pool should be deleted")
	// This should do nothing.
	err = clpKeeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	assert.Error(t, err)
}

func TestKeeper_RemoveLiquidity(t *testing.T) {

}

func TestKeeper_InitiateSwap(t *testing.T) {
	ctx, keeper, bankKeeper, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	msg := types.NewMsgBurn(1, cosmosReceivers[0], ethereumSender, amount, "stake", amount)
	coins := sdk.NewCoins(sdk.NewCoin("stake", amount), sdk.NewCoin(types.CethSymbol, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	err := keeper.ProcessBurn(ctx, cosmosReceivers[0], &msg)
	require.NoError(t, err)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))
}

func TestKeeper_FinalizeSwap(t *testing.T) {

}

func TestKeeper_ParseToInt(t *testing.T) {

}
