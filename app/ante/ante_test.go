package ante_test

import (
	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/app/ante"
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
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
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) { return ctx, nil }
	//Error will always be nil
	msg := dispensationtypes.NewMsgCreateDistribution(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), dispensationtypes.DistributionType_DISTRIBUTION_TYPE_AIRDROP, []banktypes.Output{}, "")
	tx := legacytx.StdTx{
		Msgs:          []sdk.Msg{&msg},
		Fee:           legacytx.StdFee{},
		Signatures:    nil,
		Memo:          "",
		TimeoutHeight: 0,
	}
	newCtx, _ := decorator.AnteHandle(ctx, tx, false, next)
	assert.Equal(t, loweredGasPrice.String(), newCtx.MinGasPrices().String())
	assert.Equal(t, highGasPrice.String(), ctx.MinGasPrices().String())

	newCtx, _ = decorator.AnteHandle(ctx, legacytx.StdTx{}, false, next)
	assert.Equal(t, highGasPrice.String(), newCtx.MinGasPrices().String())

}
