package types

const (
	QueryAllDistributions   = "distributions"
	QueryRecordsByDistrName = "records_by_name"
	QueryRecordsByRecipient = "records_by_recipient"
)

func NewQueryRecordsByDistributionName(distributionName string, status ClaimStatus) QueryRecordsByDistributionName {
	return QueryRecordsByDistributionName{DistributionName: distributionName, Status: status}
}

func NewQueryRecordsByRecipientAddr(address string) QueryRecordsByRecipientAddr {
	return QueryRecordsByRecipientAddr{Address: address}
}
