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
