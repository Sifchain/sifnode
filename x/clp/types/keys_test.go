package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeys_GetPoolKey(t *testing.T) {
	poolKey, err := GetPoolKey("ceth", "rwn")
	assert.NoError(t, err)
	poolKey2, err := GetPoolKey("cdash", "rwn")
	assert.NoError(t, err)
	assert.NotEqual(t, poolKey, poolKey2, "Generated keys must be unique")
	poolAddress, err := sdk.AccAddressFromHex(hex.EncodeToString(poolKey))
	assert.NoError(t, err, "Address should be convertible to cosmos sdk address because of padding")
	poolAddress2, err := sdk.AccAddressFromHex(hex.EncodeToString(poolKey2))
	assert.NoError(t, err, "Address should be convertible to cosmos sdk address because of padding")
	assert.NotEqual(t, poolAddress, poolAddress2, "Generated addresses must be unique")
	assert.IsType(t, sdk.AccAddress{}, poolAddress)
	assert.IsType(t, sdk.AccAddress{}, poolAddress2)
}

func TestGetLiquidityProviderKey(t *testing.T) {
	paddedbytes, err := pkcs7Pad([]byte("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v"), 20)
	assert.NoError(t, err)
	hx := hex.EncodeToString(paddedbytes)
	lpaddress, err := sdk.AccAddressFromHex(hx)
	assert.NoError(t, err)
	lpKey := GetLiquidityProviderKey("ceth", lpaddress.String())
	lpKey2 := GetLiquidityProviderKey("cdash", lpaddress.String())
	assert.NotEqual(t, lpKey, lpKey2, "Generated keys must be different for different pools")
}
