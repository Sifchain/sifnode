package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/app"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	oracleKeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
)

func TestKeeper_SetAdminAccount(t *testing.T) {
	app := app.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addresses, _ := oracleKeeper.CreateTestAddrs(2)
	app.OracleKeeper.SetAdminAccount(ctx, addresses[0])
	adminAccount := app.OracleKeeper.GetAdminAccount(ctx)
	assert.Equal(t, adminAccount, addresses[0])
}

func TestKeeper_IsAdminAccount(t *testing.T) {
	app := app.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addresses, _ := oracleKeeper.CreateTestAddrs(2)
	assert.False(t, app.OracleKeeper.IsAdminAccount(ctx, addresses[0]))
	app.OracleKeeper.SetAdminAccount(ctx, addresses[0])
	assert.True(t, app.OracleKeeper.IsAdminAccount(ctx, addresses[0]))
	assert.False(t, app.OracleKeeper.IsAdminAccount(ctx, addresses[1]))
}
