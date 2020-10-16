package common

import (
	"github.com/Sifchain/sifnode/x/clp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCoin(t *testing.T) {
	coin := NewCoin(clp.GetNativeAsset(), sdk.NewUint(230000000))
	assert.Equal(t, coin.Asset, clp.GetNativeAsset(), "")
	assert.Equal(t, coin.Amount, sdk.NewUint(230000000))
}
func TestCoin_IsEmpty(t *testing.T) {
	coin := NewCoin(clp.GetNativeAsset(), sdk.NewUint(230000000))
	assert.False(t, coin.IsEmpty(), "")
}

func TestCoin_IsNative(t *testing.T) {
	coin := NewCoin(clp.GetNativeAsset(), sdk.NewUint(230000000))
	assert.True(t, coin.IsNative(), "")
}
