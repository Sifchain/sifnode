package clp

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePool(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	pool := GenerateRandomPool(1)[0]
	signer := GenerateAddress()
	msgCreatePool := NewMsgCreatePool(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
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

func TestNewHandler(t *testing.T) {
	ctx, keeper := CreateTestInputDefault(t, false, 1000)
	asset1 := NewAsset("ETHEREUM", "ETH", "ETH")
	asset2 := NewAsset("TEZOS", "XCT", "XCT")
	pool := GenerateRandomPool(1)[0]
	signer := GenerateAddress()
	pool.NativeAssetBalance = 100
	pool.ExternalAssetBalance = 1
	handler := NewHandler(keeper)
	msgCreatePool := NewMsgCreatePool(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err := handler(ctx, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgAddLiquidity := NewMsgAddLiquidity(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err = handler(ctx, msgAddLiquidity)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgRemoveLiquidity := NewMsgRemoveLiquidity(signer, pool.ExternalAsset, 10, 1)
	res, err = handler(ctx, msgRemoveLiquidity)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgSwap := NewMsgSwap(signer, asset1, asset2, 1000)
	res, err = handler(ctx, msgSwap)
	require.Error(t, err)
	require.Nil(t, res)
	msgRemoveLiquidity = NewMsgRemoveLiquidity(signer, pool.ExternalAsset, 10000, -1)
	res, err = handler(ctx, msgRemoveLiquidity)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgDecommissionPool := NewMsgDecommissionPool(signer, pool.ExternalAsset.Ticker)
	res, err = handler(ctx, msgDecommissionPool)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, nil)
	require.Error(t, err)
	require.Nil(t, res)
}
