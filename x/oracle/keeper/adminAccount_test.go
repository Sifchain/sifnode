package keeper_test

import (
	"testing"

	oracleKeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetAdminAccount(t *testing.T) {
	ctx, keeper, _, _, _, _ := oracleKeeper.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	addresses, _ := oracleKeeper.CreateTestAddrs(2)
	keeper.SetAdminAccount(ctx, addresses[0])
	adminAccount := keeper.GetAdminAccount(ctx)
	assert.Equal(t, adminAccount, addresses[0])
}

func TestKeeper_IsAdminAccount(t *testing.T) {
	ctx, keeper, _, _, _, _ := oracleKeeper.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	addresses, _ := oracleKeeper.CreateTestAddrs(2)
	assert.False(t, keeper.IsAdminAccount(ctx, addresses[0]))
	keeper.SetAdminAccount(ctx, addresses[0])
	assert.True(t, keeper.IsAdminAccount(ctx, addresses[0]))
	assert.False(t, keeper.IsAdminAccount(ctx, addresses[1]))
}
