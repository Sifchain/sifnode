package clp_test

import (
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExportGenesis(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	// Generate State
	poolscount, lpCount := CreateState(ctx, keeper, t)
	state := clp.ExportGenesis(ctx, keeper)
	assert.Equal(t, len(state.PoolList), poolscount)
	assert.Equal(t, len(state.LiquidityProviderList), lpCount)

}

func TestInitGenesis(t *testing.T) {
	ctx1, keeper1 := test.CreateTestAppClp(false)
	ctx2, keeper2 := test.CreateTestAppClp(false)
	// Generate State
	poolscount, lpCount := CreateState(ctx1, keeper1, t)
	state := clp.ExportGenesis(ctx1, keeper1)
	assert.Equal(t, len(state.PoolList), poolscount)
	assert.Equal(t, len(state.LiquidityProviderList), lpCount)
	state2 := clp.ExportGenesis(ctx2, keeper2)
	assert.Equal(t, len(state2.PoolList), 0)
	assert.Equal(t, len(state2.LiquidityProviderList), 0)

	valUpdates := clp.InitGenesis(ctx2, keeper2, state)
	assert.Equal(t, len(valUpdates), 0)

	poollist := keeper2.GetPools(ctx2)
	assert.Equal(t, len(poollist), poolscount)
	lpList := keeper2.GetLiquidityProviders(ctx2)
	assert.Equal(t, len(lpList), lpCount)
	assert.Equal(t, keeper2.GetParams(ctx2).MinCreatePoolThreshold, types.DefaultMinCreatePoolThreshold)
	assert.Equal(t, keeper2.GetParams(ctx2).MinCreatePoolThreshold, keeper1.GetParams(ctx1).MinCreatePoolThreshold)
}

func TestValidateGenesis(t *testing.T) {
	ctx, keeper := test.CreateTestAppClp(false)
	// Generate State
	poolscount, lpCount := CreateState(ctx, keeper, t)
	state := clp.ExportGenesis(ctx, keeper)
	assert.Equal(t, len(state.PoolList), poolscount)
	assert.Equal(t, len(state.LiquidityProviderList), lpCount)
	err := clp.ValidateGenesis(state)
	assert.NoError(t, err)

}
func CreateState(ctx sdk.Context, keeper clp.Keeper, t *testing.T) (int, int) {
	// Setting Pools
	pools := test.GenerateRandomPool(10)
	for _, pool := range pools {
		err := keeper.SetPool(ctx, pool)
		assert.NoError(t, err)
	}
	getpools := keeper.GetPools(ctx)
	assert.Greater(t, len(getpools), 0, "More than one pool added")
	assert.LessOrEqual(t, len(getpools), len(pools), "Set pool will ignore duplicates")

	poolscount := len(getpools)

	//Setting Liquidity providers
	lpList := test.GenerateRandomLP(10)
	for _, lp := range lpList {
		keeper.SetLiquidityProvider(ctx, lp)
	}
	v1 := test.GenerateValidatorAddress("")
	keeper.SetValidatorWhiteList(ctx, []sdk.ValAddress{v1})
	assetList := keeper.GetAssetsForLiquidityProvider(ctx, lpList[0].LiquidityProviderAddress)
	assert.LessOrEqual(t, len(assetList), len(lpList))
	lpCount := len(assetList)
	return poolscount, lpCount
}
