package clp_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
)

func TestHandler(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	handler := clp.NewHandler(app.ClpKeeper)
	res, err := handler(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestCreatePool(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	handler := clp.NewHandler(app.ClpKeeper)
	signer := test.GenerateAddress("")
	//Parameters for create pool
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	ok := app.ClpKeeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")
	ok = app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")
}

func TestAddLiquidity(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress("")
	clpKeeper := app.ClpKeeper
	handler := clp.NewHandler(clpKeeper)
	//Parameters for add liquidity
	initialBalance := sdk.NewUintFromString("100000000000000000000") // Initial account balance for all assets created
	poolBalance := sdk.NewUintFromString("1000000000000000000")      // Amount funded to pool , This same amount is used both for native and external asset
	addLiquidityAmount := sdk.NewUintFromString("1000000000000000000")
	asset := clptypes.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(initialBalance))
	err := sifapp.AddCoinsToAccount(clptypes.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msg := clptypes.NewMsgAddLiquidity(signer, asset, addLiquidityAmount, addLiquidityAmount)
	res, err := handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := clptypes.NewMsgCreatePool(signer, asset, poolBalance, poolBalance)
	res, err = handler(ctx, &msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = clptypes.NewMsgAddLiquidity(signer, asset, sdk.ZeroUint(), addLiquidityAmount)
	res, err = handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, res)

}
