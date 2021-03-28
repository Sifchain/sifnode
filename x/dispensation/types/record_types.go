package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type DistributionRecord struct {
	AirdropName string         `json:"airdrop_name"`
	Address     sdk.AccAddress `json:"address" yaml:"address"`
	Coins       sdk.Coins      `json:"coins" yaml:"coins"`
}

func NewDistributionRecord(airdropName string, address sdk.AccAddress, coins sdk.Coins) DistributionRecord {
	return DistributionRecord{AirdropName: airdropName, Address: address, Coins: coins}
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
	return strings.TrimSpace(fmt.Sprintf(`AirdropName: %s
	Address: %s
	Coins: %s`, dr.AirdropName, dr.Address.String(), dr.Coins.String()))
}

func (ar DistributionRecord) Add(ar2 DistributionRecord) DistributionRecord {
	ar.Coins = ar.Coins.Add(ar2.Coins...)
	return ar
}

type AirdropRecord struct {
	AirdropName string `json:"airdrop_name"`
}

func NewAirdropRecord(airdropName string) AirdropRecord {
	return AirdropRecord{AirdropName: airdropName}
}

func (ar AirdropRecord) Validate() bool {
	if ar.AirdropName == "" {
		return false
	}
	return true
}

func (ar AirdropRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AirdropName: %s`, ar.AirdropName))
}
