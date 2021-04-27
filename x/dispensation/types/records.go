package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

//This package is used to keep historical data. This will later be used to distribute rewards over different blocks through a gov proposal

type ClaimStatus int64

const Pending ClaimStatus = 1
const Completed ClaimStatus = 2

func (d ClaimStatus) String() string {
	switch d {
	case Pending:
		return "Pending"
	case Completed:
		return "Completed"
	default:
		return "All"
	}
}

// DistributionRecord is created for every recipient for a distribution
// TODO : Add ClaimStatus to the prefixed key for records
type DistributionRecord struct {
	ClaimStatus
	DistributionName            string         `json:"distribution_name"`
	RecipientAddress            sdk.AccAddress `json:"recipient_address"`
	Coins                       sdk.Coins      `json:"coins"`
	DistributionStartHeight     int64          `json:"distribution_start_height"`
	DistributionCompletedHeight int64          `json:"distribution_completed_height"`
}
type DistributionRecords []DistributionRecord

func NewDistributionRecord(distributionName string, recipientAddress sdk.AccAddress, coins sdk.Coins, start int64, end int64) DistributionRecord {
	return DistributionRecord{DistributionName: distributionName, RecipientAddress: recipientAddress, Coins: coins, DistributionStartHeight: start, DistributionCompletedHeight: end}
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

func (dr DistributionRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(`DistributionName: %s
	RecipientAddress: %s
	Coins: %s
    ClaimStatus :%s
    DistributionStartHeight :%d 
    DistributionCompletedHeight :%d `, dr.DistributionName, dr.RecipientAddress.String(), dr.Coins.String(), dr.ClaimStatus.String(), dr.DistributionStartHeight, dr.DistributionCompletedHeight))
}

func (dr DistributionRecord) Add(dr2 DistributionRecord) DistributionRecord {
	dr.Coins = dr.Coins.Add(dr2.Coins...)
	return dr
}

type DistributionType int64

const Airdrop DistributionType = 1

func (d DistributionType) String() string {
	switch d {
	case Airdrop:
		return "Airdrop"
	default:
		return "Invalid"
	}
}

// A Distribution object is created for every new distribution type
type Distribution struct {
	DistributionType DistributionType `json:"distribution_type"`
	DistributionName string           `json:"distribution_name"`
}
type Distributions []Distribution

func NewDistribution(t DistributionType, name string) Distribution {
	return Distribution{DistributionType: t, DistributionName: name}
}

func (ar Distribution) Validate() bool {
	if ar.DistributionName == "" {
		return false
	}
	return true
}

func (ar Distribution) String() string {
	return strings.TrimSpace(fmt.Sprintf(`DistributionName: %s DistributionType :%s`, ar.DistributionName, ar.DistributionType.String()))
}

type UserClaim struct {
	UserAddress sdk.AccAddress `json:"user_address"`
}

func NewUserClaim(userAddress sdk.AccAddress) UserClaim {
	return UserClaim{UserAddress: userAddress}
}

func (uc UserClaim) Validate() bool {
	if uc.UserAddress.Empty() {
		return false
	}
	return true
}

func (uc UserClaim) String() string {
	return strings.TrimSpace(fmt.Sprintf(`UserAddress : %s`, uc.UserAddress.String()))
}
