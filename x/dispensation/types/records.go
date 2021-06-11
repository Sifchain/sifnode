package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

//This package is used to keep historical data. This will later be used to distribute rewards over different blocks through a gov proposal

type DistributionStatus int64

const Pending DistributionStatus = 1
const Completed DistributionStatus = 2

func (ds DistributionStatus) String() string {
	switch ds {
	case Pending:
		return "Pending"
	case Completed:
		return "Completed"
	default:
		return "All"
	}
}

func IsValidStatus(s string) (DistributionStatus, bool) {
	switch s {
	case "Pending":
		return Pending, true
	case "Completed":
		return Completed, true
	default:
		return -1, false
	}
}

// DistributionRecord is created for every recipient for a distribution
// TODO : Add DistributionStatus to the prefixed key for records
type DistributionRecord struct {
	DistributionStatus          DistributionStatus `json:"distribution_status"`
	DistributionName            string             `json:"distribution_name"`
	DistributionType            DistributionType   `json:"distribution_type"`
	RecipientAddress            sdk.AccAddress     `json:"recipient_address"`
	Coins                       sdk.Coins          `json:"coins"`
	DistributionStartHeight     int64              `json:"distribution_start_height"`
	DistributionCompletedHeight int64              `json:"distribution_completed_height"`
	AuthorizedRunner            sdk.AccAddress     `json:"authorized_runner"`
}

type DistributionRecords []DistributionRecord

func (records DistributionRecords) String() string {
	if len(records) == 0 {
		return ""
	}
	var rc string
	for _, record := range records {
		rc = rc + record.RecipientAddress.String() + ","
	}
	rc = rc[:len(rc)-1]
	return rc
}

func NewDistributionRecord(distributionName string, distributionType DistributionType, recipientAddress sdk.AccAddress, coins sdk.Coins, start int64, end int64, authorizedRunner sdk.AccAddress) DistributionRecord {
	return DistributionRecord{
		DistributionName:            distributionName,
		DistributionType:            distributionType,
		RecipientAddress:            recipientAddress,
		Coins:                       coins,
		DistributionStartHeight:     start,
		DistributionCompletedHeight: end,
		AuthorizedRunner:            authorizedRunner}
}

func (dr DistributionRecord) DoesClaimExist() bool {
	if dr.DistributionType == LiquidityMining || dr.DistributionType == ValidatorSubsidy {
		return true
	}
	return false
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
	_, err := sdk.AccAddressFromBech32(dr.AuthorizedRunner.String())
	if err != nil {
		return false
	}
	_, err = sdk.AccAddressFromBech32(dr.RecipientAddress.String())
	if err != nil {
		return false
	}
	_, ok := IsValidDistributionType(dr.DistributionType.String())
	if !ok {
		return false
	}
	return true
}

func (dr DistributionRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(`DistributionName: %s
    DistributionType : %s, 
	RecipientAddress: %s
	Coins: %s
    DistributionStatus :%s
    DistributionStartHeight :%d 
    DistributionCompletedHeight :%d `, dr.DistributionName, dr.DistributionType, dr.RecipientAddress.String(), dr.Coins.String(), dr.DistributionStatus.String(), dr.DistributionStartHeight, dr.DistributionCompletedHeight))
}

func (dr DistributionRecord) Add(dr2 DistributionRecord) DistributionRecord {
	dr.Coins = dr.Coins.Add(dr2.Coins...)
	return dr
}

// The same type is also used Claims
type DistributionType int64

const DistributionTypeUnknown DistributionType = 0
const Airdrop DistributionType = 1
const LiquidityMining DistributionType = 2
const ValidatorSubsidy DistributionType = 3

func (dt DistributionType) String() string {
	switch dt {
	case Airdrop:
		return "Airdrop"
	case LiquidityMining:
		return "LiquidityMining"
	case ValidatorSubsidy:
		return "ValidatorSubsidy"
	default:
		return "Invalid"
	}
}

func IsValidDistributionType(distributionType string) (DistributionType, bool) {
	switch distributionType {
	case "Airdrop":
		return Airdrop, true
	case "LiquidityMining":
		return LiquidityMining, true
	case "ValidatorSubsidy":
		return ValidatorSubsidy, true
	default:
		return 0, false
	}
}

func IsValidClaim(claimType string) (DistributionType, bool) {
	switch claimType {
	case "LiquidityMining":
		return LiquidityMining, true
	case "ValidatorSubsidy":
		return ValidatorSubsidy, true
	default:
		return 0, false
	}
}

// A Distribution object is created for every new distribution type
type Distribution struct {
	DistributionType DistributionType `json:"distribution_type"`
	DistributionName string           `json:"distribution_name"`
	Runner           sdk.AccAddress   `json:"runner"`
}

type Distributions []Distribution

func NewDistribution(t DistributionType, name string, runner sdk.AccAddress) Distribution {
	return Distribution{DistributionType: t, DistributionName: name, Runner: runner}
}

func (d Distribution) Validate() bool {
	if d.DistributionName == "" {
		return false
	}
	_, ok := IsValidDistributionType(d.DistributionType.String())
	if !ok {
		return false
	}
	_, err := sdk.AccAddressFromBech32(d.Runner.String())
	if err != nil {
		return false
	}

	return true
}

func (d Distribution) String() string {
	return strings.TrimSpace(fmt.Sprintf(`DistributionName: %s DistributionType :%s`, d.DistributionName, d.DistributionType.String()))
}

type UserClaim struct {
	UserAddress   sdk.AccAddress   `json:"user_address"`
	UserClaimType DistributionType `json:"user_claim_type"`
	UserClaimTime time.Time        `json:"user_claim_time"`
}

type UserClaims []UserClaim

func NewUserClaim(userAddress sdk.AccAddress, userClaimType DistributionType, time time.Time) UserClaim {
	return UserClaim{UserAddress: userAddress, UserClaimType: userClaimType, UserClaimTime: time}
}

func (uc UserClaim) Validate() bool {
	if uc.UserAddress.Empty() {
		return false
	}
	_, err := sdk.AccAddressFromBech32(uc.UserAddress.String())
	if err != nil {
		return false
	}
	_, ok := IsValidClaim(uc.UserClaimType.String())
	if !ok {
		return false
	}
	return true
}

func (uc UserClaim) String() string {
	return strings.TrimSpace(fmt.Sprintf(`UserAddress : %s | UserClaimType : %s`, uc.UserAddress.String(), uc.UserClaimType))
}
