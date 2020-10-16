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
	pool1 := GenerateRandomPool(1)[0]
	pool2 := GenerateRandomPool(1)[0]
	signer := GenerateAddress()
	msg := NewMsgSwap(signer, pool1.ExternalAsset, pool2.ExternalAsset, 1)
	res, err := handleMsgSwap(ctx, keeper, msg)
	require.Error(t, err)
	require.Nil(t, res)
	msgCreatePool := NewMsgCreatePool(signer, pool1.ExternalAsset, pool1.NativeAssetBalance, pool1.ExternalAssetBalance)
	res, err = handleMsgCreatePool(ctx, keeper, msgCreatePool)
	require.NoError(t, err)
	require.NotNil(t, res)
	msg = NewMsgSwap(signer, pool1.ExternalAsset, pool2.ExternalAsset, 1)
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
	msg := NewMsgDecommissionPool(signer, pool.ExternalAsset.Ticker, pool.ExternalAsset.SourceChain)
	res, err = handleMsgDecommissionPool(ctx, keeper, msg)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgN := NewMsgAddLiquidity(signer, pool.ExternalAsset, pool.NativeAssetBalance, pool.ExternalAssetBalance)
	res, err = handleMsgAddLiquidity(ctx, keeper, msgN)
	require.Error(t, err)
	require.Nil(t, res)

}
