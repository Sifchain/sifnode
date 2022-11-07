package clp_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"

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
	ctx, app := test.CreateTestAppClp(false)
	handler := clp.NewHandler(app.ClpKeeper)
	res, err := handler(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestCreatePool(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	handler := clp.NewHandler(app.ClpKeeper)
	signer := test.GenerateAddress("")
	//Parameters for create pool
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	ok := app.ClpKeeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")
	assert.True(t, ok, "")
	MinThreshold := sdk.NewUint(app.ClpKeeper.GetParams(ctx).MinCreatePoolThreshold)
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
	ok = app.ClpKeeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")
	newAsset := clptypes.NewAsset("Asset")
	// Not whitelisted
	msgNonWhitelisted := clptypes.NewMsgCreatePool(signer, clptypes.NewAsset(newAsset.Symbol), poolBalance, poolBalance)
	_, err = handler(ctx, &msgNonWhitelisted)
	require.Error(t, err)
	// Whitelist Asset
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{Denom: newAsset.Symbol, Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}})
	newAssetCoin := sdk.NewCoin(newAsset.Symbol, sdk.Int(initialBalance))
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.Coins{newAssetCoin}.Sort())
	require.NoError(t, err)
	// Create Pool
	res, err = handler(ctx, &msgNonWhitelisted)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestAddLiquidity(t *testing.T) {
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
	ok := clpKeeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = clpKeeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")
	signer2 := test.GenerateAddress(test.AddressKey2)
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer2, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msg = clptypes.NewMsgAddLiquidity(signer2, asset, addLiquidityAmount, addLiquidityAmount)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	lpList, _, err := clpKeeper.GetLiquidityProvidersForAssetPaginated(ctx, asset, &query.PageRequest{})
	require.NoError(t, err)
	assert.Equal(t, 2, len(lpList))
	newAsset := clptypes.NewAsset("Asset")
	msgNonWhitelisted := clptypes.NewMsgAddLiquidity(signer, newAsset, sdk.NewUint(1000), sdk.NewUint(1000))
	_, err = handler(ctx, &msgNonWhitelisted)
	require.Error(t, err)
}

func TestAddLiquidity_LargeValue(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
	//Parameters for add liquidity
	poolBalanceRowan := sdk.NewUintFromString("162057826929020210025062784")
	poolBalanceCacoin := sdk.NewUintFromString("1000000000000000000000") // Amount funded to pool , This same amount is used both for native and external asset
	addLiquidityAmountRowan := sdk.NewUintFromString("1000000000000000000000")
	addLiquidityAmountCaCoin := sdk.NewUintFromString("8999998679900000000000000000000")
	asset := clptypes.NewAsset("cacoin")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(poolBalanceCacoin).Add(sdk.Int(addLiquidityAmountCaCoin)))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(poolBalanceRowan).Add(sdk.Int(addLiquidityAmountRowan)))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
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
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	newLP := test.GenerateAddress(test.AddressKey2)
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
	externalDenom := "eth"
	initialBalance := sdk.NewUintFromString("100000000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("10000000000000000000")        // Amount funded to pool , This same amount is used both for native and external asset
	wBasis := sdk.NewInt(1000)
	asymmetry := sdk.NewInt(10000)
	asset := clptypes.NewAsset(externalDenom)
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, newLP, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)

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
	UnlockAllliquidity(app, ctx, asset, signer, t)

	coins := CalculateWithdraw(t, clpKeeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.EqualError(t, err, "Cannot remove liquidity asymmetrically")
	require.Nil(t, res)
	for _, coin := range coins {
		ok := clpKeeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}
	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(10000)
	coins = CalculateWithdraw(t, clpKeeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.EqualError(t, err, "Cannot remove liquidity asymmetrically")
	require.Nil(t, res)
	for _, coin := range coins {
		ok := clpKeeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}
	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.ZeroInt()
	coins = CalculateWithdraw(t, clpKeeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, coin := range coins {
		ok := clpKeeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}
	wBasis = sdk.NewInt(1000)
	asymmetry = sdk.NewInt(-10000)
	coins = CalculateWithdraw(t, clpKeeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.EqualError(t, err, "Cannot remove liquidity asymmetrically")
	require.Nil(t, res)
	for _, coin := range coins {
		ok := clpKeeper.HasBalance(ctx, signer, coin)
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

	UnlockAllliquidity(app, ctx, asset, newLP, t)

	msg = clptypes.NewMsgRemoveLiquidity(newLP, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.EqualError(t, err, "Cannot remove liquidity asymmetrically")
	require.Nil(t, res, "Cannot withdraw now as new LP hasnt added liquidity")
}

func TestSwap(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
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
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin1, nativeCoin))
	require.NoError(t, err)
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin2))
	require.NoError(t, err)
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
	receivedAmount := CalculateSwapReceived(t, clpKeeper, app.TokenRegistryKeeper, ctx, assetEth, assetDash, swapSentAssetETH)
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
	ok := clpKeeper.HasBalance(ctx, signer, CoinsExt1)
	assert.True(t, ok, "")
	ok = clpKeeper.HasBalance(ctx, signer, CoinsNative)
	assert.True(t, ok, "")
	ok = clpKeeper.HasBalance(ctx, signer, CoinsExt2)
	assert.True(t, ok, "")
	msg = clptypes.NewMsgSwap(signer, assetEth, assetDash, swapSentAssetETH, swapSentAssetETH)
	res, err = handler(ctx, &msg)
	require.ErrorIs(t, err, clptypes.ErrReceivedAmountBelowExpected)
	require.Nil(t, res)
	msgE := clptypes.NewMsgSwap(signer, assetEth, clptypes.NewAsset("Asset"), swapSentAssetETH, swapSentAssetETH)
	_, err = handler(ctx, &msgE)
	assert.Error(t, err)
	msgE = clptypes.NewMsgSwap(signer, clptypes.NewAsset("Asset"), assetDash, swapSentAssetETH, swapSentAssetETH)
	_, err = handler(ctx, &msgE)
	assert.Error(t, err)
}

func TestDecommisionPool(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
	//Parameters for Decommission
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")
	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err := handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	// SIGNER became new LP
	lpNewBalance := initialBalance.Sub(sdk.Uint(sdk.Int(poolBalance)))
	lpCoinsExt := sdk.NewCoin(asset.Symbol, sdk.Int(lpNewBalance))
	lpCoinsNative := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(lpNewBalance))
	ok := clpKeeper.HasBalance(ctx, signer, lpCoinsExt)
	assert.True(t, ok, "")
	ok = clpKeeper.HasBalance(ctx, signer, lpCoinsNative)
	assert.True(t, ok, "")
	UnlockAllliquidity(app, ctx, asset, signer, t)

	msgrm := clptypes.NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(5001), sdk.NewInt(1))
	res, err = handler(ctx, &msgrm)
	require.EqualError(t, err, "Cannot remove liquidity asymmetrically")
	require.Nil(t, res)

	msgrm = clptypes.NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(5001), sdk.NewInt(0))
	res, err = handler(ctx, &msgrm)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg := clptypes.NewMsgDecommissionPool(signer, asset.Symbol)
	_, err = handler(ctx, &msg)
	require.Error(t, err)
	v := test.GenerateWhitelistAddress("")
	clpKeeper.SetClpWhiteList(ctx, []sdk.AccAddress{v})
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
	ok = clpKeeper.HasBalance(ctx, signer, lpCoinsExt)
	assert.True(t, ok, "")
	ok = clpKeeper.HasBalance(ctx, signer, lpCoinsNative)
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
	ctx, app := test.CreateTestAppClp(false)
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	_, err = app.TokenRegistryKeeper.GetEntry(registry, pool.ExternalAsset.Symbol)
	swapFeeRate := sdk.NewDecWithPrec(3, 3)
	assert.NoError(t, err)
	if asymmetry.IsPositive() {
		swapResult, _, _, _, err := clpkeeper.SwapOne(clptypes.GetSettlementAsset(), swapAmount, asset, pool, sdk.OneDec(), swapFeeRate)
		assert.NoError(t, err)
		externalAssetCoin = sdk.NewCoin(asset.Symbol, sdk.Int(withdrawExternalAssetAmount.Add(swapResult)))
		nativeAssetCoin = sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.Int(withdrawNativeAssetAmount))
	}
	if asymmetry.IsNegative() {
		swapResult, _, _, _, err := clpkeeper.SwapOne(asset, swapAmount, clptypes.GetSettlementAsset(), pool, sdk.OneDec(), swapFeeRate)
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

func CalculateSwapReceived(t *testing.T, keeper clpkeeper.Keeper, tokenRegistryKeeper tokenregistrytypes.Keeper, ctx sdk.Context, assetSent clptypes.Asset, assetReceived clptypes.Asset, swapAmount sdk.Uint) sdk.Uint {
	swapFeeRate := sdk.NewDecWithPrec(3, 3)
	inPool, err := keeper.GetPool(ctx, assetSent.Symbol)
	assert.NoError(t, err)
	outPool, err := keeper.GetPool(ctx, assetReceived.Symbol)
	assert.NoError(t, err)
	registry := tokenRegistryKeeper.GetRegistry(ctx)
	_, err = tokenRegistryKeeper.GetEntry(registry, inPool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	emitAmount, _, _, _, err := clpkeeper.SwapOne(assetSent, swapAmount, clptypes.GetSettlementAsset(), inPool, sdk.OneDec(), swapFeeRate)
	assert.NoError(t, err)
	_, err = tokenRegistryKeeper.GetEntry(registry, outPool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	emitAmount2, _, _, _, err := clpkeeper.SwapOne(clptypes.GetSettlementAsset(), emitAmount, assetReceived, outPool, sdk.OneDec(), swapFeeRate)
	assert.NoError(t, err)
	return emitAmount2
}

func TestUnlockLiquidity(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	newLP := test.GenerateAddress(test.AddressKey2)
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
	externalDenom := "eth"
	initialBalance := sdk.NewUintFromString("100000000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("10000000000000000000")        // Amount funded to pool , This same amount is used both for native and external asset
	asset := clptypes.NewAsset(externalDenom)
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	err = sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, newLP, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	wBasis := sdk.NewInt(1000)
	asymmetry := sdk.NewInt(10000)
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err := handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	coins := CalculateWithdraw(t, clpKeeper, ctx, asset, signer.String(), wBasis.String(), asymmetry)
	msg := clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res)

	UnlockAllliquidity(app, ctx, asset, signer, t)
	lp, err := app.ClpKeeper.GetLiquidityProvider(ctx, externalDenom, signer.String())
	assert.NoError(t, err)
	beforeUnlocks := lp.Unlocks

	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.EqualError(t, err, "Cannot remove liquidity asymmetrically")
	require.Nil(t, res)

	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(5001), sdk.NewInt(0))
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, coin := range coins {
		ok := clpKeeper.HasBalance(ctx, signer, coin)
		assert.True(t, ok, "")
	}
	ctx = ctx.WithBlockHeight(3)

	lp, err = app.ClpKeeper.GetLiquidityProvider(ctx, externalDenom, signer.String())
	assert.NoError(t, err)
	afterUnlocks := lp.Unlocks
	// Unlocks expired but still not pruned
	assert.NotNil(t, afterUnlocks)
	// Unlocks reduced by liquidity removal
	assert.True(t, beforeUnlocks[0].Units.GT(afterUnlocks[0].Units))

	msg = clptypes.NewMsgRemoveLiquidity(signer, asset, wBasis, asymmetry)
	res, err = handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res)
	// Remove Liquidity prunes unlocks
	lp, err = app.ClpKeeper.GetLiquidityProvider(ctx, externalDenom, signer.String())
	assert.NoError(t, err)
	assert.Nil(t, lp.Unlocks)

	// Test flow with no unbond request made.
	err = app.ClpKeeper.UseUnlockedLiquidity(ctx, clptypes.LiquidityProvider{Asset: &clptypes.Asset{Symbol: "ceth"}}, sdk.NewUint(1), true)
	require.NoError(t, err)
}

func UnlockAllliquidity(app *sifapp.SifchainApp, ctx sdk.Context, asset clptypes.Asset, lp sdk.AccAddress, t *testing.T) {
	nlp, err := app.ClpKeeper.GetLiquidityProvider(ctx, asset.Symbol, lp.String())
	assert.NoError(t, err)
	nlp.Unlocks = append(nlp.Unlocks, &clptypes.LiquidityUnlock{
		RequestHeight: 0,
		Units:         nlp.LiquidityProviderUnits,
	})
	app.ClpKeeper.SetLiquidityProvider(ctx, &nlp)
}
