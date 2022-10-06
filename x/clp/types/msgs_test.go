package types

import (
	"bytes"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
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

func GetETHAsset() Asset {
	return NewAsset("eth")
}

func GetROWANAsset() Asset {
	return NewAsset("rowan")
}

func GetWrongAsset() Asset {
	return NewAsset("01234567890123456789012345678901234567890123456789012345678901234567890123456789")
}

func TestNewMsgCreatePool(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := GetETHAsset()
	newpool := NewMsgCreatePool(signer, asset, sdk.NewUint(1000), sdk.NewUint(100))
	err := newpool.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, newpool.GetSigners()[0], signer)
	wrongAsset := GetWrongAsset()
	newpool = NewMsgCreatePool(signer, wrongAsset, sdk.NewUint(1000), sdk.NewUint(100))
	err = newpool.ValidateBasic()
	assert.Error(t, err)

	str := newpool.Route()
	assert.Equal(t, str, "clp")
	str = newpool.Type()
	assert.Equal(t, str, "create_pool")
	newpool = NewMsgCreatePool(nil, asset, sdk.NewUint(1000), sdk.NewUint(100))
	err = newpool.ValidateBasic()
	assert.Error(t, err, "invalid address")
	newpool = NewMsgCreatePool(signer, GetROWANAsset(), sdk.NewUint(1000), sdk.NewUint(100))
	err = newpool.ValidateBasic()
	assert.Error(t, err, "External Asset cannot be rowan")
	newpool = NewMsgCreatePool(signer, asset, sdk.NewUint(0), sdk.NewUint(100))
	err = newpool.ValidateBasic()
	assert.Error(t, err, "amount is invalid")
	newpool = NewMsgCreatePool(signer, asset, sdk.NewUint(1000), sdk.NewUint(0))
	err = newpool.ValidateBasic()
	assert.Error(t, err, "amount is invalid")
}

func TestNewMsgDecommissionPool(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := GetETHAsset()
	tx := NewMsgDecommissionPool(signer, asset.Symbol)
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := GetWrongAsset()
	tx = NewMsgDecommissionPool(signer, wrongAsset.Symbol)
	err = tx.ValidateBasic()
	assert.Error(t, err)

	str := tx.Route()
	assert.Equal(t, str, "clp")
	str = tx.Type()
	assert.Equal(t, str, "decommission_pool")
	tx = NewMsgDecommissionPool(nil, asset.Symbol)
	err = tx.ValidateBasic()
	assert.Error(t, err, "invalid address")
}

func TestNewMsgSwap(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := GetETHAsset()
	tx := NewMsgSwap(signer, asset, GetSettlementAsset(), sdk.NewUint(100), sdk.NewUint(90))
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := GetWrongAsset()
	tx = NewMsgSwap(signer, wrongAsset, GetSettlementAsset(), sdk.NewUint(100), sdk.NewUint(90))
	err = tx.ValidateBasic()
	assert.Error(t, err)
	tx = NewMsgSwap(signer, asset, GetSettlementAsset(), sdk.NewUint(0), sdk.NewUint(90))
	err = tx.ValidateBasic()
	assert.Error(t, err)

	str := tx.Route()
	assert.Equal(t, str, "clp")
	str = tx.Type()
	assert.Equal(t, str, "swap")
	tx = NewMsgSwap(nil, asset, GetSettlementAsset(), sdk.NewUint(100), sdk.NewUint(90))
	err = tx.ValidateBasic()
	assert.Error(t, err, "invalid address")
	tx = NewMsgSwap(signer, asset, wrongAsset, sdk.NewUint(100), sdk.NewUint(90))
	err = tx.ValidateBasic()
	assert.Error(t, err, "asset is invalid")
	tx = NewMsgSwap(signer, asset, asset, sdk.NewUint(100), sdk.NewUint(100))
	err = tx.ValidateBasic()
	assert.Error(t, err, "Sent And Received asset cannot be the same")
	tx = NewMsgSwap(signer, asset, GetSettlementAsset(), sdk.NewUint(0), sdk.NewUint(100))
	err = tx.ValidateBasic()
	assert.Error(t, err, "amount is invalid")
}

func TestNewMsgAddLiquidity(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := GetETHAsset()
	tx := NewMsgAddLiquidity(signer, asset, sdk.NewUint(100), sdk.NewUint(100))
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := GetWrongAsset()
	tx = NewMsgAddLiquidity(signer, wrongAsset, sdk.NewUint(100), sdk.NewUint(100))
	err = tx.ValidateBasic()
	assert.Error(t, err)

	str := tx.Route()
	assert.Equal(t, str, "clp")
	str = tx.Type()
	assert.Equal(t, str, "add_liquidity")
	tx = NewMsgAddLiquidity(nil, asset, sdk.NewUint(100), sdk.NewUint(100))
	err = tx.ValidateBasic()
	assert.Error(t, err, "invalid address")
	tx = NewMsgAddLiquidity(signer, GetSettlementAsset(), sdk.NewUint(100), sdk.NewUint(100))
	err = tx.ValidateBasic()
	assert.Error(t, err, "External asset cannot be rowan")
	tx = NewMsgAddLiquidity(signer, asset, sdk.ZeroUint(), sdk.ZeroUint())
	err = tx.ValidateBasic()
	assert.Error(t, err, "Both asset ammounts cannot be 0 0 / 0")
}

func TestNewMsgRemoveLiquidity(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	asset := GetETHAsset()
	tx := NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(100), sdk.NewInt(100))
	err := tx.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, tx.GetSigners()[0], signer)
	wrongAsset := GetWrongAsset()
	tx = NewMsgRemoveLiquidity(signer, wrongAsset, sdk.NewInt(100), sdk.NewInt(100))
	err = tx.ValidateBasic()
	assert.Error(t, err)
	tx = NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(-100), sdk.NewInt(100))
	err = tx.ValidateBasic()
	assert.Error(t, err)
	tx = NewMsgRemoveLiquidity(signer, asset, sdk.NewInt(100), sdk.NewInt(100000))
	err = tx.ValidateBasic()
	assert.Error(t, err)

	str := tx.Route()
	assert.Equal(t, str, "clp")
	str = tx.Type()
	assert.Equal(t, str, "remove_liquidity")
	tx = NewMsgRemoveLiquidity(nil, asset, sdk.NewInt(100), sdk.NewInt(100))
	err = tx.ValidateBasic()
	assert.Error(t, err, "invalid address")
}

func TestNewMsgAddProviderDistributionPeriodRequest(t *testing.T) {
	signer := NewSigner("A58856F0FD53BF058B4909A21AEC019107BA6")
	var periods []*ProviderDistributionPeriod

	validPeriod := ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 10, DistributionPeriodBlockRate: sdk.NewDecWithPrec(1, 2), DistributionPeriodMod: 1}
	startBeforeEnd := ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 8, DistributionPeriodBlockRate: sdk.NewDecWithPrec(1, 2), DistributionPeriodMod: 1}
	rateTooLow := ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 12, DistributionPeriodBlockRate: sdk.NewDec(-1), DistributionPeriodMod: 1}
	rateTooHigh := ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 12, DistributionPeriodBlockRate: sdk.NewDec(2), DistributionPeriodMod: 1}
	moduloTooLow := ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 12, DistributionPeriodBlockRate: sdk.NewDecWithPrec(1, 2), DistributionPeriodMod: 0}

	periods = append(periods, &validPeriod)
	tx := MsgAddProviderDistributionPeriodRequest{Signer: signer.String(), DistributionPeriods: periods}
	err := tx.ValidateBasic()
	assert.NoError(t, err)

	periods = append(periods, &startBeforeEnd)
	tx = MsgAddProviderDistributionPeriodRequest{Signer: signer.String(), DistributionPeriods: periods}
	err = tx.ValidateBasic()
	assert.Error(t, err)

	periods[1] = &rateTooLow
	tx = MsgAddProviderDistributionPeriodRequest{Signer: signer.String(), DistributionPeriods: periods}
	err = tx.ValidateBasic()
	assert.Error(t, err)

	periods[1] = &rateTooHigh
	tx = MsgAddProviderDistributionPeriodRequest{Signer: signer.String(), DistributionPeriods: periods}
	err = tx.ValidateBasic()
	assert.Error(t, err)

	periods[1] = &moduloTooLow
	tx = MsgAddProviderDistributionPeriodRequest{Signer: signer.String(), DistributionPeriods: periods}
	err = tx.ValidateBasic()
	assert.Error(t, err)
}
