package types

import (
	// "fmt"

	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func Test_NewPool(t *testing.T) {
	newAsset := NewAsset("eth0123456789012345678901234567890123456789012345678901234567890123456789")
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	poolUnits := sdk.NewUint(1)
	pool, _ := NewPool(&newAsset, nativeAssetAmount, externalAssetAmount, poolUnits)
	assert.Equal(t, pool.ExternalAssetBalance, externalAssetAmount)
	assert.Equal(t, pool.NativeAssetBalance, nativeAssetAmount)
	assert.Equal(t, pool.ExternalAsset, &newAsset)
}

func Test_TypesValidate(t *testing.T) {
	newAsset := NewAsset("eth")
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	poolUnits := sdk.NewUint(1)
	pool, _ := NewPool(&newAsset, nativeAssetAmount, externalAssetAmount, poolUnits)
	boolean := pool.Validate()
	assert.True(t, boolean)
	newAsset = NewAsset("eth0123456789012345678901234567890123456789012345678901234567890123456789")
	pool, _ = NewPool(&newAsset, nativeAssetAmount, externalAssetAmount, poolUnits)
	boolean = pool.Validate()
	assert.False(t, boolean)
}

/*
func Test_NewPoolsResponse() {
	_ := NewPoolsResponse(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProvider() {
	_ := NewLiquidityProvider(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProviderResponse() {
	_ := NewLiquidityProviderResponse(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProviderDataResponse() {
	_ := NewLiquidityProviderDataResponse(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProviderData() {
	_ := NewLiquidityProviderData(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}
*/
