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

	NativeSymbol        = "rowan"
	PoolThrehold        = "1000000000000000000"
	NativeAssetDecimals = 18

	MaxSymbolLength = 71
	MaxWbasis       = 10000
)

var (
	PoolPrefix                          = []byte{0x00} // key for storing Pools
	LiquidityProviderPrefix             = []byte{0x01} // key for storing Liquidity Providers
	WhiteListValidatorPrefix            = []byte{0x02} // Key to store WhiteList , allowed to decommission pools
	PmtpRateParamsPrefix                = []byte{0x03} // Key to store the Pmtp rate params
	PmtpEpochPrefix                     = []byte{0x04} // Key to store the Epoch
	PmtpParamsPrefix                    = []byte{0x05} // Key to store the Pmtp params
	RewardParamPrefix                   = []byte{0x06}
	SymmetryThresholdPrefix             = []byte{0x07}
	LiquidityProtectionParamsPrefix     = []byte{0x08} // Key to store the Liquidity Protection params
	LiquidityProtectionRateParamsPrefix = []byte{0x09} // Key to store the Liquidity Protection rate params
	ProviderDistributionParamsPrefix    = []byte{0x0a}
	RewardsBlockDistributionPrefix      = []byte{0x0b}
	SwapFeeRatePrefix                   = []byte{0x0c}
	RemovalRequestPrefix                = []byte{0x0d}
	RemovalQueuePrefix                  = []byte{0x0e}
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

func GetDefaultRewardParams() *RewardParams {
	return &RewardParams{
		LiquidityRemovalLockPeriod:   0,
		LiquidityRemovalCancelPeriod: 12 * 60 * 24 * 30,
		RewardPeriods:                nil,
		RewardPeriodStartTime:        "",
	}
}

func GetDefaultPmtpParams() *PmtpParams {
	return &PmtpParams{
		PmtpPeriodGovernanceRate: sdk.MustNewDecFromStr("0.0"),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     0,
		PmtpPeriodEndBlock:       0,
	}
}

func GetDefaultLiquidityProtectionParams() *LiquidityProtectionParams {
	return &LiquidityProtectionParams{
		MaxRowanLiquidityThreshold:      sdk.NewUint(1000000000000),
		MaxRowanLiquidityThresholdAsset: "cusdc",
		EpochLength:                     14400,
		IsActive:                        false,
	}
}

// GetRemovalRequestKey generates a key to store a removal request,
// the key is in the format: lpaddress_id
func GetRemovalRequestKey(request RemovalRequest) []byte {
	key := []byte(fmt.Sprintf("%s_%d", request.Msg.Signer, request.Id))
	return append(RemovalRequestPrefix, key...)
}

func GetRemovalRequestLPPrefix(lpaddress string) []byte {
	key := []byte(fmt.Sprintf("%s", lpaddress))
	return append(RemovalRequestPrefix, key...)
}

func GetRemovalQueueKey(symbol string) []byte {
	key := []byte(fmt.Sprintf("_%s", symbol))
	return append(RemovalQueuePrefix, key...)
}
