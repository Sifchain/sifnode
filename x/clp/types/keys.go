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
	RewardParamPrefix        = []byte{0x06}
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

func DefaultRewardsPeriod() []*RewardPeriod {
	rp_1_allocation := sdk.NewUintFromString("10000000000000000000000")
	cethMultiplier := sdk.MustNewDecFromStr("1.5")
	rewardPeriods := []*RewardPeriod{
		{
			Id:         "RP_1",
			StartBlock: 1,
			EndBlock:   12 * 60 * 24 * 7,
			Allocation: &rp_1_allocation,
			Multipliers: []*PoolMultiplier{{
				Asset:      "ceth",
				Multiplier: &cethMultiplier,
			}},
		},
	}
	return rewardPeriods
}
func GetDefaultRewardParams() *RewardParams {
	zero := sdk.ZeroDec()
	return &RewardParams{
		LiquidityRemovalLockPeriod:   12 * 60 * 24 * 7,
		LiquidityRemovalCancelPeriod: 12 * 60 * 24 * 30,
		DefaultMultiplier:            &zero,
		RewardPeriods:                DefaultRewardsPeriod(),
	}
}

func GetDefaultPmtpParams() *PmtpParams {
	return &PmtpParams{
		PmtpPeriodGovernanceRate: sdk.MustNewDecFromStr("0.0010"),
		PmtpPeriodEpochLength:    14440,
		PmtpPeriodStartBlock:     211,
		PmtpPeriodEndBlock:       72210,
	}
}
