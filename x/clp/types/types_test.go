package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

const (
	_ = "A58856F0FD53BF058B4909A21AEC019107BA6"
)

func Test_NewPool(t *testing.T) {
	newAsset := NewAsset("eth0123456789012345678901234567890123456789012345678901234567890123456789")
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	poolUnits := sdk.NewUint(1)
	pool := NewPool(&newAsset, nativeAssetAmount, externalAssetAmount, poolUnits)
	assert.Equal(t, pool.ExternalAssetBalance, externalAssetAmount)
	assert.Equal(t, pool.NativeAssetBalance, nativeAssetAmount)
	assert.Equal(t, pool.ExternalAsset, &newAsset)
}

func Test_TypesValidate(t *testing.T) {
	newAsset := NewAsset("eth")
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	poolUnits := sdk.NewUint(1)
	pool := NewPool(&newAsset, nativeAssetAmount, externalAssetAmount, poolUnits)
	boolean := pool.Validate()
	assert.True(t, boolean)
	newAsset = NewAsset("eth0123456789012345678901234567890123456789012345678901234567890123456789")
	pool = NewPool(&newAsset, nativeAssetAmount, externalAssetAmount, poolUnits)
	boolean = pool.Validate()
	assert.False(t, boolean)
}

func Test_NewLiquidityProvider(t *testing.T) {
	newAsset := NewAsset("eth")
	address, _ := sdk.AccAddressFromBech32("A58856F0FD53BF058B4909A21AEC019107BA6")
	liquidityProviderUnits := sdk.NewUint(1)
	liquidityProvider := NewLiquidityProvider(&newAsset, liquidityProviderUnits, address, 0)
	assert.Equal(t, liquidityProvider.Asset, &newAsset)
	assert.Equal(t, liquidityProvider.LiquidityProviderAddress, "")
	assert.Equal(t, liquidityProvider.LiquidityProviderUnits, liquidityProviderUnits)
	boolean := liquidityProvider.Validate()
	assert.True(t, boolean)
	wrongAsset := NewAsset("")
	liquidityProvider = NewLiquidityProvider(&wrongAsset, liquidityProviderUnits, address, 0)
	boolean = liquidityProvider.Validate()
	assert.False(t, boolean)
}

func Test_NewLiquidityProviderResponse(t *testing.T) {
	newAsset := NewAsset("eth")
	address, _ := sdk.AccAddressFromBech32("A58856F0FD53BF058B4909A21AEC019107BA6")
	liquidityProviderUnits := sdk.NewUint(1)
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	liquidityProvider := NewLiquidityProvider(&newAsset, liquidityProviderUnits, address, 0)
	liquidityProviderResponse := NewLiquidityProviderResponse(liquidityProvider, int64(10), nativeAssetAmount.String(), externalAssetAmount.String())
	assert.Equal(t, liquidityProviderResponse.ExternalAssetBalance, externalAssetAmount.String())
	assert.Equal(t, liquidityProviderResponse.NativeAssetBalance, nativeAssetAmount.String())
	assert.Equal(t, liquidityProviderResponse.Height, int64(10))
	assert.Equal(t, liquidityProviderResponse.LiquidityProvider, &liquidityProvider)
}
