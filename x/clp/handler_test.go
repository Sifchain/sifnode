package clp_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
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

	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	ok := keeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = keeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")

	assert.True(t, ok, "")

	MinThreshold := sdk.NewUint(keeper.GetParams(ctx).MinCreatePoolThreshold)
	// Will fail if we are below minimum
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, MinThreshold.Sub(sdk.NewUint(1)), sdk.ZeroUint())
	res, err := handler(ctx, &msgCreatePool) //clp.handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	// Will fail if we ask for too much.
	msgCreatePool = clptypes.NewMsgCreatePool(signer, asset, initialBalance.Add(sdk.NewUint(1)), initialBalance.Add(sdk.NewUint(1)))
	res, err = handler(ctx, &msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	// Ask for the right amount.
	msgCreatePool = clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err = handler(ctx, &msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Can't create it a second time.
	res, err = handler(ctx, &msgCreatePool) //handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.Error(t, err)
	require.Nil(t, res)

	externalCoin = sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance.Sub(poolBalance)))
	nativeCoin = sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance.Sub(poolBalance)))
	ok = keeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = keeper.HasBalance(ctx, signer, nativeCoin)
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

	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	msg := clptypes.NewMsgAddLiquidity(signer, asset, addLiquidityAmount, addLiquidityAmount)
	res, err := handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err = handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = clptypes.NewMsgAddLiquidity(signer, asset, sdk.ZeroUint(), addLiquidityAmount)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	// Subtracted twice , during create and add
	externalCoin = sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance.Sub(addLiquidityAmount).Sub(addLiquidityAmount)))
	nativeCoin = sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance.Sub(addLiquidityAmount).Sub(sdk.ZeroUint())))

	ok := keeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = keeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")

	signer2 := test.GenerateAddress(test.AddressKey2)
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer2, sdk.NewCoins(externalCoin, nativeCoin))
	msg = clptypes.NewMsgAddLiquidity(signer2, asset, addLiquidityAmount, addLiquidityAmount)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	lpList, _, err := keeper.GetLiquidityProvidersForAssetPaginated(ctx, asset, &query.PageRequest{})
	require.NoError(t, err)
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

	asset := clptypes.NewAsset("cacoin")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(poolBalanceCacoin).Add(sdk.Int(addLiquidityAmountCaCoin)))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(poolBalanceRowan).Add(sdk.Int(addLiquidityAmountRowan)))
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalanceRowan, poolBalanceCacoin)
	res, err := handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg := clptypes.NewMsgAddLiquidity(signer, asset, addLiquidityAmountRowan, addLiquidityAmountCaCoin)
	res, err = handler(ctx, &msg)
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

	asset := clptypes.NewAsset(externalDenom)

	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	_ = keeper.GetBankKeeper().AddCoins(ctx, newLP, sdk.NewCoins(externalCoin, nativeCoin))

	msg := clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err := handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res)

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(10000)
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err = handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	coins := CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, coin := range coins {
		ok := keeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(10000)
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, coin := range coins {
		ok := keeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.ZeroInt()
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, coin := range coins {
		ok := keeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}

	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(-10000)
	coins = CalculateWithdraw(t, keeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, coin := range coins {
		ok := keeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}

	wBasis = sdk.NewInt(10000)
	asymmetry = sdk.ZeroInt()
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res, "Cannot withdraw pool is too shallow")

	wBasis = sdk.NewInt(10000)
	asymmetry = sdk.NewInt(100)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res, "Cannot withdraw pool is too shallow")

	msgAdd := clptypes.NewMsgAddLiquidity(newLP, asset, poolBalance, poolBalance)
	res, err = handler(ctx, &msgAdd)
	require.NoError(t, err)
	require.NotNil(t, res)

	wBasis = sdk.NewInt(10000)
	asymmetry = sdk.NewInt(10000)
	msg = clptypes.NewMsgRemoveLiquidity(newLP, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res, "Can withdraw now as new LP has added liquidity")
}

func TestSwap(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	handler := clp.NewHandler(keeper)
	assetEth := clptypes.NewAsset("eth")
	assetDash := clptypes.NewAsset("dash")

	// Test Parameters for swap

	// initialBalance: Initial account balance for all assets created.
	initialBalance := sdk.NewUintFromString("1000000000000000000000")
	// poolBalance: Amount funded to pool. The same amount is used both for native and external asset.
	poolBalance := sdk.NewUintFromString("1000000000000000000")
	swapSentAssetETH := sdk.NewUintFromString("1000000000000000")

	externalCoin1 := sdk.NewCoin(assetEth.Symbol, sdk.Int(initialBalance))
	externalCoin2 := sdk.NewCoin(assetDash.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	// Signer is given ETH and RWN (Signer will creat pool and become LP)
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin1, nativeCoin))
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin2))

	msg := clptypes.NewMsgSwap(signer, assetEth, assetDash, sdk.NewUint(1), sdk.NewUint(10))
	res, err := handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := clptypes.NewMsgCreatePool(signer, assetEth, poolBalance, poolBalance)
	res, err = handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	msgCreatePool = clptypes.NewMsgCreatePool(signer, assetDash, poolBalance, poolBalance)
	res, err = handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	receivedAmount := CalculateSwapReceived(t, keeper, ctx, assetEth, assetDash, swapSentAssetETH)

	msg = clptypes.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, receivedAmount)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Created ETH pool and Send amount for swap
	CoinsExt1 := sdk.NewCoin(assetEth.Symbol, sdk.Int(initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance))).Sub(sdk.Uint(sdk.Int(swapSentAssetETH)))))
	// Creating two pools
	CoinsNative := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance))).Sub(sdk.Uint(sdk.Int(poolBalance)))))
	// Created one pool and Received swap amount
	CoinsExt2 := sdk.NewCoin(assetDash.Symbol, sdk.Int(initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance))).Add(sdk.Uint(sdk.Int(receivedAmount)))))

	ok := keeper.HasBalance(ctx, signer, CoinsExt1)
	assert.True(t, ok, "")
	ok = keeper.HasBalance(ctx, signer, CoinsNative)
	assert.True(t, ok, "")
	ok = keeper.HasBalance(ctx, signer, CoinsExt2)
	assert.True(t, ok, "")

	msg = clptypes.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, swapSentAssetETH)
	res, err = handler(ctx, &msg)
	require.ErrorIs(t, err, clptypes.ErrReceivedAmountBelowExpected)
	require.Nil(t, res)
}

func TestDecommisionPool(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")

	handler := clp.NewHandler(keeper)

	//Parameters for Decommission
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")

	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)
	_ = keeper.GetBankKeeper().AddCoins(ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))

	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err := handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// SIGNER became new LP
	lpNewBalance := initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance)))
	lpCoinsExt := sdk.NewCoin(asset.Symbol, sdk.Int(lpNewBalance))
	lpCoinsNative := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(lpNewBalance))
	ok := keeper.HasBalance(ctx, signer, lpCoinsExt)
	assert.True(t, ok, "")
	ok = keeper.HasBalance(ctx, signer, lpCoinsNative)
	assert.True(t, ok, "")

	msgrm := clptypes.NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(5001), sdk.NewInt(1))

	res, err = handler(ctx, &msgrm)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg := clptypes.NewMsgDecommissionPool(signer, asset.Symbol)
	_, err = handler(ctx, &msg)
	require.Error(t, err)

	v := test.GenerateWhitelistAddress("")
	keeper.SetClpWhiteList(ctx, []sdk.AccAddress{v})

	msg = clptypes.NewMsgDecommissionPool(signer, asset.Symbol)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	msgN := clptypes.NewMsgAddLiquidity(signer, asset, sdk.NewUint(1000), sdk.NewUint(1000))
	res, err = handler(ctx, &msgN)
	require.Error(t, err)
	require.Nil(t, res)

	// LP refunded coins when decommison
	lpNewBalance = initialBalance

	lpCoinsExt = sdk.NewCoin(asset.Symbol, sdk.Int(lpNewBalance))
	lpCoinsNative = sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(lpNewBalance))
	ok = keeper.HasBalance(ctx, signer, lpCoinsExt)
	assert.True(t, ok, "")
	ok = keeper.HasBalance(ctx, signer, lpCoinsNative)
	assert.True(t, ok, "")
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
	if asymmetry.IsPositive() {
		swapResult, _, _, _, err := clpkeeper.SwapOne(clptypes.GetSettlementAsset(), swapAmount, asset, pool)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount.Add(swapResult)))
		nativeAssetCoin = sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount))
	}
	if asymmetry.IsNegative() {
		swapResult, _, _, _, err := clpkeeper.SwapOne(asset, swapAmount, clptypes.GetSettlementAsset(), pool)
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

func CalculateSwapReceived(t *testing.T, keeper clpkeeper.Keeper, ctx sdk.Context, assetSent clptypes.Asset, assetReceived clptypes.Asset, swapAmount sdk.Uint) sdk.Uint {
	inPool, err := keeper.GetPool(ctx, assetSent.Symbol)
	assert.NoError(t, err)
	outPool, err := keeper.GetPool(ctx, assetReceived.Symbol)
	assert.NoError(t, err)
	emitAmount, _, _, _, err := clpkeeper.SwapOne(assetSent, swapAmount, clptypes.GetSettlementAsset(), inPool)
	assert.NoError(t, err)
	emitAmount2, _, _, _, err := clpkeeper.SwapOne(clptypes.GetSettlementAsset(), emitAmount, assetReceived, outPool)
	assert.NoError(t, err)
	return emitAmount2
}
