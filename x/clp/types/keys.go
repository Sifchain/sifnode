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
)

var (
	PoolPrefix              = []byte{0x00} // key for storing Pools
	LiquidityProviderPrefix = []byte{0x01} // key for storing Liquidity Providers
)

func GetPoolKey(ticker string, native string) ([]byte, error) {
	addr, err := GetPoolAddress(ticker, native)
	if err != nil {
		return nil, err
	}
	key := []byte(addr)
	return append(PoolPrefix, key...), nil
}

func GetPoolAddress(ticker string, native string) (string, error) {
	addr, err := GetAddress(ticker, native)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func GetAddress(ticker, native string) (sdk.AccAddress, error) {
	addressBytes := []byte(fmt.Sprintf("%s_%s", ticker, native))
	paddedbytes, err := pkcs7Pad(addressBytes, 20)
	if err != nil {
		return nil, err
	}
	hx := hex.EncodeToString(paddedbytes)
	return sdk.AccAddressFromHex(hx)
}

func GetLiquidityProviderKey(ticker string, lp string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", ticker, lp))
	return append(LiquidityProviderPrefix, key...)
}

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
