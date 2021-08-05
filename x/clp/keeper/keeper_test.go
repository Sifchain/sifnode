package keeper_test

import (
	"math"
	"testing"

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
	_ = clpKeeper.Logger(ctx)
	pool.ExternalAsset.Symbol = ""
	err := clpKeeper.SetPool(ctx, &pool)
	assert.Error(t, err)
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

func TestKeeper_SetPool(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	err := clpKeeper.SetPool(ctx, &pool)
	assert.NoError(t, err)
	getpool, err := clpKeeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	assert.Equal(t, clpKeeper.ExistsPool(ctx, pool.ExternalAsset.Symbol), true)
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
	assert.Equal(t, lp, lpList[0])
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
	err := clpKeeper.GetBankKeeper().AddCoins(ctx, user1, sdk.NewCoins(nativeCoin))
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
	assert.LessOrEqual(t, len(assetList), len(lpList))
}

func TestKeeper_GetModuleAccount(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	moduleAccount := clpKeeper.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	assert.Equal(t, moduleAccount.GetName(), types.ModuleName)
	assert.Equal(t, moduleAccount.GetPermissions(), []string{authtypes.Burner, authtypes.Minter})
}
