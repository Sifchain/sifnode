package keeper_test

import (
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_Errors(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, keeper := test.CreateTestAppClp(false)
	_ = keeper.Logger(ctx)
	pool.ExternalAsset.Symbol = ""
	err := keeper.SetPool(ctx, pool)
	assert.Error(t, err)
	getpools := keeper.GetPools(ctx)
	assert.Equal(t, len(getpools), 0, "No pool added")

	lp := test.GenerateRandomLP(1)[0]
	lp.Asset.Symbol = ""
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	assert.Error(t, err)
	assert.NotEqual(t, getlp, lp)
	assert.NotNil(t, test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA7"))
}

func TestKeeper_SetPool(t *testing.T) {

	pool := test.GenerateRandomPool(1)[0]
	ctx, keeper := test.CreateTestAppClp(false)
	err := keeper.SetPool(ctx, pool)
	assert.NoError(t, err)
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	assert.Equal(t, keeper.ExistsPool(ctx, pool.ExternalAsset.Symbol), true)
}

func TestKeeper_GetPools(t *testing.T) {
	pools := test.GenerateRandomPool(10)
	ctx, keeper := test.CreateTestAppClp(false)
	for _, pool := range pools {
		err := keeper.SetPool(ctx, pool)
		assert.NoError(t, err)
	}
	getpools := keeper.GetPools(ctx)
	assert.Greater(t, len(getpools), 0, "More than one pool added")
	assert.LessOrEqual(t, len(getpools), len(pools), "Set pool will ignore duplicates")
}

func TestKeeper_DestroyPool(t *testing.T) {
	pool := test.GenerateRandomPool(1)[0]
	ctx, keeper := test.CreateTestAppClp(false)
	err := keeper.SetPool(ctx, pool)
	assert.NoError(t, err)
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	err = keeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	_, err = keeper.GetPool(ctx, pool.ExternalAsset.Symbol)
	assert.Error(t, err, "Pool should be deleted")
	// This should do nothing.
	err = keeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	assert.Error(t, err)
}

func TestKeeper_SetLiquidityProvider(t *testing.T) {
	lp := test.GenerateRandomLP(1)[0]
	ctx, keeper := test.CreateTestAppClp(false)
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
	lpList := keeper.GetLiquidityProvidersForAsset(ctx, lp.Asset)
	assert.Equal(t, lp, lpList[0])
}

func TestKeeper_DestroyLiquidityProvider(t *testing.T) {
	lp := test.GenerateRandomLP(1)[0]
	ctx, keeper := test.CreateTestAppClp(false)
	keeper.SetLiquidityProvider(ctx, lp)
	getlp, err := keeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, lp)
	assert.True(t, keeper.GetLiquidityProviderIterator(ctx).Valid())
	keeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	_, err = keeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	assert.Error(t, err, "LiquidityProvider has been deleted")
	// This should do nothing
	keeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	assert.False(t, keeper.GetLiquidityProviderIterator(ctx).Valid())
}

func TestKeeper_BankKeeper(t *testing.T) {
	user1 := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	user2 := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA7")
	ctx, keeper := test.CreateTestAppClp(false)
	initialBalance := sdk.NewUint(10000)
	sendingBalance := sdk.NewUint(1000)
	nativeCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(initialBalance))
	sendingCoin := sdk.NewCoin(clp.NativeSymbol, sdk.Int(sendingBalance))
	_, err := keeper.GetBankKeeper().AddCoins(ctx, user1, sdk.Coins{nativeCoin})
	assert.NoError(t, err)
	assert.True(t, keeper.HasCoins(ctx, user1, sdk.Coins{nativeCoin}))
	assert.NoError(t, keeper.SendCoins(ctx, user1, user2, sdk.Coins{sendingCoin}))
	assert.True(t, keeper.HasCoins(ctx, user2, sdk.Coins{sendingCoin}))
}

func TestKeeper_GetAssetsForLiquidityProvider(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	lpList := test.GenerateRandomLP(10)
	for _, lp := range lpList {
		keeper.SetLiquidityProvider(ctx, lp)
	}
	assetList := keeper.GetAssetsForLiquidityProvider(ctx, lpList[0].LiquidityProviderAddress)
	assert.LessOrEqual(t, len(assetList), len(lpList))
}

func TestKeeper_GetModuleAccount(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	moduleAccount := keeper.GetSupplyKeeper().GetModuleAccount(ctx, clp.ModuleName)
	assert.Equal(t, moduleAccount.GetName(), clp.ModuleName)
	assert.Equal(t, moduleAccount.GetPermissions(), []string{supply.Burner, supply.Minter})
}
