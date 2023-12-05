package keeper_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_CreatePool_Error(t *testing.T) {
	// nativeAssetAmount sdk.Uint, externalAssetAmount
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth0123456789012345678901234567890123456789012345678901234567890123456789")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUintFromString("10000")))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUintFromString("10000")))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Test Create Pool with invalid pool asset name
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.Error(t, err, "Invalid pool asset name.")
	assert.Nil(t, pool)
}

func TestKeeper_CreatePool_Range(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	nativeAssetAmount2 := sdk.NewUintFromString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	externalAssetAmount2 := sdk.NewUintFromString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUintFromString("0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUintFromString("0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount2, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.Error(t, err, "Unable to parse to Int")
	assert.Nil(t, pool)
	msgCreatePool = types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount2)
	// Create Pool
	pool, err = app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.Error(t, err, "Unable to parse to Int")
	assert.Nil(t, pool)

}

func TestKeeper_CreatePool_And_AddLiquidity_RemoveLiquidity(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	signer2 := test.GenerateAddress("")
	asset2 := types.NewAsset("xxx")
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	nativeAssetAmount2 := sdk.NewUintFromString("0xffff")
	externalAssetAmount2 := sdk.NewUintFromString("0xffff")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))

	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	msgCreatePool := types.NewMsgCreatePool(nil, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool with empty address
	_, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.Error(t, err, "empty address string is not allowed")

	msgCreatePool = types.NewMsgCreatePool(signer, asset, nativeAssetAmount2, externalAssetAmount2)
	// Create Pool with user does not have enough balance
	_, err = app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.Error(t, err, "user does not have enough balance of the required coin")

	msgCreatePool = types.NewMsgCreatePool(signer2, asset2, nativeAssetAmount, externalAssetAmount)
	// Create Pool with user does not have enough balance for singer2
	_, err = app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.Error(t, err, "user does not have enough balance of the required coin")

	msgCreatePool = types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	require.NoError(t, err, "Error Generating new pool", err)
	_, err = app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), nil)
	assert.Error(t, err, "MsgCreatePool can not be nil")
	msg := types.NewMsgAddLiquidity(signer, asset, nativeAssetAmount, externalAssetAmount)
	app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), signer)
	lp, err := app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.Equal(t, lp.LiquidityProviderAddress, "sif15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqhns3lt")
	assert.NoError(t, err)
	assert.Equal(t, pool.ExternalAssetBalance, externalAssetAmount)
	assert.Equal(t, pool.NativeAssetBalance, nativeAssetAmount)
	msg = types.NewMsgAddLiquidity(signer2, asset2, nativeAssetAmount, externalAssetAmount)
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.Error(t, err, "insufficient funds")
	msg = types.NewMsgAddLiquidity(nil, asset, nativeAssetAmount, externalAssetAmount)
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.Error(t, err, "empty address string is not allowed")
	msg = types.NewMsgAddLiquidity(signer2, asset2, sdk.NewUintFromString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), sdk.NewUintFromString("998"))
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.Error(t, err, "Unable to parse to Int")
	msg = types.NewMsgAddLiquidity(signer2, asset2, sdk.NewUintFromString("998"), sdk.NewUintFromString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.Error(t, err, "Unable to parse to Int")
	msg = types.NewMsgAddLiquidity(signer2, asset2, sdk.NewUintFromString("99800000001"), sdk.NewUintFromString("99800000001"))
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.Error(t, err, "Not enough money in account")
	subCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(100)))
	errorRemoveLiquidity := app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(989), sdk.NewUint(10001), sdk.NewUint(10001))
	assert.NoError(t, errorRemoveLiquidity)
	ok := app.ClpKeeper.HasBalance(ctx, signer, subCoin)
	assert.True(t, ok, "")

	subCoin = sdk.NewCoin(asset2.Symbol, sdk.Int(sdk.NewUint(100)))
	pool.GetExternalAsset().Symbol = ""
	errorRemoveLiquidity = app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(989), sdk.NewUint(1000), sdk.NewUint(1000))
	assert.Error(t, errorRemoveLiquidity, "Unable to set pool")

	pool.GetExternalAsset().Symbol = "eth"

	subCoin = sdk.NewCoin(asset2.Symbol, sdk.Int(sdk.NewUint(100)))
	errorRemoveLiquidity = app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(989), sdk.NewUint(1000), sdk.NewUint(1000))
	assert.Error(t, errorRemoveLiquidity, "pool does not have sufficient balance")

	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	subCoin = sdk.NewCoin(asset2.Symbol, sdk.Int(sdk.NewUint(100)))
	errorRemoveLiquidity = app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(989), sdk.NewUint(1000), sdk.NewUint(1000))
	assert.Error(t, errorRemoveLiquidity, "pool does not have sufficient balance")

	subCoin = sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(100)))
	errorRemoveLiquidity = app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(989), sdk.NewUint(99), sdk.NewUint(99))
	assert.Error(t, errorRemoveLiquidity, "Cannot withdraw pool is too shallow")

	subCoin = sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(100)))
	errorRemoveLiquidity = app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(989), sdk.NewUint(10001), sdk.NewUint(10001))
	assert.NoError(t, errorRemoveLiquidity)
	res := app.ClpKeeper.HasBalance(ctx, signer, subCoin)
	assert.True(t, res, "Cannot withdraw pool is too shallow")
	subCoin = sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(100)))
	errorRemoveLiquidity = app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(0), sdk.NewUint(10001), sdk.NewUint(10001))
	assert.NoError(t, errorRemoveLiquidity)
	lp.LiquidityProviderAddress = ""
	_ = app.ClpKeeper.RemoveLiquidity(ctx, *pool, subCoin, subCoin, *lp, sdk.NewUint(0), sdk.NewUint(10001), sdk.NewUint(10001))
	assert.Error(t, err, "empty address string is not allowed")
}

func TestKeeper_CreateLiquidityProvider(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	asset := types.NewAsset("eth")
	lpAddress, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
	require.NoError(t, err, "Error Creating Liquidity Provider :", err)
	lp := app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), lpAddress)
	assert.NoError(t, err)
	assert.Equal(t, lp.LiquidityProviderAddress, "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
}

func TestKeeper_RemoveLiquidityProvider(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	asset2 := types.NewAsset("xxx")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	subCoin := sdk.NewUintFromString("1")
	newAssetCoin := sdk.NewCoin(asset.Symbol, sdk.Int(subCoin))
	newAssetCoin2 := sdk.NewCoin(asset2.Symbol, sdk.Int(subCoin))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, _ := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	msg := types.NewMsgAddLiquidity(signer, asset, nativeAssetAmount, externalAssetAmount)
	app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), signer)
	lp, _ := app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	getlp, _ := app.ClpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	lp1, _ := app.ClpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.True(t, app.ClpKeeper.GetLiquidityProviderIterator(ctx).Valid())
	app.ClpKeeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	_, err := app.ClpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.Error(t, err, "LiquidityProvider has been deleted")
	lp1.LiquidityProviderAddress = ""
	err = app.ClpKeeper.RemoveLiquidityProvider(ctx, sdk.Coins{newAssetCoin2}.Sort(), lp1)
	assert.Error(t, err, "empty address string is not allowed")
	err = app.ClpKeeper.RemoveLiquidityProvider(ctx, sdk.Coins{newAssetCoin2}.Sort(), getlp)
	assert.Error(t, err, "unable to add balance")
	err = app.ClpKeeper.RemoveLiquidityProvider(ctx, sdk.Coins{newAssetCoin}.Sort(), getlp)
	assert.NoError(t, err)
	assert.False(t, app.ClpKeeper.GetLiquidityProviderIterator(ctx).Valid())
	msg = types.NewMsgAddLiquidity(signer, asset, nativeAssetAmount, externalAssetAmount)
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.NoError(t, err)

	msg = types.NewMsgAddLiquidity(signer, asset, nativeAssetAmount, externalAssetAmount)
	pool.GetExternalAsset().Symbol = ""
	_, err = app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.Error(t, err, "Unable to set pool")

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
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	require.NoError(t, err, "Error Generating new pool :", err)

	err = app.ClpKeeper.DecommissionPool(ctx, *pool)
	require.NoError(t, err)
	_, err = app.ClpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.Error(t, err, "Pool should be deleted")
	err = app.ClpKeeper.DecommissionPool(ctx, *pool)
	assert.Error(t, err, "Unable to destroy pool")
}

func TestKeeper_InitiateSwap(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	asset := types.NewAsset("eth")
	asset1 := types.NewAsset("xxx")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	externalCoin1 := sdk.NewCoin(asset1.Symbol, sdk.Int(sdk.NewUint(10000000000000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	err := app.ClpKeeper.InitiateSwap(ctx, externalCoin, signer)
	require.NoError(t, err)
	ok := app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")
	err = app.ClpKeeper.InitiateSwap(ctx, externalCoin1, signer)
	assert.Error(t, err, "user does not have enough balance of the required coin")

}

func TestKeeper_FinalizeSwap(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	assetEth := types.NewAsset("eth")
	assetDash := types.NewAsset("dash")
	externalCoin := sdk.NewCoin(assetEth.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	msgCreatePool := types.NewMsgCreatePool(signer, assetEth, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	_, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	require.NoError(t, err, "Error Generating new pool :", err)
	externalCoin = sdk.NewCoin(assetDash.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin = sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	msgCreatePool = types.NewMsgCreatePool(signer, assetDash, nativeAssetAmount, externalAssetAmount)
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	require.NoError(t, err, "Error Generating new pool :", err)
	// Test Parameters for swap
	// initialBalance: Initial account balance for all assets created.
	initialBalance := sdk.NewUintFromString("1000000000000000000000")
	// poolBalance: Amount funded to pool. The same amount is used both for native and external asset.
	externalCoin1 := sdk.NewCoin("eth", sdk.Int(initialBalance))
	externalCoin2 := sdk.NewCoin("dash", sdk.Int(initialBalance))
	// Signer is given ETH and RWN (Signer will creat pool and become LP)
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin1, nativeCoin))
	_ = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin2, nativeCoin))
	msg := types.NewMsgSwap(signer, assetEth, assetDash, sdk.NewUint(1), sdk.NewUint(10))
	err = app.ClpKeeper.FinalizeSwap(ctx, "", *pool, msg)
	assert.Error(t, err, "Unable to parse to Int")

	err = app.ClpKeeper.FinalizeSwap(ctx, "1", *pool, msg)
	require.NoError(t, err)

	msg = types.NewMsgSwap(signer, types.NewAsset("xxx"), types.NewAsset("xxxx"), sdk.NewUint(1), sdk.NewUint(10))
	err = app.ClpKeeper.FinalizeSwap(ctx, "1", *pool, msg)
	assert.Error(t, err, "insufficient funds")

	msg = types.NewMsgSwap(nil, assetEth, assetDash, sdk.NewUint(1), sdk.NewUint(10))
	err = app.ClpKeeper.FinalizeSwap(ctx, "1", *pool, msg)
	assert.Error(t, err, "empty address string is not allowed")
	msg = types.NewMsgSwap(signer, assetEth, assetDash, sdk.NewUint(1), sdk.NewUint(10))
	pool.ExternalAsset.Symbol = ""
	err = app.ClpKeeper.FinalizeSwap(ctx, "1", *pool, msg)
	assert.Error(t, err, "Unable to set pool")
}

func TestKeeper_ParseToInt(t *testing.T) {
	_, app := test.CreateTestAppClp(false)
	res, boolean := app.ClpKeeper.ParseToInt("1")
	assert.True(t, boolean)
	assert.Equal(t, res.String(), "1")
}

func TestKeeper_ParseToInt_WithBigNumber(t *testing.T) {
	_, app := test.CreateTestAppClp(false)
	_, boolean := app.ClpKeeper.ParseToInt("10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.False(t, boolean)
}
