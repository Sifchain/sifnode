package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryAllDistributions   = "distributions"
	QueryRecordsByDistrName = "records_by_name"
	QueryRecordsByRecipient = "records_by_recipient"
	QueryClaimsByType       = "claims_by_type"
)

type QueryRecordsByDistributionName struct {
	DistributionName string             `json:"distribution_name"`
	Status           DistributionStatus `json:"status"`
}

func NewQueryRecordsByDistributionName(distributionName string, status DistributionStatus) QueryRecordsByDistributionName {
	return QueryRecordsByDistributionName{DistributionName: distributionName, Status: status}
}

type QueryRecordsByRecipientAddr struct {
	Address sdk.AccAddress `json:"address"`
}

func NewQueryRecordsByRecipientAddr(address sdk.AccAddress) QueryRecordsByRecipientAddr {
	return QueryRecordsByRecipientAddr{Address: address}
}

type QueryUserClaims struct {
	UserClaimType DistributionType `json:"user_claim_type"`
}

func NewQueryUserClaims(userClaimType DistributionType) QueryUserClaims {
	return QueryUserClaims{UserClaimType: userClaimType}
}
