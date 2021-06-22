package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//This package is used to keep historical data. This will later be used to distribute rewards over different blocks through a gov proposal

func NewDistributionRecord(status DistributionStatus, distributionType DistributionType, distributionName string, recipientAddress string, coins sdk.Coins, start int64, end int64, runner string) DistributionRecord {
	return DistributionRecord{
		DistributionStatus:          status,
		AuthorizedRunner:            runner,
		DistributionType:            distributionType,
		DistributionName:            distributionName,
		RecipientAddress:            recipientAddress,
		Coins:                       coins,
		DistributionStartHeight:     start,
		DistributionCompletedHeight: end,
	}

}

func (dr DistributionRecord) Validate() bool {
	if len(dr.RecipientAddress) == 0 {
		return false
	}
	if !dr.Coins.IsValid() {
		return false
	}
	if !dr.Coins.IsAllPositive() {
		return false
	}
	return true
}

func (dr DistributionRecord) Add(dr2 DistributionRecord) DistributionRecord {
	dr.Coins = dr.Coins.Add(dr2.Coins...)
	return dr
}

func (dr DistributionRecord) DoesTypeSupportClaim() bool {
	if dr.DistributionType == DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING || dr.DistributionType == DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY {
		return true
	}
	return false
}

func NewUserClaim(userAddress string, userClaimType DistributionType, time string) UserClaim {
	return UserClaim{UserAddress: userAddress, UserClaimType: userClaimType, UserClaimTime: time}
}

func (uc UserClaim) Validate() bool {
	if len(uc.UserAddress) == 0 {
		return false
	}
	return true
}

func NewDistribution(t DistributionType, name string) Distribution {
	return Distribution{DistributionType: t, DistributionName: name}
}

func (ar Distribution) Validate() bool {
	if ar.DistributionName == "" {
		return false
	}
	return true
}
func GetDistributionStatus(status string) (DistributionStatus, bool) {
	switch status {
	case "Completed":
		return DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, true
	case "Pending":
		return DistributionStatus_DISTRIBUTION_STATUS_PENDING, true
	case "Failed":
		return DistributionStatus_DISTRIBUTION_STATUS_FAILED, true
	default:
		return DistributionStatus_DISTRIBUTION_STATUS_UNSPECIFIED, false
	}
}

func GetClaimType(claimType string) (DistributionType, bool) {
	switch claimType {
	case "LiquidityMining":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "ValidatorSubsidy":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
}

func GetDistributionTypeFromShortString(distributionType string) (DistributionType, bool) {
	switch distributionType {
	case "Airdrop":
		return DistributionType_DISTRIBUTION_TYPE_AIRDROP, true
	case "LiquidityMining":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "ValidatorSubsidy":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
}
func IsValidDistributionType(distributionType string) (DistributionType, bool) {
	switch distributionType {
	case "DISTRIBUTION_TYPE_AIRDROP":
		return DistributionType_DISTRIBUTION_TYPE_AIRDROP, true
	case "DISTRIBUTION_TYPE_LIQUIDITY_MINING":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, false
	}
}

func IsValidClaimType(claimType string) (DistributionType, bool) {
	switch claimType {
	case "DISTRIBUTION_TYPE_LIQUIDITY_MINING":
		return DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, true
	case "DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY":
		return DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, true
	default:
		return 0, false
	}
}
