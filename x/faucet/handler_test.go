package faucet_test

import (
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/faucet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandler(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	handler := faucet.NewHandler(app.FaucetKeeper)
	res, err := handler(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res)

}

func TestRequestCoins(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	handler := faucet.NewHandler(app.FaucetKeeper)
	keeper := app.FaucetKeeper
	signer := test.GenerateAddress("")
	requestCoins := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))}
	faucetFunding := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(10000000000))}
	msg := faucet.NewMsgRequestCoins(signer, requestCoins)
	res, err := handler(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res)
	ok := keeper.GetBankKeeper().HasCoins(ctx, signer, requestCoins)
	assert.False(t, ok, "Faucet does not have funds")

	_, err = keeper.GetBankKeeper().AddCoins(ctx, faucet.GetFaucetModuleAddress(), faucetFunding)
	require.NoError(t, err)

	res, err = handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	ok = keeper.GetBankKeeper().HasCoins(ctx, signer, requestCoins)
	assert.True(t, ok, "")

}

func TestAddCoins(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	handler := faucet.NewHandler(app.FaucetKeeper)
	keeper := app.FaucetKeeper
	signer := test.GenerateAddress("")
	requestCoins := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))}
	faucetFunding := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(10000000000))}
	_, err := keeper.GetBankKeeper().AddCoins(ctx, signer, faucetFunding)
	require.NoError(t, err)

	msg := faucet.NewMsgAddCoins(signer, faucetFunding)
	res, err := handler(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	ok := keeper.GetBankKeeper().HasCoins(ctx, faucet.GetFaucetModuleAddress(), faucetFunding)
	assert.True(t, ok, "")

	msgR := faucet.NewMsgRequestCoins(signer, requestCoins)
	res, err = handler(ctx, msgR)
	require.NoError(t, err)
	require.NotNil(t, res)
	ok = keeper.GetBankKeeper().HasCoins(ctx, signer, requestCoins)
	assert.True(t, ok, "")

}
