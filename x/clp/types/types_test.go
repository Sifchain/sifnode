package types

import (
	// "fmt"
	"bytes"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

const (
	AddressKey1 = "A58856F0FD53BF058B4909A21AEC019107BA6"
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

func Test_NewLiquidityProvider(t *testing.T) {
	newAsset := NewAsset("eth")
	address := GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	liquidityProviderUnits := sdk.NewUint(1)
	liquidityProvider := NewLiquidityProvider(&newAsset, liquidityProviderUnits, address)
	assert.Equal(t, liquidityProvider.Asset, &newAsset)
	assert.Equal(t, liquidityProvider.LiquidityProviderAddress, "cosmos15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqjwl8sq")
	assert.Equal(t, liquidityProvider.LiquidityProviderUnits, liquidityProviderUnits)
	boolean := liquidityProvider.Validate()
	assert.True(t, boolean)
	wrongAsset := NewAsset("")
	liquidityProvider = NewLiquidityProvider(&wrongAsset, liquidityProviderUnits, address)
	boolean = liquidityProvider.Validate()
	assert.False(t, boolean)

}

func Test_NewLiquidityProviderResponse(t *testing.T) {
	newAsset := NewAsset("eth")
	address := GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	liquidityProviderUnits := sdk.NewUint(1)
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	liquidityProvider := NewLiquidityProvider(&newAsset, liquidityProviderUnits, address)
	liquidityProviderResponse := NewLiquidityProviderResponse(liquidityProvider, int64(10), nativeAssetAmount.String(), externalAssetAmount.String())
	assert.Equal(t, liquidityProviderResponse.ExternalAssetBalance, externalAssetAmount.String())
	assert.Equal(t, liquidityProviderResponse.NativeAssetBalance, nativeAssetAmount.String())
	assert.Equal(t, liquidityProviderResponse.Height, int64(10))
	assert.Equal(t, liquidityProviderResponse.LiquidityProvider, &liquidityProvider)
}

func GenerateAddress(key string) sdk.AccAddress {
	if key == "" {
		key = AddressKey1
	}
	var buffer bytes.Buffer
	buffer.WriteString(key)
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
