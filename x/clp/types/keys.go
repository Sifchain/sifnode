package types

import (
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

	NativeSymbol = "rowan"
	PoolThrehold = "1000000000000000000"

	MaxSymbolLength = 71
	MaxWbasis       = 10000
)

var (
	PoolPrefix               = []byte{0x00} // key for storing Pools
	LiquidityProviderPrefix  = []byte{0x01} // key for storing Liquidity Providers
	WhiteListValidatorPrefix = []byte{0x02} // Key to store WhiteList , allowed to decommission pools
	PmtpRateParamsPrefix     = []byte{0x03} // Key to store the Pmtp rate params
	PmtpEpochPrefix          = []byte{0x04} // Key to store the Epoch
	PmtpParamsPrefix         = []byte{0x05} // Key to store the Pmtp params
	AdminAccountStorePrefix  = []byte{0x06} // Key to store the admin account
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

func GetDefaultPmtpParams() *PmtpParams {
	//return &PmtpParams{
	//	PmtpPeriodGovernanceRate: sdk.MustNewDecFromStr("0.10"),
	//	PmtpPeriodEpochLength:    14440,
	//	PmtpPeriodStartBlock:     211,
	//	PmtpPeriodEndBlock:       72210,
	//}
	return &PmtpParams{
		PmtpPeriodGovernanceRate: sdk.MustNewDecFromStr("0.10"),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     1,
		PmtpPeriodEndBlock:       1000,
	}
}
