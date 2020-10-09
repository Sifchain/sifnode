package types

import "fmt"

const (
	// ModuleName is the name of the module
	ModuleName = "clp"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName

	NativeToken = "ROWAN"
)

var (
	PoolPrefix              = []byte{0x00} // key for storing Pools
	LiquidityProviderPrefix = []byte{0x01} // key for storing Liquidity Providers
)

func GetPoolKey(ticker string, native string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", ticker, native))
	return append(PoolPrefix, key...)
}

func GetLiquidityProviderKey(ticker string, lp string) []byte {
	key := []byte(fmt.Sprintf("%s_%s", ticker, lp))
	return append(LiquidityProviderPrefix, key...)
}
