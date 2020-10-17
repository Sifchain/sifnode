package clp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePool(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	signer := GenerateAddress()
	asset := NewAsset("ETHEREUM", "ETH", "eth")
	externalCoin := sdk.NewCoin(asset.Ticker, sdk.NewInt(10000))
	nattiveCoin := sdk.NewCoin(NativeTicker, sdk.NewInt(10000))
	keeper.BankKeeper.AddCoins(ctx, signer, sdk.Coins{externalCoin, nattiveCoin})
	msgCreatePool := NewMsgCreatePool(signer, asset, 1000, 1000)
	res, err := handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestAddLiqudity(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	pool := GenerateRandomPool(1)[0]
	signer := GenerateAddress()
	msg := NewMsgAddLiquidity(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err := handleMsgAddLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = NewMsgAddLiquidity(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err = handleMsgAddLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestRemoveLiquidity(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	pool := GenerateRandomPool(1)[0]
	signer := GenerateAddress()
	msg := NewMsgRemoveLiquidity(signer, pool.ExternalAsset, 10000, 1)
	res, err := handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = NewMsgRemoveLiquidity(signer, pool.ExternalAsset, 10, 1)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestSwap(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	asset1 := NewAsset("ETHEREUM", "ETH", "ETH")
	asset2 := NewAsset("TEZOS", "XCT", "XCT")
	signer := GenerateAddress()
	msg := NewMsgSwap(signer, asset1, asset2, 1)
	res, err := handleMsgSwap(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, asset1, 10000, 10000)
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgCreatePool = NewMsgCreatePool(signer, asset2, 10000, 10000)
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = NewMsgSwap(signer, asset1, asset2, 1000)
	res, err = handleMsgSwap(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestDecommisionPool(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	pool := GenerateRandomPool(1)[0]
	signer := GenerateAddress()
	pool.NativeAssetBalance = 100
	pool.ExternalAssetBalance = 1
	msgCreatePool := NewMsgCreatePool(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err := handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgrm := NewMsgRemoveLiquidity(signer, pool.ExternalAsset, 10000, -1)
	res, err = handleMsgRemoveLiquidity(ctx, keeper, msgrm)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg := NewMsgDecommissionPool(signer, pool.ExternalAsset.Ticker)
	res, err = handleMsgDecommissionPool(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgN := NewMsgAddLiquidity(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err = handleMsgAddLiquidity(ctx, keeper, msgN)
	require.Error(t, err)
	require.Nil(t, res)

}
