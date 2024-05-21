package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewAsset(t *testing.T) {
	asset := NewAsset("eth")
	assert.Equal(t, asset.Symbol, "eth")
}

func Test_AssetValidate(t *testing.T) {
	asset := NewAsset("eth")
	boolean := asset.Validate()
	assert.True(t, boolean)

	// to long
	asset = NewAsset("abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg")
	boolean = asset.Validate()
	assert.False(t, boolean)
}

func Test_Equals(t *testing.T) {
	asset := NewAsset("eth")
	asset1 := NewAsset("eth")
	boolean := asset.Equals(asset1)
	assert.True(t, boolean)
	asset1 = NewAsset("ethx")
	boolean = asset.Equals(asset1)
	assert.False(t, boolean)
}

func Test_IsEmpty(t *testing.T) {
	asset := NewAsset("eth")
	boolean := asset.IsEmpty()
	assert.False(t, boolean)
	asset = NewAsset("")
	boolean = asset.IsEmpty()
	assert.True(t, boolean)
}

func Test_GetSettlementAsset(t *testing.T) {
	asset := GetSettlementAsset()
	assert.Equal(t, asset, NewAsset("rowan"))
}

func Test_GetCLPModuleAddress(t *testing.T) {
	clpModuleAddress := GetCLPModuleAddress()
	assert.Equal(t, clpModuleAddress.String(), "cosmos1pjm228rsgwqf23arkx7lm9ypkyma7mzr5e99gl")
}

func Test_GetDefaultCLPAdmin(t *testing.T) {
	defaultCLPAdmin := GetDefaultCLPAdmin()
	assert.Equal(t, defaultCLPAdmin.String(), "cosmos1ny48eeuk4dm9f63dy0lwfgjhnvud9yvt8tcaat")
}
