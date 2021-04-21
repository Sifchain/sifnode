package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryAllDistributions   = "distributions"
	QueryRecordsByDistrName = "records_by_name"
	QueryRecordsByRecipient = "records_by_recipient"
)

func NewQueryRecordsByDistributionName(distributionName string, status sdk.Uint) QueryRecordsByDistributionName {
	return QueryRecordsByDistributionName{DistributionName: distributionName, Status: status}
}

func NewQueryRecordsByRecipientAddr(address string) QueryRecordsByRecipientAddr {
	return QueryRecordsByRecipientAddr{Address: address}
}
