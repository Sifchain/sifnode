package clp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandler(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	handler := NewHandler(keeper)
	res, err := handler(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestCreatePool(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	//Parameters for create pool
	initialBalance := 10000 // Initial account balance for all assets created
	poolBalance := 1000     // Amount funded to pool , This same amount is used both for native and external asset

	asset := NewAsset("ETHEREUM", "ETH", "ceth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(initialBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance)))
	_, _ = keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")

	MinThreshold := keeper.GetParams(ctx).MinCreatePoolThreshold
	// Will fail if we are below minimum
	msgCreatePool := NewMsgCreatePool(signer, asset, MinThreshold-1, 0)
	res, err := handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	// Will fail if we ask for too much.
	msgCreatePool = NewMsgCreatePool(signer, asset, uint(initialBalance+1), uint(initialBalance+1))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	// Ask for the right amount.
	msgCreatePool = NewMsgCreatePool(signer, asset, uint(poolBalance), uint(poolBalance))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Can't create it a second time.
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	externalCoin = sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(initialBalance-poolBalance)))
	nativeCoin = sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance-poolBalance)))
	ok = keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")
}

func TestAddLiqudity(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	//Parameters for add liquidity
	initialBalance := 10000 // Initial account balance for all assets created
	poolBalance := 1000     // Amount funded to pool , This same amount is used both for native and external asset
	addLiquidityAmount := 1000

	asset := NewAsset("ETHEREUM", "ETH", "ceth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(initialBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance)))
	_, _ = keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msg := NewMsgAddLiquidity(signer, asset, uint(addLiquidityAmount), uint(addLiquidityAmount))
	res, err := handleMsgAddLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, asset, uint(poolBalance), uint(poolBalance))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = NewMsgAddLiquidity(signer, asset, uint(addLiquidityAmount), uint(addLiquidityAmount))
	res, err = handleMsgAddLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	// Subtracted twice , during create and add
	externalCoin = sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(initialBalance-addLiquidityAmount-addLiquidityAmount)))
	nativeCoin = sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance-addLiquidityAmount-addLiquidityAmount)))
	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")
}

func TestRemoveLiquidity(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	externalDenom := "ceth"
	nativeDenom := GetSettlementAsset().Ticker
	//Parameters for Remove Liquidity
	initialBalance := 10000 // Initial account balance for all assets created
	poolBalance := 1000     // Amount funded to pool , This same amount is used both for native and external asset
	wBasis := 1000
	asymmetry := 10000

	asset := NewAsset("ETHEREUM", "ETH", externalDenom)
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(initialBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance)))
	_, _ = keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msg := NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err := handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)

	wBasis = 1000
	asymmetry = 10000
	msgCreatePool := NewMsgCreatePool(signer, asset, uint(poolBalance), uint(poolBalance))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	nativeAmountOLD := keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom)
	externalAmountOLD := keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom)
	coins := CalculateWithdraw(t, keeper, ctx, asset, signer.String(), uint(wBasis), asymmetry)
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom).Int64(), nativeAmountOLD.Int64())
	assert.Greater(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom).Int64(), externalAmountOLD.Int64())
	ok := keeper.BankKeeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = 1000
	asymmetry = 10000
	nativeAmountOLD = keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom)
	externalAmountOLD = keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom)
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), uint(wBasis), asymmetry)
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom).Int64(), nativeAmountOLD.Int64())
	assert.Greater(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom).Int64(), externalAmountOLD.Int64())
	ok = keeper.BankKeeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = 1000
	asymmetry = 0
	nativeAmountOLD = keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom)
	externalAmountOLD = keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom)
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), uint(wBasis), asymmetry)
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Greater(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom).Int64(), nativeAmountOLD.Int64())
	assert.Greater(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom).Int64(), externalAmountOLD.Int64())
	ok = keeper.BankKeeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = 1000
	asymmetry = -10000
	nativeAmountOLD = keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom)
	externalAmountOLD = keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom)
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), uint(wBasis), asymmetry)
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Greater(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(nativeDenom).Int64(), nativeAmountOLD.Int64())
	assert.Equal(t, keeper.BankKeeper.GetCoins(ctx, signer).AmountOf(externalDenom).Int64(), externalAmountOLD.Int64())
	ok = keeper.BankKeeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = 10000
	asymmetry = 0
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res, "Cannot withdraw pool is too shallow")

	wBasis = 10000
	asymmetry = 100
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res, "Cannot withdraw pool is too shallow")

	newLP := GenerateAddress2()
	_, _ = keeper.BankKeeper.AddCoins(ctx, newLP, sdk.Coins{externalCoin, nativeCoin})
	msgAdd := NewMsgAddLiquidity(newLP, asset, uint(1000), uint(1000))
	res, err = handleMsgAddLiquidity(ctx, keeper, msgAdd)
	require.NoError(t, err)
	require.NotNil(t, res)

	wBasis = 10000
	asymmetry = 10000
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res, "Can withdraw now as new LP has added liquidity")

}
func TestSwap(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	handler := NewHandler(keeper)
	assetEth := NewAsset("ETHEREUM", "ETH", "ceth")
	assetDash := NewAsset("DASH", "DASH", "cdash")

	// Test Parameters for swap
	initialBalance := 10000 // Initial account balance for all assets created
	poolBalance := 1000     // Amount funded to pool , This same amount is used both for native and external asset
	swapSentAssetETH := 100 // Amount Swapped

	externalCoin1 := sdk.NewCoin(assetEth.Ticker, sdk.NewInt(int64(initialBalance)))
	externalCoin2 := sdk.NewCoin(assetDash.Ticker, sdk.NewInt(int64(initialBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance)))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)
	_, _ = keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin1, nativeCoin})
	_, _ = keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin2})

	msg := NewMsgSwap(signer, assetEth, assetDash, 1)
	res, err := handler(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, assetEth, uint(poolBalance), uint(poolBalance))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	msgCreatePool = NewMsgCreatePool(signer, assetDash, uint(poolBalance), uint(poolBalance))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	receivedAmount := CalculateSwapReceived(t, keeper, ctx, assetEth, assetDash, uint(swapSentAssetETH))

	msg = NewMsgSwap(signer, assetEth, assetDash, uint(swapSentAssetETH))
	res, err = handleMsgSwap(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	CoinsExt1 := sdk.NewCoin(assetEth.Ticker, sdk.NewInt(int64(initialBalance-poolBalance-swapSentAssetETH)))     // Created ETH pool and Send amount for swap
	CoinsNative := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance-poolBalance-poolBalance)))           // Creating two pools
	CoinsExt2 := sdk.NewCoin(assetDash.Ticker, sdk.NewInt(int64(initialBalance-poolBalance+int(receivedAmount)))) // Created one pool and Received swap amount
	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{CoinsExt1, CoinsNative, CoinsExt2})
	assert.True(t, ok, "")

}

func TestDecommisionPool(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	handler := NewHandler(keeper)

	//Parameters for Decommission
	initialBalance := 10000 // Initial account balance for all assets created
	poolBalance := 100      // Amount funded to pool , This same amount is used both for native and external asset

	asset := NewAsset("ETHEREUM", "ETH", "ceth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(initialBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(initialBalance)))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)
	_, _ = keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msgCreatePool := NewMsgCreatePool(signer, asset, uint(poolBalance), uint(poolBalance))
	res, err := handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// SIGNER became new LP
	lpNewBalance := initialBalance - poolBalance
	lpCoinsExt := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(lpNewBalance)))
	lpCoinsNative := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(lpNewBalance)))
	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{lpCoinsExt, lpCoinsNative})
	assert.True(t, ok, "")

	msgrm := NewMsgRemoveLiquidity(signer, asset, 5001, 1)

	res, err = handler(ctx, msgrm)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg := NewMsgDecommissionPool(signer, asset.Ticker)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	msgN := NewMsgAddLiquidity(signer, asset, 1000, 1000)
	res, err = handler(ctx, msgN)
	require.Error(t, err)
	require.Nil(t, res)

	// LP refunded coins when decommison
	lpNewBalance = initialBalance

	lpCoinsExt = sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(lpNewBalance)))
	lpCoinsNative = sdk.NewCoin(NativeTicker, sdk.NewInt(int64(lpNewBalance)))
	ok = keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{lpCoinsExt, lpCoinsNative})
	assert.True(t, ok, "")

}

func CalculateWithdraw(t *testing.T, keeper Keeper, ctx sdk.Context, asset Asset, signer string, wBasisPoints uint, asymmetry int) sdk.Coins {
	pool, err := keeper.GetPool(ctx, asset.Ticker)
	assert.NoError(t, err)
	lp, err := keeper.GetLiquidityProvider(ctx, asset.Ticker, signer)
	assert.NoError(t, err)
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, _, swapAmount := calculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance, pool.ExternalAssetBalance, lp.LiquidityProviderUnits,
		int(wBasisPoints), asymmetry)
	externalAssetCoin := sdk.Coin{}
	nativeAssetCoin := sdk.Coin{}
	if asymmetry > 0 {
		swapResult, _, _, _, err := swapOne(GetSettlementAsset(), swapAmount, asset, pool)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Ticker, sdk.NewIntFromUint64(uint64(withdrawExternalAssetAmount+swapResult)))
		nativeAssetCoin = sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(uint64(withdrawNativeAssetAmount)))
	}
	if asymmetry < 0 {
		swapResult, _, _, _, err := swapOne(asset, swapAmount, GetSettlementAsset(), pool)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Ticker, sdk.NewIntFromUint64(uint64(withdrawExternalAssetAmount)))
		nativeAssetCoin = sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(uint64(withdrawNativeAssetAmount+swapResult)))
	}
	if asymmetry == 0 {
		externalAssetCoin = sdk.NewCoin(asset.Ticker, sdk.NewIntFromUint64(uint64(withdrawExternalAssetAmount)))
		nativeAssetCoin = sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(uint64(withdrawNativeAssetAmount)))
	}

	return sdk.Coins{externalAssetCoin, nativeAssetCoin}

}

func CalculateSwapReceived(t *testing.T, keeper Keeper, ctx sdk.Context, assetSent Asset, assetReceived Asset, swapAmount uint) uint {
	inPool, err := keeper.GetPool(ctx, assetSent.Ticker)
	assert.NoError(t, err)
	outPool, err := keeper.GetPool(ctx, assetReceived.Ticker)
	assert.NoError(t, err)
	emitAmount, _, _, _, err := swapOne(assetSent, swapAmount, GetSettlementAsset(), inPool)
	assert.NoError(t, err)
	emitAmount2, _, _, _, err := swapOne(GetSettlementAsset(), emitAmount, assetReceived, outPool)
	assert.NoError(t, err)
	return emitAmount2
}
