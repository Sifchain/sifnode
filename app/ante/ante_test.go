package ante_test

import (
	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/app/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestReduceGasPriceDecorator_AnteHandle(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)
	_ = sifapp.AddTestAddrs(app, ctx, 6, initTokens)
	decorator := ante.ReduceGasPriceDecorator{}
	highGasPrice := sdk.DecCoin{
		Denom:  "rowan",
		Amount: sdk.MustNewDecFromStr("0.5"),
	}
	loweredGasPrice := sdk.DecCoin{
		Denom:  "rowan",
		Amount: sdk.MustNewDecFromStr("0.00000005"),
	}
	ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(highGasPrice))
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error){return ctx, nil}
	//Error will always be nil
	newCtx,_ := decorator.AnteHandle(ctx,legacytx.StdTx{},false,next)
	assert.Equal(t, loweredGasPrice.String(),newCtx.MinGasPrices().String())
	assert.Equal(t, highGasPrice.String(),ctx.MinGasPrices().String())
}