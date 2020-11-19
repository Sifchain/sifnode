package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "clp"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName

	NativeTicker = "rwn"
	NativeChain  = "SIFCHAIN"
	NativeSymbol = "RWN"

	AddressLength        = 20
	MaxTickerLength      = 6
	MaxSymbolLength      = 6
	MaxSourceChainLength = 20
	MaxWbasis            = 10000
)

var (
	PoolPrefix              = []byte{0x00} // key for storing Pools
	LiquidityProviderPrefix = []byte{0x01} // key for storing Liquidity Providers
)

func GetPoolKey(externalTicker string, nativeTicker string) ([]byte, error) {
	addr, err := GetPoolAddress(externalTicker, nativeTicker)
	if err != nil {
		return nil, err
	}
	key := []byte(addr)
	return append(PoolPrefix, key...), nil
}

//Generate a new pool address from a string
//The external asset ticker and the native asset ticket ,in combination is used to generate an unique address
func GetPoolAddress(externalTicker string, nativeTicker string) (string, error) {
	addr, err := GetAddress(externalTicker, nativeTicker)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func GetAddress(externalTicker, nativeTicker string) (sdk.AccAddress, error) {
	addressBytes := []byte(fmt.Sprintf("%s_%s", externalTicker, nativeTicker))
	paddedbytes, err := pkcs7Pad(addressBytes, AddressLength)
	if err != nil {
		return nil, err
	}
	hx := hex.EncodeToString(paddedbytes)
	return sdk.AccAddressFromHex(hx)
}

// Generate key to store a Liquidity Provider
// The key is of the format ticker_lpaddress
// Example : eth_sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v
func GetLiquidityProviderKey(externalTicker string, lp string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", externalTicker, lp))
	return append(LiquidityProviderPrefix, key...)
}

// Padding extra bytes to meet the size requirments of the cosmos address variable
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}
