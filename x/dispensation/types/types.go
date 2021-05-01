package types

// ----------------------------------------------------------------------------
// Client Types

func NewQueryAllDistributionsResponse(distributions Distributions, height int64) QueryAllDistributionsResponse {
	return QueryAllDistributionsResponse{
		Distributions: distributions.Distributions,
		Height: height,
	}
}

func NewQueryRecordsByDistributionNameResponse(distributionRecords DistributionRecords, height int64) QueryRecordsByDistributionNameResponse {
	return QueryRecordsByDistributionNameResponse{
		DistributionRecords: &distributionRecords,
		Height: height,
	}
}

func NewQueryRecordsByRecipientAddrResponse(distributionRecords DistributionRecords, height int64) QueryRecordsByRecipientAddrResponse {
	return QueryRecordsByRecipientAddrResponse{
		DistributionRecords: &distributionRecords,
		Height: height,
	}
}