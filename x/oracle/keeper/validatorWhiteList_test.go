package keeper_test

import (
	"testing"

	oracleKeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	ctx, keeper, _, _, _, _ := oracleKeeper.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	_, addresses := oracleKeeper.CreateTestAddrs(2)
	keeper.SetOracleWhiteList(ctx, addresses)
	vList := keeper.GetOracleWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, keeper.ExistsOracleWhiteList(ctx))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	ctx, keeper, _, _, _, _ := oracleKeeper.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	_, addresses := oracleKeeper.CreateTestAddrs(2)
	keeper.SetOracleWhiteList(ctx, addresses)
	assert.True(t, keeper.ValidateAddress(ctx, addresses[0]))
	assert.True(t, keeper.ValidateAddress(ctx, addresses[1]))
	_, addresses = oracleKeeper.CreateTestAddrs(3)
	assert.False(t, keeper.ValidateAddress(ctx, addresses[2]))
}
