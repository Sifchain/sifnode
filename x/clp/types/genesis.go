package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	admin := GetDefaultCLPAdmin()
	return &GenesisState{
		Params:                        DefaultParams(),
		AddressWhitelist:              []string{admin.String()},
		PoolList:                      []*Pool{},
		LiquidityProviders:            []*LiquidityProvider{},
		RewardsBucketList:             []RewardsBucket{},
		RewardParams:                  *GetDefaultRewardParams(),
		PmtpParams:                    *GetDefaultPmtpParams(),
		PmtpEpoch:                     PmtpEpoch{},
		PmtpRateParams:                PmtpRateParams{},
		LiquidityProtectionParams:     *GetDefaultLiquidityProtectionParams(),
		LiquidityProtectionRateParams: LiquidityProtectionRateParams{},
		SwapFeeParams:                 *GetDefaultSwapFeeParams(),
		ProviderDistributionParams:    *GetDefaultProviderDistributionParams(),
	}
}

func GetGenesisStateFromAppState(marshaler codec.JSONCodec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := marshaler.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	return genesisState
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in rewardsBucket
	rewardsBucketIndexMap := make(map[string]struct{})

	for _, elem := range gs.RewardsBucketList {
		index := string(RewardsBucketKey(elem.Denom))
		if _, ok := rewardsBucketIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for rewardsBucket")
		}
		rewardsBucketIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
