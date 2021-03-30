package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

//This package is used to keep historical data. This will later be used to distribute rewards over different blocks through a gov proposal

// DistributionRecord is created for every recipient for a distribution
// TODO add a claim status for the distribution record which can be used to break the Distribution into two different processes . Distribute and Claim
type DistributionRecord struct {
	DistributionName string         `json:"airdrop_name"`
	Address          sdk.AccAddress `json:"address" yaml:"address"`
	Coins            sdk.Coins      `json:"coins" yaml:"coins"`
}
type DistributionRecords []DistributionRecord

func NewDistributionRecord(name string, address sdk.AccAddress, coins sdk.Coins) DistributionRecord {
	return DistributionRecord{DistributionName: name, Address: address, Coins: coins}
}

func (dr DistributionRecord) Validate() bool {
	if len(dr.Address) == 0 {
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
	Address: %s
	Coins: %s`, dr.DistributionName, dr.Address.String(), dr.Coins.String()))
}

func (ar DistributionRecord) Add(ar2 DistributionRecord) DistributionRecord {
	ar.Coins = ar.Coins.Add(ar2.Coins...)
	return ar
}

// A Distribution object is created for every new distribution type
type DistributionType int64

const Airdrop DistributionType = 1

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
	return strings.TrimSpace(fmt.Sprintf(`DistributionName: %s`, ar.DistributionName))
}
