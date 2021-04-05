package clp_test

import (
	"github.com/Sifchain/sifnode/x/clp"
	keeper2 "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandler(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	handler := clp.NewHandler(keeper)
	res, err := handler(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res)

}

func TestCreatePool(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	handler := clp.NewHandler(keeper)
	signer := test.GenerateAddress("")
	//Parameters for create pool
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset

	asset := clp.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance))
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	ok := keeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")

	MinThreshold := sdk.NewUint(uint64(keeper.GetParams(ctx).MinCreatePoolThreshold))
	// Will fail if we are below minimum
	msgCreatePool := clp.NewMsgCreatePool(signer, asset, MinThreshold.Sub(sdk.NewUint(1)), sdk.ZeroUint())
	res, err := handler(ctx, msgCreatePool) //clp.handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	// Will fail if we ask for too much.
	msgCreatePool = clp.NewMsgCreatePool(signer, asset, initialBalance.Add(sdk.NewUint(1)), initialBalance.Add(sdk.NewUint(1)))
	res, err = handler(ctx, msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	// Ask for the right amount.
	msgCreatePool = clp.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err = handler(ctx, msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Can't create it a second time.
	res, err = handler(ctx, msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	externalCoin = sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance.Sub(poolBalance)))
	nativeCoin = sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance.Sub(poolBalance)))
	ok = keeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")
}

func TestAddLiquidity(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	handler := clp.NewHandler(keeper)
	//Parameters for add liquidity
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	addLiquidityAmount := sdk.NewUintFromString("1000000000000000000")

	asset := clp.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance))
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msg := clp.NewMsgAddLiquidity(signer, asset, addLiquidityAmount, addLiquidityAmount)
	res, err := handler(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := clp.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err = handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = clp.NewMsgAddLiquidity(signer, asset, sdk.ZeroUint(), addLiquidityAmount)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	// Subtracted twice , during create and add
	externalCoin = sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance.Sub(addLiquidityAmount).Sub(addLiquidityAmount)))
	nativeCoin = sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance.Sub(addLiquidityAmount).Sub(sdk.ZeroUint())))
	ok := keeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")

	signer2 := test.GenerateAddress(test.AddressKey2)
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer2, sdk.Coins{externalCoin, nativeCoin})
	msg = clp.NewMsgAddLiquidity(signer2, asset, addLiquidityAmount, addLiquidityAmount)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	lpList := keeper.GetLiquidityProvidersForAsset(ctx, asset)
	assert.Equal(t, 2, len(lpList))

}

func TestAddLiquidity_LargeValue(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	handler := clp.NewHandler(keeper)
	//Parameters for add liquidity
	poolBalanceRowan := sdk.NewUintFromString("162057826929020210025062784")
	poolBalanceCacoin := sdk.NewUintFromString("1000000000000000000000") // Amount funded to pool , This same amount is used both for native and external asset
	addLiquidityAmountRowan := sdk.NewUintFromString("1000000000000000000000")
	addLiquidityAmountCaCoin := sdk.NewUintFromString("8999998679900000000000000000000")

	asset := clp.NewAsset("cacoin")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(poolBalanceCacoin).Add(sdk.Int(addLiquidityAmountCaCoin)))
	nativeCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(poolBalanceRowan).Add(sdk.Int(addLiquidityAmountRowan)))
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msgCreatePool := clp.NewMsgCreatePool(signer, asset, poolBalanceRowan, poolBalanceCacoin)
	res, err := handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg := clp.NewMsgAddLiquidity(signer, asset, addLiquidityAmountRowan, addLiquidityAmountCaCoin)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

}

func TestRemoveLiquidity(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	newLP := test.GenerateAddress(test.AddressKey2)
	handler := clp.NewHandler(keeper)
	externalDenom := "eth"
	initialBalance := sdk.NewUintFromString("100000000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("10000000000000000000")        // Amount funded to pool , This same amount is used both for native and external asset
	wBasis := sdk.NewInt(1000)
	asymmetry := sdk.NewInt(10000)

	asset := clp.NewAsset(externalDenom)

	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance))
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, newLP, sdk.Coins{externalCoin, nativeCoin})

	msg := clp.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err := handler(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res)

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(10000)
	msgCreatePool := clp.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err = handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	coins := CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clp.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	ok := keeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(10000)
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clp.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	ok = keeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.ZeroInt()
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clp.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	ok = keeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(-10000)
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clp.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	ok = keeper.HasCoins(ctx, signer, coins)
	assert.True(t, ok, "")

	wBasis = sdk.NewInt(10000)
	asymmetry = sdk.ZeroInt()
	msg = clp.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res, "Cannot withdraw pool is too shallow")

	wBasis = sdk.NewInt(10000)
	asymmetry = sdk.NewInt(100)
	msg = clp.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res, "Cannot withdraw pool is too shallow")

	msgAdd := clp.NewMsgAddLiquidity(newLP, asset, poolBalance, poolBalance)
	res, err = handler(ctx, msgAdd)
	require.NoError(t, err)
	require.NotNil(t, res)

	wBasis = sdk.NewInt(10000)
	asymmetry = sdk.NewInt(10000)
	msg = clp.NewMsgRemoveLiquidity(newLP, asset, wBasis, asymmetry)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res, "Can withdraw now as new LP has added liquidity")

}
func TestSwap(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	handler := clp.NewHandler(keeper)
	assetEth := clp.NewAsset("eth")
	assetDash := clp.NewAsset("dash")

	// Test Parameters for swap

	initialBalance := sdk.NewUintFromString("1000000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")       // Amount funded to pool , This same amount is used both for native and external asset
	swapSentAssetETH := sdk.NewUintFromString("1000000000000000")

	externalCoin1 := sdk.NewCoin(assetEth.Symbol, sdk.Int(initialBalance))
	externalCoin2 := sdk.NewCoin(assetDash.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.Coins{externalCoin1, nativeCoin})
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.Coins{externalCoin2})

	msg := clp.NewMsgSwap(signer, assetEth, assetDash, sdk.NewUint(1), sdk.NewUint(10))
	res, err := handler(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := clp.NewMsgCreatePool(signer, assetEth, poolBalance, poolBalance)
	res, err = handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	msgCreatePool = clp.NewMsgCreatePool(signer, assetDash, poolBalance, poolBalance)
	res, err = handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	receivedAmount := CalculateSwapReceived(t, keeper, ctx, assetEth, assetDash, swapSentAssetETH)

	msg = clp.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, receivedAmount)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	CoinsExt1 := sdk.NewCoin(assetEth.Symbol, sdk.Int(initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance))).Sub(sdk.Uint(sdk.Int(swapSentAssetETH))))) // Created ETH pool and Send amount for swap
	CoinsNative := sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance))).Sub(sdk.Uint(sdk.Int(poolBalance)))))   // Creating two pools
	CoinsExt2 := sdk.NewCoin(assetDash.Symbol, sdk.Int(initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance))).Add(sdk.Uint(sdk.Int(receivedAmount)))))  // Created one pool and Received swap amount
	ok := keeper.HasCoins(ctx, signer, sdk.Coins{CoinsExt1, CoinsNative, CoinsExt2})
	assert.True(t, ok, "")

	msg = clp.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, swapSentAssetETH)
	res, err = handler(ctx, msg)
	require.Error(t, err)
	require.NotNil(t, res)

}

func TestDecommisionPool(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")

	handler := clp.NewHandler(keeper)

	//Parameters for Decommission
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")

	asset := clp.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)
	_, _ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msgCreatePool := clp.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err := handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// SIGNER became new LP
	lpNewBalance := initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance)))
	lpCoinsExt := sdk.NewCoin(asset.Symbol, sdk.Int(lpNewBalance))
	lpCoinsNative := sdk.NewCoin(clp.NativeSymbol, sdk.Int(lpNewBalance))
	ok := keeper.HasCoins(ctx, signer, sdk.Coins{lpCoinsExt, lpCoinsNative})
	assert.True(t, ok, "")

	msgrm := clp.NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(5001), sdk.NewInt(1))

	res, err = handler(ctx, msgrm)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg := clp.NewMsgDecommissionPool(signer, asset.Symbol)
	_, err = handler(ctx, msg)
	require.Error(t, err)

	v := test.GenerateWhitelistAddress("")
	keeper.SetClpWhiteList(ctx, []sdk.AccAddress{v})

	msg = clp.NewMsgDecommissionPool(signer, asset.Symbol)
	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	msgN := clp.NewMsgAddLiquidity(signer, asset, sdk.NewUint(1000), sdk.NewUint(1000))
	res, err = handler(ctx, msgN)
	require.Error(t, err)
	require.Nil(t, res)

	// LP refunded coins when decommison
	lpNewBalance = initialBalance

	lpCoinsExt = sdk.NewCoin(asset.Symbol, sdk.Int(lpNewBalance))
	lpCoinsNative = sdk.NewCoin(clp.NativeSymbol, sdk.Int(lpNewBalance))
	ok = keeper.HasCoins(ctx, signer, sdk.Coins{lpCoinsExt, lpCoinsNative})
	assert.True(t, ok, "")

}

func CalculateWithdraw(t *testing.T, keeper clp.Keeper, ctx sdk.Context, asset clp.Asset, signer string, wBasisPoints string, asymmetry sdk.Int) sdk.Coins {
	pool, err := keeper.GetPool(ctx, asset.Symbol)
	assert.NoError(t, err)
	lp, err := keeper.GetLiquidityProvider(ctx, asset.Symbol, signer)
	assert.NoError(t, err)
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, _, swapAmount := keeper2.CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		wBasisPoints, asymmetry)
	externalAssetCoin := sdk.Coin{}
	nativeAssetCoin := sdk.Coin{}
	if asymmetry.IsPositive() {
		swapResult, _, _, _, err := keeper2.SwapOne(clp.GetSettlementAsset(), swapAmount, asset, pool)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount.Add(swapResult)))
		nativeAssetCoin = sdk.NewCoin(clp.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount))
	}
	if asymmetry.IsNegative() {
		swapResult, _, _, _, err := keeper2.SwapOne(asset, swapAmount, clp.GetSettlementAsset(), pool)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount))
		nativeAssetCoin = sdk.NewCoin(clp.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount.Add(swapResult)))
	}
	if asymmetry.IsZero() {
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount))
		nativeAssetCoin = sdk.NewCoin(clp.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount))
	}

	return sdk.Coins{externalAssetCoin, nativeAssetCoin}

}

func CalculateSwapReceived(t *testing.T, keeper clp.Keeper, ctx sdk.Context, assetSent clp.Asset, assetReceived clp.Asset, swapAmount sdk.Uint) sdk.Uint {
	inPool, err := keeper.GetPool(ctx, assetSent.Symbol)
	assert.NoError(t, err)
	outPool, err := keeper.GetPool(ctx, assetReceived.Symbol)
	assert.NoError(t, err)
	emitAmount, _, _, _, err := keeper2.SwapOne(assetSent, swapAmount, clp.GetSettlementAsset(), inPool)
	assert.NoError(t, err)
	emitAmount2, _, _, _, err := keeper2.SwapOne(clp.GetSettlementAsset(), emitAmount, assetReceived, outPool)
	assert.NoError(t, err)
	return emitAmount2
}
