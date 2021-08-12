package clp_test

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"

	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestExportGenesis(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	// pop
	// Generate State
	poolscount, lpCount := CreateState(ctx, app.ClpKeeper, t)
	state := clp.ExportGenesis(ctx, app.ClpKeeper)
	assert.Equal(t, len(state.PoolList), poolscount)
	assert.Equal(t, len(state.LiquidityProviders), lpCount)

}

func TestInitGenesis(t *testing.T) {
	ctx1, app1 := test.CreateTestAppClp(false)
	ctx2, app2 := test.CreateTestAppClp(false)
	// Generate State
	poolscount, lpCount := CreateState(ctx1, app1.ClpKeeper, t)
	state := clp.ExportGenesis(ctx1, app1.ClpKeeper)
	assert.Equal(t, len(state.PoolList), poolscount)
	assert.Equal(t, len(state.LiquidityProviders), lpCount)
	state2 := clp.ExportGenesis(ctx2, app2.ClpKeeper)
	assert.Equal(t, len(state2.PoolList), 0)
	assert.Equal(t, len(state2.LiquidityProviders), 0)

	valUpdates := clp.InitGenesis(ctx2, app2.ClpKeeper, state)
	assert.Equal(t, len(valUpdates), 0)

	poolsList, _, err := app2.ClpKeeper.GetPoolsPaginated(ctx2, &query.PageRequest{Limit: math.MaxUint64})
	assert.NoError(t, err)
	assert.Equal(t, len(poolsList), poolscount)
	lpList, _, err := app2.ClpKeeper.GetAllLiquidityProvidersPaginated(ctx2, &query.PageRequest{Limit: math.MaxUint64})
	assert.NoError(t, err)
	assert.Equal(t, len(lpList), lpCount)
	assert.Equal(t, app2.ClpKeeper.GetParams(ctx2).MinCreatePoolThreshold, types.DefaultMinCreatePoolThreshold)
	assert.Equal(t, app2.ClpKeeper.GetParams(ctx2).MinCreatePoolThreshold, app1.ClpKeeper.GetParams(ctx1).MinCreatePoolThreshold)
}

func TestValidateGenesis(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	// Generate State
	poolscount, lpCount := CreateState(ctx, app.ClpKeeper, t)
	state := clp.ExportGenesis(ctx, app.ClpKeeper)
	assert.Equal(t, len(state.PoolList), poolscount)
	assert.Equal(t, len(state.LiquidityProviders), lpCount)
	err := clp.ValidateGenesis(state)
	assert.NoError(t, err)

}

func CreateState(ctx sdk.Context, keeper keeper.Keeper, t *testing.T) (int, int) {
	// Setting Pools
	pools := test.GenerateRandomPool(10)
	for i := range pools {
		pool := pools[i]
		err := keeper.SetPool(ctx, &pool)
		assert.NoError(t, err)
	}
	poolsList, _, err := keeper.GetPoolsPaginated(ctx, &query.PageRequest{})
	assert.NoError(t, err)
	poolsCount := len(poolsList)
	assert.Greater(t, poolsCount, 0, "More than one pool added")
	assert.LessOrEqual(t, poolsCount, len(pools), "Set pool will ignore duplicates")
	// Setting Liquidity providers
	lpList := test.GenerateRandomLP(10)
	for _, lp := range lpList {
		lp := lp
		keeper.SetLiquidityProvider(ctx, &lp)
	}
	v1 := test.GenerateWhitelistAddress("")
	keeper.SetClpWhiteList(ctx, []sdk.AccAddress{v1})
	accAddr, err := sdk.AccAddressFromBech32(lpList[1].LiquidityProviderAddress)
	assert.NoError(t, err)
	assetList, _, err := keeper.GetAssetsForLiquidityProviderPaginated(ctx, accAddr, &query.PageRequest{Limit: math.MaxUint64})
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(assetList), len(lpList))
	lpCount := len(assetList)
	return poolsCount, lpCount
}
