package clp_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
)

func TestCreatePool(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	asset := clptypes.NewAsset("eth")
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	handler := clp.NewHandler(app.ClpKeeper)
	signer := test.GenerateAddress("")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	ok := app.ClpKeeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, initialBalance, poolBalance)
	res, err := handler(ctx, &msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// check for failure if we try to create a pool twice
	msgCreatePool = clptypes.NewMsgCreatePool(signer, asset, initialBalance, poolBalance)
	_, err = handler(ctx, &msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err, clptypes.ErrPoolTooShallow)
}

func TestGetPool(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	asset := clptypes.NewAsset("eth")
	_, err := app.ClpKeeper.GetPool(ctx, asset.Symbol)
	require.Error(t, err, clptypes.ErrPoolDoesNotExist)

	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	handler := clp.NewHandler(app.ClpKeeper)
	signer := test.GenerateAddress("")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	ok := app.ClpKeeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, initialBalance, poolBalance)
	res, err := handler(ctx, &msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	_, err = app.ClpKeeper.GetPool(ctx, asset.Symbol)
	assert.NoError(t, err)

}

func TestAddLiquidityErrorCases(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
	//Parameters for add liquidity
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	addLiquidityAmount := sdk.NewUintFromString("1000000000000000000")
	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err := handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg := clptypes.NewMsgAddLiquidity(signer, asset, sdk.ZeroUint(), addLiquidityAmount)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	asset1 := clptypes.NewAsset("btc")
	msg1 := clptypes.NewMsgAddLiquidity(signer, asset1, sdk.ZeroUint(), addLiquidityAmount)
	_, err = handler(ctx, &msg1)
	require.Error(t, err, clptypes.ErrTokenNotSupported)
	asset1 = clptypes.NewAsset("eth")
	msg1 = clptypes.NewMsgAddLiquidity(signer, asset1, sdk.ZeroUint(), addLiquidityAmount)
	_, err = handler(ctx, &msg1)
	require.NoError(t, err)

}

func TestPoolMultiplyCases(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	newLP := test.GenerateAddress(test.AddressKey2)
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
	externalDenom := "eth"
	assetDash := clptypes.NewAsset("dash")
	initialBalance := sdk.NewUintFromString("9999999999999")      // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("100000000000000000000") // Amount funded to pool , This same amount is used both for native and external asset
	wBasis := sdk.NewInt(1000)
	asymmetry := sdk.NewInt(10000)
	asset := clptypes.NewAsset(externalDenom)
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, newLP, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	// Fail if amount is greater than user has
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	_, err = handler(ctx, &msgCreatePool)
	require.Error(t, err, clptypes.ErrBalanceNotAvailable)

	// Fail if amount is less than or equal to minimum
	poolBalance = sdk.NewUintFromString("100000") // Amount funded to pool , This same amount is used both for native and external asset
	msgCreatePool = clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	_, err = handler(ctx, &msgCreatePool)
	require.Error(t, err, clptypes.ErrTotalAmountTooLow)
	// Only works the first time, fails later
	initialBalance = sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance = sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	addLiquidityAmount := sdk.NewUintFromString("1000000000000000000")
	externalCoin = sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin = sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	externalCoin = sdk.NewCoin(assetDash.Symbol, sdk.Int(initialBalance))
	nativeCoin = sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, newLP, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool = clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	_, err = handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	// check for failure if we try to create a pool twice
	msgCreatePool = clptypes.NewMsgCreatePool(signer, asset, initialBalance, poolBalance)
	_, err = handler(ctx, &msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err, clptypes.ErrUnableToCreatePool)
	// ensure we can add liquidity, money gets transferred
	msg := clptypes.NewMsgAddLiquidity(signer, asset, sdk.ZeroUint(), addLiquidityAmount)
	res, err := handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	// ensure we can remove liquidity, money gets transferred
	coins := CalculateWithdraw(t, clpKeeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	reMsg := clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &reMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, coin := range coins {
		ok := clpKeeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}
	// check for failure if we try to remove more
	wBasis = sdk.NewInt(10000)
	asymmetry = sdk.ZeroInt()
	reMsg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	_, err = handler(ctx, &reMsg)
	require.Error(t, err, clptypes.ErrPoolTooShallow)
	// check for failure if we try to add too much liquidity: TestAddLiquidity_LargeValue
	// check for failure if we try to swap too much for user
	swapSentAssetETH := sdk.NewUintFromString("1000000000000000000000000000")
	assetEth := clptypes.NewAsset("eth")
	swMsg := clptypes.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, sdk.NewUintFromString("10000000000000"))
	_, err = handler(ctx, &swMsg)
	require.Error(t, err, clptypes.ErrPoolDoesNotExist)

	poolBalance = sdk.NewUintFromString("1000000000000000000")
	msgCreatePool = clptypes.NewMsgCreatePool(signer, assetDash, poolBalance, poolBalance)
	_, err = handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	swMsg = clptypes.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, sdk.NewUintFromString("10000000000000"))
	_, err = handler(ctx, &swMsg)
	require.Error(t, err, clptypes.ErrBalanceNotAvailable)
	// check for failure if we try to swap and receive amount is below expected
	swapSentAssetETH = sdk.NewUintFromString("99999999")
	swMsg = clptypes.NewMsgSwap(signer, assetDash, assetEth, swapSentAssetETH, sdk.NewUintFromString("10000000000000"))
	_, err = handler(ctx, &swMsg)
	require.Error(t, err, clptypes.ErrReceivedAmountBelowExpected)
	// now try to do a swap that works
	swapSentAssetETH = sdk.NewUintFromString("10000000000009000009")
	swMsg = clptypes.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, sdk.NewUintFromString("100000000009"))
	_, err = handler(ctx, &swMsg)
	require.NoError(t, err)
}

func CalculateWithdraw(t *testing.T, keeper clpkeeper.Keeper, ctx sdk.Context, asset clptypes.Asset, signer string, wBasisPoints string, asymmetry sdk.Int) sdk.Coins {
	pool, err := keeper.GetPool(ctx, asset.Symbol)
	assert.NoError(t, err)
	lp, err := keeper.GetLiquidityProvider(ctx, asset.Symbol, signer)
	assert.NoError(t, err)
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, _, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		wBasisPoints, asymmetry)
	externalAssetCoin := sdk.Coin{}
	nativeAssetCoin := sdk.Coin{}
	ctx, app := test.CreateTestAppClp(false)
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := app.TokenRegistryKeeper.GetEntry(registry, pool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	if asymmetry.IsPositive() {
		normalizationFactor, adjustExternalToken := keeper.GetNormalizationFactor(eAsset.Decimals)
		swapResult, _, _, _, err := clpkeeper.SwapOne(clptypes.GetSettlementAsset(), swapAmount, asset, pool, normalizationFactor, adjustExternalToken)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount.Add(swapResult)))
		nativeAssetCoin = sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount))
	}
	if asymmetry.IsNegative() {
		normalizationFactor, adjustExternalToken := keeper.GetNormalizationFactor(eAsset.Decimals)
		swapResult, _, _, _, err := clpkeeper.SwapOne(asset, swapAmount, clptypes.GetSettlementAsset(), pool, normalizationFactor, adjustExternalToken)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount))
		nativeAssetCoin = sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount.Add(swapResult)))
	}
	if asymmetry.IsZero() {
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount))
		nativeAssetCoin = sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount))
	}
	return sdk.NewCoins(externalAssetCoin, nativeAssetCoin)
}
