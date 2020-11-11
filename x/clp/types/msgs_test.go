package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func NewSigner(signer string) sdk.AccAddress {
	var buffer bytes.Buffer
	buffer.WriteString(signer)
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

func TestNewMsgCreatePool(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := NewAsset("Ethereum", "ETH", "ceth")
	newpool := NewMsgCreatePool(signer, asset, sdk.NewUint(1000), sdk.NewUint(100))
	err := newpool.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, newpool.GetSigners()[0], signer)
	wrongAsset := NewAsset("Ethereum", "ETH", "deadbeef")
	newpool = NewMsgCreatePool(signer, wrongAsset, sdk.NewUint(1000), sdk.NewUint(100))
	err = newpool.ValidateBasic()
	assert.Error(t, err)
}

func TestNewMsgDecommissionPool(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := NewAsset("Ethereum", "ETH", "ceth")
	tx := NewMsgDecommissionPool(signer, asset.Ticker)
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := NewAsset("Ethereum", "ETH", "deadbeef")
	tx = NewMsgDecommissionPool(signer, wrongAsset.Ticker)
	err = tx.ValidateBasic()
	assert.Error(t, err)
}

func TestNewMsgSwap(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := NewAsset("Ethereum", "ETH", "ceth")
	tx := NewMsgSwap(signer, asset, GetSettlementAsset(), sdk.NewUint(100))
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := NewAsset("Ethereum", "ETH", "deadbeef")
	tx = NewMsgSwap(signer, wrongAsset, GetSettlementAsset(), sdk.NewUint(100))
	err = tx.ValidateBasic()
	assert.Error(t, err)
	tx = NewMsgSwap(signer, asset, GetSettlementAsset(), sdk.NewUint(0))
	err = tx.ValidateBasic()
	assert.Error(t, err)
}

func TestNewMsgAddLiquidity(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := NewAsset("Ethereum", "ETH", "ceth")
	tx := NewMsgAddLiquidity(signer, asset, sdk.NewUint(100), sdk.NewUint(100))
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := NewAsset("Ethereum", "ETH", "deadbeef")
	tx = NewMsgAddLiquidity(signer, wrongAsset, sdk.NewUint(100), sdk.NewUint(100))
	err = tx.ValidateBasic()
	assert.Error(t, err)
}

func TestNewMsgRemoveLiquidity(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := NewAsset("Ethereum", "ETH", "ceth")
	tx := NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(100), sdk.NewInt(100))
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := NewAsset("Ethereum", "ETH", "deadbeef")
	tx = NewMsgRemoveLiquidity(signer, wrongAsset, sdk.NewInt(100), sdk.NewInt(100))
	err = tx.ValidateBasic()
	assert.Error(t, err)
	tx = NewMsgRemoveLiquidity(signer, wrongAsset, sdk.NewInt(-100), sdk.NewInt(100))
	err = tx.ValidateBasic()
	assert.Error(t, err)
	tx = NewMsgRemoveLiquidity(signer, wrongAsset, sdk.NewInt(100), sdk.NewInt(100000))
	err = tx.ValidateBasic()
	assert.Error(t, err)
}
