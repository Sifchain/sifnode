package types

const (
	QueryAllDistributions   = "distributions"
	QueryRecordsByDistrName = "records_by_name"
	QueryRecordsByRecipient = "records_by_recipient"
	QueryClaimsByType       = "claims_by_type"
)

func NewQueryRecordsByDistributionName(distributionName string, status DistributionStatus) QueryRecordsByDistributionNameRequest {
	return QueryRecordsByDistributionNameRequest{DistributionName: distributionName, Status: status}
}

func NewQueryRecordsByRecipientAddr(address string) QueryRecordsByRecipientAddrRequest {
	return QueryRecordsByRecipientAddrRequest{Address: address}
}
