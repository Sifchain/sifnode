package keeper_test

import (
	"fmt"
	"math"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestKeeper_Errors(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	num := clpKeeper.GetMinCreatePoolThreshold(ctx)
	assert.Equal(t, num, uint64(100))
	_ = clpKeeper.GetParams(ctx)
	_ = clpKeeper.Logger(ctx)
	pool.ExternalAsset.Symbol = ""
	err := clpKeeper.SetPool(ctx, &pool)
	assert.Error(t, err, "Unable to set pool")
	boolean := clpKeeper.ValidatePool(pool)
	assert.False(t, boolean)
	getpools, _, err := clpKeeper.GetPoolsPaginated(ctx, &query.PageRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(getpools), 0, "No pool added")
	lp := test.GenerateRandomLP(1)[0]
	lp.Asset.Symbol = ""
	clpKeeper.SetLiquidityProvider(ctx, &lp)
	getlp, err := clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.Error(t, err)
	assert.NotEqual(t, getlp, lp)
	assert.NotNil(t, test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA7"))
}

func TestKeeper_CalculateAssetsForLP(t *testing.T) {
	_, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	tokens := []string{"cada", "cbch", "cbnb", "cbtc", "ceos", "ceth", "ctrx", "cusdt"}
	pools, lpList := test.GeneratePoolsAndLPs(keeper, ctx, tokens)
	native, external, _, _ := clpkeeper.CalculateAllAssetsForLP(pools[0], lpList[0])
	assert.Equal(t, "100", external.String())
	assert.Equal(t, "1000", native.String())
}

func TestKeeper_SetPool_ValidatePool(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	err := clpKeeper.SetPool(ctx, &pool)
	assert.NoError(t, err)
	getpool, err := clpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	assert.Equal(t, clpKeeper.ExistsPool(ctx, pool.ExternalAsset.Symbol), true)
	boolean := clpKeeper.ValidatePool(pool)
	assert.True(t, boolean)
}

func TestKeeper_GetPools(t *testing.T) {
	pools := test.GenerateRandomPool(10)
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	for i := range pools {
		pool := pools[i]
		err := clpKeeper.SetPool(ctx, &pool)
		assert.NoError(t, err)
	}
	getpools, _, err := clpKeeper.GetPoolsPaginated(ctx, &query.PageRequest{})
	assert.NoError(t, err)
	assert.Greater(t, len(getpools), 0, "More than one pool added")
	assert.LessOrEqual(t, len(getpools), len(pools), "Set pool will ignore duplicates")
}

func TestKeeper_DestroyPool(t *testing.T) {
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

func TestKeeper_SetLiquidityProvider(t *testing.T) {
	lp := test.GenerateRandomLP(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	clpKeeper.SetLiquidityProvider(ctx, &lp)
	getlp, err := clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
	lpList, _, err := clpKeeper.GetLiquidityProvidersForAssetPaginated(ctx, *lp.Asset, &query.PageRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &lp, lpList[0])
}

func TestKeeper_DestroyLiquidityProvider(t *testing.T) {
	lp := test.GenerateRandomLP(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	clpKeeper.SetLiquidityProvider(ctx, &lp)
	getlp, err := clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
	assert.True(t, clpKeeper.GetLiquidityProviderIterator(ctx).Valid())
	clpKeeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	_, err = clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.Error(t, err, "LiquidityProvider has been deleted")
	// This should do nothing
	clpKeeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.False(t, clpKeeper.GetLiquidityProviderIterator(ctx).Valid())
}

func TestKeeper_BankKeeper(t *testing.T) {
	user1 := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	user2 := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA7")
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	initialBalance := sdk.NewUint(10000)
	sendingBalance := sdk.NewUint(1000)
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(initialBalance))
	sendingCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sendingBalance))
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, user1, sdk.NewCoins(nativeCoin))
	assert.NoError(t, err)
	assert.True(t, clpKeeper.HasBalance(ctx, user1, nativeCoin))
	assert.NoError(t, clpKeeper.SendCoins(ctx, user1, user2, sdk.NewCoins(sendingCoin)))
	assert.True(t, clpKeeper.HasBalance(ctx, user2, sendingCoin))
}

func TestKeeper_GetAssetsForLiquidityProvider(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	lpList := test.GenerateRandomLP(10)
	for i := range lpList {
		lp := lpList[i]
		clpKeeper.SetLiquidityProvider(ctx, &lp)
	}
	lpaddr, err := sdk.AccAddressFromBech32(lpList[0].LiquidityProviderAddress)
	require.NoError(t, err)
	assetList, _, err := clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Limit: math.MaxUint64})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(assetList), len(lpList))
}

func TestKeeper_GetModuleAccount(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	moduleAccount := clpKeeper.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	assert.Equal(t, moduleAccount.GetName(), types.ModuleName)
	assert.Equal(t, moduleAccount.GetPermissions(), []string{authtypes.Burner, authtypes.Minter})
}

func TestKeeper_GetLiquidityProviderData(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	queryLimit := 5
	tokens := []string{"cada", "cbch", "cbnb", "cbtc", "ceos", "ceth", "ctrx", "cusdt"}
	pools, lpList := test.GeneratePoolsAndLPs(clpKeeper, ctx, tokens)
	lpaddr, err := sdk.AccAddressFromBech32(lpList[0].LiquidityProviderAddress)
	assert.NoError(t, err)
	assetList, pageRes, err := clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Limit: uint64(queryLimit)})
	assert.NoError(t, err)
	assert.Len(t, assetList, queryLimit)
	assert.NotNil(t, pageRes.NextKey)
	assetList, pageRes, err = clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Key: pageRes.NextKey, Limit: uint64(queryLimit)})
	assert.NoError(t, err)
	assert.Len(t, assetList, len(tokens)-queryLimit)
	assert.Nil(t, pageRes.NextKey)
	assetList, pageRes, err = clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Limit: uint64(200)})
	assert.NoError(t, err)
	assert.Len(t, assetList, len(tokens))
	assert.Nil(t, pageRes.NextKey)
	lpDataList := make([]*types.LiquidityProviderData, 0, len(assetList))
	for i := range assetList {
		asset := assetList[i]
		pool, err := clpKeeper.GetPool(ctx, asset.Symbol)
		if err != nil {
			continue
		}
		lp, err := clpKeeper.GetLiquidityProvider(ctx, asset.Symbol, lpaddr.String())
		if err != nil {
			continue
		}
		native, external, _, _ := clpkeeper.CalculateAllAssetsForLP(pool, lp)
		lpData := types.NewLiquidityProviderData(lp, native.String(), external.String())
		lpDataList = append(lpDataList, &lpData)
	}
	lpDataResponse := types.NewLiquidityProviderDataResponse(lpDataList, ctx.BlockHeight())
	assert.NotNil(t, lpDataResponse)
	assert.Equal(t, len(pools), len(lpDataResponse.LiquidityProviderData))
	assert.Equal(t, len(lpList), len(lpDataResponse.LiquidityProviderData))
	for i := 0; i < len(lpDataResponse.LiquidityProviderData); i++ {
		lpData := lpDataResponse.LiquidityProviderData[i]
		assert.Contains(t, lpList, *lpData.LiquidityProvider)
		assert.Equal(t, lpList[0].LiquidityProviderAddress, lpData.LiquidityProvider.LiquidityProviderAddress)
		assert.Equal(t, assetList[i], lpData.LiquidityProvider.Asset)
		assert.Equal(t, fmt.Sprint(100*uint64(i+1)), lpData.ExternalAssetBalance)
		assert.Equal(t, fmt.Sprint(1000*uint64(i+1)), lpData.NativeAssetBalance)
	}
}

func TestKeeper_SwapOne(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	wBasis := sdk.NewInt(1000)
	asymmetry := sdk.NewInt(10000)
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.NoError(t, err)
	msg := types.NewMsgAddLiquidity(signer, asset, nativeAssetAmount, externalAssetAmount)
	app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), signer)
	lp, err := app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.NoError(t, err)
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := app.TokenRegistryKeeper.GetEntry(registry, pool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	// asymmetry is positive
	normalizationFactor, adjustExternalToken := app.ClpKeeper.GetNormalizationFactor(eAsset.Decimals)
	_, _, _, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), wBasis.String(), asymmetry)
	swapResult, liquidityFee, priceImpact, _, err := clpkeeper.SwapOne(types.GetSettlementAsset(), swapAmount, asset, *pool, normalizationFactor, adjustExternalToken)
	assert.NoError(t, err)
	assert.Equal(t, swapResult.String(), "10")
	assert.Equal(t, liquidityFee.String(), "978")
	assert.Equal(t, priceImpact.String(), "0")
}

func TestKeeper_SetInputs(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, _ := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	X, x, Y, toRowan := clpkeeper.SetInputs(sdk.NewUint(1), asset, *pool)
	assert.Equal(t, X, sdk.NewUint(998))
	assert.Equal(t, x, sdk.NewUint(1))
	assert.Equal(t, Y, sdk.NewUint(998))
	assert.Equal(t, toRowan, false)
}

func TestKeeper_GetSwapFee(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, _ := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	eAsset, _ := app.TokenRegistryKeeper.GetEntry(registry, pool.ExternalAsset.Symbol)
	normalizationFactor, adjustExternalToken := app.ClpKeeper.GetNormalizationFactor(eAsset.Decimals)
	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, *pool, normalizationFactor, adjustExternalToken)
	assert.Equal(t, swapResult.String(), "1")
}
