package types

// ----------------------------------------------------------------------------
// Client Types

func NewDistributionRecordsResponse(distributionRecords DistributionRecords, height int64) DistributionRecordsResponse {
	return DistributionRecordsResponse{DistributionRecords: &distributionRecords, Height: height}
}

func NewDistributionsResponse(distributions Distributions, height int64) DistributionsResponse {
	return DistributionsResponse{Distributions: distributions.Distributions, Height: height}
}
