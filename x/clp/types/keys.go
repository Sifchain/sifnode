package types

import (
	"fmt"
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

	NativeSymbol      = "rowan"
	SwapType          = "swap"
	PoolThrehold      = "1000000000000000000"
	PoolUnitsMinValue = "1000000000"
	TxFeeMultiplier   = "1"

	MaxSymbolLength   = 10
	MaxWbasis         = 10000
	MaxTokenPrecision = 18
	MinTokenPrecision = 6
)

var (
	PoolPrefix               = []byte{0x00} // key for storing Pools
	LiquidityProviderPrefix  = []byte{0x01} // key for storing Liquidity Providers
	WhiteListValidatorPrefix = []byte{0x02} // Key to store WhiteList , allowed to decommission pools
)

// Generates a key for storing a specific pool
// The key is of the format externalticker_nativeticker
// Example : eth_rwn and converted into bytes after adding a prefix
func GetPoolKey(externalTicker string, nativeTicker string) ([]byte, error) {
	key := []byte(fmt.Sprintf("%s_%s", externalTicker, nativeTicker))
	return append(PoolPrefix, key...), nil
}

// Generate key to store a Liquidity Provider
// The key is of the format ticker_lpaddress
// Example : eth_sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v and converted into bytes after adding a prefix
func GetLiquidityProviderKey(externalTicker string, lp string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", externalTicker, lp))
	return append(LiquidityProviderPrefix, key...)
}

func GetNormalizationMap() map[string]int64 {
	m := make(map[string]int64)
	m["cel"] = 4
	m["ausdc"] = 6
	m["usdt"] = 6
	m["usdc"] = 6
	m["cro"] = 8
	m["cdai"] = 8
	m["wbtc"] = 8
	m["ceth"] = 8
	m["renbtc"] = 8
	m["cusdc"] = 8
	m["husd"] = 8
	m["ampl"] = 9
	return m
}
