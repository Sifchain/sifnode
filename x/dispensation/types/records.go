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

func NewDistributionRecord(distributionName string, recipientAddress string, coins sdk.Coins, start sdk.Int, end sdk.Int) DistributionRecord {
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
