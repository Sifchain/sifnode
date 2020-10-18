package clp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePool(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()

	intitalBalance := 10000
	poolfundingAmount := 1000
	asset := NewAsset("ETHEREUM", "ETH", "eth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(intitalBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance)))
	keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")
	msgCreatePool := NewMsgCreatePool(signer, asset, 1000, 1000)
	res, err := handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	externalCoin = sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(intitalBalance-poolfundingAmount)))
	nativeCoin = sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance-poolfundingAmount)))
	ok = keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")

}

func TestAddLiqudity(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	intitalBalance := 10000
	addLiqudityAmount := 1000
	asset := NewAsset("ETHEREUM", "ETH", "eth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(intitalBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance)))
	keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msg := NewMsgAddLiquidity(signer, asset, uint(addLiqudityAmount), uint(addLiqudityAmount))
	res, err := handleMsgAddLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, asset, uint(addLiqudityAmount), uint(addLiqudityAmount))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = NewMsgAddLiquidity(signer, asset, uint(addLiqudityAmount), uint(addLiqudityAmount))
	res, err = handleMsgAddLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	// Subtracted twice , during create and add
	externalCoin = sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(intitalBalance-addLiqudityAmount-addLiqudityAmount)))
	nativeCoin = sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance-addLiqudityAmount-addLiqudityAmount)))
	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")
}

func TestRemoveLiquidity(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()

	intitalBalance := 10000
	wBasis := 1000
	asset := NewAsset("ETHEREUM", "ETH", "eth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(intitalBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance)))
	keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msg := NewMsgRemoveLiquidity(signer, asset, wBasis, 1)
	res, err := handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, asset, uint(wBasis), uint(wBasis))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = NewMsgRemoveLiquidity(signer, asset, wBasis, 1)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	//subtracted during create added during remove
	externalCoin = sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(intitalBalance-wBasis+(wBasis/10))))
	nativeCoin = sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance-wBasis+(wBasis/10))))
	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})
	assert.True(t, ok, "")
}

func TestSwap(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	asset1 := NewAsset("ETHEREUM", "ETH", "eth")
	asset2 := NewAsset("DASH", "DASH", "dash")

	intitalBalance := 10000
	poolbalance := 1000
	externalCoin1 := sdk.NewCoin(asset1.Ticker, sdk.NewInt(int64(intitalBalance)))
	externalCoin2 := sdk.NewCoin(asset2.Ticker, sdk.NewInt(int64(intitalBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance)))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)

	keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin1, nativeCoin})
	keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin2})

	msg := NewMsgSwap(signer, asset1, asset2, 1)
	res, err := handleMsgSwap(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, asset1, uint(poolbalance), uint(poolbalance))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	msgCreatePool = NewMsgCreatePool(signer, asset2, uint(poolbalance), uint(poolbalance))
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg = NewMsgSwap(signer, asset1, asset2, 100)
	res, err = handleMsgSwap(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	// TODO USE FORMULA
	CoinsExt1 := sdk.NewCoin(asset1.Ticker, sdk.NewInt(8900))
	CoinsNative := sdk.NewCoin(NativeTicker, sdk.NewInt(8000))
	CoinsExt2 := sdk.NewCoin(asset2.Ticker, sdk.NewInt(9082))
	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{CoinsExt1, CoinsNative, CoinsExt2})
	assert.True(t, ok, "")

}

func TestDecommisionPool(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	pool := GenerateRandomPool(1)[0]
	signer := GenerateAddress()
	pool.NativeAssetBalance = 100
	pool.ExternalAssetBalance = 1

	intitalBalance := 10000
	poolbalance := 100
	asset := NewAsset("ETHEREUM", "ETH", "eth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(intitalBalance)))
	nativeCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(intitalBalance)))
	// Signer is given ETH and RWN ( Signer will creat pool and become LP)
	keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nativeCoin})

	msgCreatePool := NewMsgCreatePool(signer, asset, uint(poolbalance), uint(poolbalance))
	res, err := handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)

	// SIGNER became new LP
	lpNewBalance := intitalBalance - poolbalance
	lpCoinsExt := sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(lpNewBalance)))
	lpCoinsNative := sdk.NewCoin(NativeTicker, sdk.NewInt(int64(lpNewBalance)))
	ok := keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{lpCoinsExt, lpCoinsNative})
	assert.True(t, ok, "")

	msgrm := NewMsgRemoveLiquidity(signer, asset, 5001, -1)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msgrm)
	require.NoError(t, err)
	require.NotNil(t, res)

	msg := NewMsgDecommissionPool(signer, asset.Ticker)
	res, err = handleMsgDecommissionPool(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgN := NewMsgAddLiquidity(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err = handleMsgAddLiquidity(ctx, keeper, msgN)
	require.Error(t, err)
	require.Nil(t, res)

	// LP refunded coins when decommison
	lpNewBalance = intitalBalance
	lpCoinsExt = sdk.NewCoin(asset.Ticker, sdk.NewInt(int64(lpNewBalance)))
	lpCoinsNative = sdk.NewCoin(NativeTicker, sdk.NewInt(int64(lpNewBalance)))
	ok = keeper.BankKeeper.HasCoins(ctx, signer, sdk.Coins{lpCoinsExt, lpCoinsNative})
	assert.True(t, ok, "")

}
