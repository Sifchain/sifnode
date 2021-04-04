package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Sifchain/sifnode/app"
	oracleKeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	app := app.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	_, addresses := oracleKeeper.CreateTestAddrs(2)
	app.OracleKeeper.SetOracleWhiteList(ctx, addresses)
	vList := app.OracleKeeper.GetOracleWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, app.OracleKeeper.ExistsOracleWhiteList(ctx))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	app := app.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	_, addresses := oracleKeeper.CreateTestAddrs(2)
	app.OracleKeeper.SetOracleWhiteList(ctx, addresses)
	assert.True(t, app.OracleKeeper.ValidateAddress(ctx, addresses[0]))
	assert.True(t, app.OracleKeeper.ValidateAddress(ctx, addresses[1]))
	_, addresses = oracleKeeper.CreateTestAddrs(3)
	assert.False(t, app.OracleKeeper.ValidateAddress(ctx, addresses[2]))
}
