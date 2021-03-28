package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type DistributionRecord struct {
	DistributionName string         `json:"airdrop_name"`
	Address          sdk.AccAddress `json:"address" yaml:"address"`
	Coins            sdk.Coins      `json:"coins" yaml:"coins"`
}

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

type DistributionList struct {
	DistributionName string `json:"airdrop_name"`
}

func NewAirdropRecord(name string) DistributionList {
	return DistributionList{DistributionName: name}
}

func (ar DistributionList) Validate() bool {
	if ar.DistributionName == "" {
		return false
	}
	return true
}

func (ar DistributionList) String() string {
	return strings.TrimSpace(fmt.Sprintf(`DistributionName: %s`, ar.DistributionName))
}
