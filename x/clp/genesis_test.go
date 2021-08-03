package clp_test

import (
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
	ctx, keeper := test.CreateTestAppClp(false)
	// Generate State
	poolsCount, lpCount := CreateState(ctx, keeper, t)
	state := clp.ExportGenesis(ctx, keeper)
	assert.Equal(t, len(state.PoolList), poolsCount)
	assert.Equal(t, len(state.LiquidityProviders), lpCount)
}

func TestInitGenesis(t *testing.T) {
	ctx1, keeper1 := test.CreateTestAppClp(false)
	ctx2, keeper2 := test.CreateTestAppClp(false)
	// Generate State
	poolsCount, lpCount := CreateState(ctx1, keeper1, t)
	state := clp.ExportGenesis(ctx1, keeper1)
	assert.Equal(t, len(state.PoolList), poolsCount)
	assert.Equal(t, len(state.LiquidityProviders), lpCount)
	state2 := clp.ExportGenesis(ctx2, keeper2)
	assert.Equal(t, len(state2.PoolList), 0)
	assert.Equal(t, len(state2.LiquidityProviders), 0)
	valUpdates := clp.InitGenesis(ctx2, keeper2, state)
	assert.Equal(t, len(valUpdates), 0)
	poolsList, _, err := keeper2.GetPoolsPaginated(ctx2, &query.PageRequest{})
	assert.NoError(t, err)
	assert.Equal(t, len(poolsList), poolsCount)
	lpList := keeper2.GetLiquidityProviders(ctx2)
	assert.Equal(t, len(lpList), lpCount)
	assert.Equal(t, keeper2.GetParams(ctx2).MinCreatePoolThreshold, types.DefaultMinCreatePoolThreshold)
	assert.Equal(t, keeper2.GetParams(ctx2).MinCreatePoolThreshold, keeper1.GetParams(ctx1).MinCreatePoolThreshold)
}

func TestValidateGenesis(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	// Generate State
	poolsCount, lpCount := CreateState(ctx, keeper, t)
	state := clp.ExportGenesis(ctx, keeper)
	assert.Equal(t, len(state.PoolList), poolsCount)
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
	assetList := keeper.GetAssetsForLiquidityProvider(ctx, accAddr)
	assert.LessOrEqual(t, len(assetList), len(lpList))
	lpCount := len(assetList)
	return poolsCount, lpCount
}
