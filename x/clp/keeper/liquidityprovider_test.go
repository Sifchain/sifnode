package keeper_test

import (
	"fmt"
	"math"
	"testing"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestKeeper_SetLiquidityProvider(t *testing.T) {
	lp := test.GenerateRandomLP(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	clpKeeper.SetLiquidityProvider(ctx, lp)
	getlp, err := clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, *lp)
	lpList, _, err := clpKeeper.GetLiquidityProvidersForAssetPaginated(ctx, *lp.Asset, &query.PageRequest{})
	assert.NoError(t, err)
	assert.Equal(t, lp, lpList[0])
}

func TestKeeper_DestroyLiquidityProvider(t *testing.T) {
	lp := test.GenerateRandomLP(1)[0]
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	clpKeeper.SetLiquidityProvider(ctx, lp)
	getlp, err := clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.NoError(t, err, "Error in get liquidityProvider")
	assert.Equal(t, getlp, *lp)
	assert.True(t, clpKeeper.GetLiquidityProviderIterator(ctx).Valid())
	clpKeeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	_, err = clpKeeper.GetLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.Error(t, err, "LiquidityProvider has been deleted")
	// This should do nothing
	clpKeeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	assert.False(t, clpKeeper.GetLiquidityProviderIterator(ctx).Valid())
}

func TestKeeper_GetAssetsForLiquidityProvider(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	lpList := test.GenerateRandomLP(10)
	for i := range lpList {
		lp := lpList[i]
		clpKeeper.SetLiquidityProvider(ctx, lp)
	}
	lpaddr, err := sdk.AccAddressFromBech32(lpList[0].LiquidityProviderAddress)
	require.NoError(t, err)
	assetList, _, err := clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Limit: math.MaxUint64})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(assetList), len(lpList))
}

func TestKeeper_GetLiquidityProviderData(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	queryLimit := 5
	tokens := []string{"cada", "cbch", "cbnb", "cbtc", "ceos", "ceth", "ctrx", "cusdt"}
	pools, lpList := test.GeneratePoolsAndLPs(clpKeeper, ctx, tokens)
	lpaddr, err := sdk.AccAddressFromBech32(lpList[0].LiquidityProviderAddress)
	require.NoError(t, err)
	assetList, pageRes, err := clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Limit: uint64(queryLimit)})
	require.NoError(t, err)
	require.Len(t, assetList, queryLimit)
	require.NotNil(t, pageRes.NextKey)
	assetList, pageRes, err = clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Key: pageRes.NextKey, Limit: uint64(queryLimit)})
	require.NoError(t, err)
	require.Len(t, assetList, len(tokens)-queryLimit)
	require.Nil(t, pageRes.NextKey)
	assetList, pageRes, err = clpKeeper.GetAssetsForLiquidityProviderPaginated(ctx, lpaddr, &query.PageRequest{Limit: uint64(200)})
	require.NoError(t, err)
	require.Len(t, assetList, len(tokens))
	require.Nil(t, pageRes.NextKey)
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
	require.NotNil(t, lpDataResponse)
	require.Equal(t, len(pools), len(lpDataResponse.LiquidityProviderData))
	require.Equal(t, len(lpList), len(lpDataResponse.LiquidityProviderData))
	for i := 0; i < len(lpDataResponse.LiquidityProviderData); i++ {
		lpData := lpDataResponse.LiquidityProviderData[i]
		require.Contains(t, lpList, *lpData.LiquidityProvider)
		require.Equal(t, lpList[0].LiquidityProviderAddress, lpData.LiquidityProvider.LiquidityProviderAddress)
		require.Equal(t, assetList[i], lpData.LiquidityProvider.Asset)
		require.Equal(t, fmt.Sprint(100*uint64(i+1)), lpData.ExternalAssetBalance)
		require.Equal(t, fmt.Sprint(1000*uint64(i+1)), lpData.NativeAssetBalance)
	}
}

// add tests for GetRewardsEligibleLiquidityProviders
func TestKeeper_GetRewardsEligibleLiquidityProviders(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	tokens := []string{"cada", "cbch", "cbnb", "cbtc", "ceos", "ceth", "ctrx", "cusdt"}
	pools, lpList := test.GeneratePoolsAndLPs(clpKeeper, ctx, tokens)
	// update lp list to contain some lps that has a last updated block that is older
	// than the current block height
	for i := range lpList {
		if (i % 2) == 0 {
			continue
		}
		lpList[i].LastUpdatedBlock = ctx.BlockHeight() - int64(clpKeeper.GetRewardsParams(ctx).RewardsLockPeriod) - 1
		clpKeeper.SetLiquidityProvider(ctx, &lpList[i])
	}
	// get rewards eligible lps
	rewardsEligibleLps, err := clpKeeper.GetRewardsEligibleLiquidityProviders(ctx)
	require.NoError(t, err)
	// check that the rewards eligible lps map contains half the pools
	require.Equal(t, len(pools)/2, len(rewardsEligibleLps))
	// check that the rewards eligible lps map contains half the lps
	for i, lp := range lpList {
		asset := lp.Asset
		assetLps := rewardsEligibleLps[*asset]
		if (i % 2) == 0 {
			require.NotContains(t, assetLps, &lp) //nolint:gosec
		} else {
			require.Contains(t, assetLps, &lp) //nolint:gosec
		}
	}
}
