package keeper_test

import (
	"bytes"
	"strconv"
	"testing"

	clp "github.com/Sifchain/sifnode/x/clp"
	k "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func NewSigner(signer string) sdk.AccAddress {
	var buffer bytes.Buffer
	buffer.WriteString(signer)
	buffer.WriteString(strconv.Itoa(100))
	res, _ := sdk.AccAddressFromHex(buffer.String())
	bech := res.String()
	addr := buffer.String()
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}
	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}
	return res
}

func TestKeeper_CreatePoolAndProvideLiquidity(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.ClpKeeper
	handler := clp.NewHandler(keeper)
	signer := test.GenerateAddress("")
	initialBalance := sdk.NewUintFromString("30000000000000000000")
	nativeAmount := sdk.NewUintFromString("2000000000000000000")
	externalAmount := sdk.NewUintFromString("5000000000000000000")
	asset := types.NewAsset("cusdt")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(initialBalance))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(initialBalance))
	coins := sdk.NewCoins(externalCoin, nativeCoin)
	err := app.BankKeeper.AddCoins(ctx, signer, coins)
	assert.NoError(t, err)

	ok := keeper.HasBalance(ctx, signer, externalCoin)
	assert.True(t, ok, "")

	ok = keeper.HasBalance(ctx, signer, nativeCoin)
	assert.True(t, ok, "")

	msg := types.NewMsgCreatePool(signer, asset, nativeAmount, externalAmount)
	assert.NoError(t, err)

	_, err = handler(ctx, &msg)
	assert.NoError(t, err)

	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	assert.Equal(t, "2000000000000000000", nativeBalance.String())
	assert.Equal(t, "5000000000000000000", externalBalance.String())

	poolUnits, lpunits, err := k.CalculatePoolUnits(msg.ExternalAsset.Symbol, sdk.ZeroUint(),
		sdk.ZeroUint(), sdk.ZeroUint(), nativeBalance, externalBalance)
	assert.NoError(t, err)
	assert.Equal(t, "2000000000000000000", lpunits.String())
	assert.Equal(t, "2000000000000000000", poolUnits.String())

	pool, err := keeper.CreatePool(ctx, poolUnits, &msg)
	assert.NoError(t, err)
	assert.Equal(t, externalCoin.Denom, pool.ExternalAsset.Symbol)
	assert.Equal(t, "2000000000000000000", pool.PoolUnits.String())
	assert.Equal(t, "2000000000000000000", pool.NativeAssetBalance.String())
	assert.Equal(t, "5000000000000000000", pool.ExternalAssetBalance.String())

	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	assert.NoError(t, err)

	lp := keeper.CreateLiquidityProvider(ctx, msg.ExternalAsset, lpunits, signer)
	assert.Equal(t, addr.String(), lp.LiquidityProviderAddress)
	assert.Equal(t, lpunits, lp.LiquidityProviderUnits)

	msgAddLiquidity := types.NewMsgAddLiquidity(signer, asset, nativeAmount, externalAmount)
	assert.NoError(t, err)

	_, err = handler(ctx, &msgAddLiquidity)
	assert.NoError(t, err)

	lp2, err := keeper.AddLiquidity(ctx, &msgAddLiquidity, *pool, externalAmount, lpunits)
	assert.NoError(t, err)
	assert.Equal(t, "6000000000000000000", lp2.LiquidityProviderUnits.String())
	assert.Equal(t, "2000000000000000000", pool.NativeAssetBalance.String())
	assert.Equal(t, "5000000000000000000", pool.ExternalAssetBalance.String())
}
