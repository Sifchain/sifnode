package types

// ----------------------------------------------------------------------------
// Client Types

type DistributionRecordsResponse struct {
	DistributionRecords
	Height int64 `json:"height"`
}

func NewDistributionRecordsResponse(distributionRecords DistributionRecords, height int64) DistributionRecordsResponse {
	return DistributionRecordsResponse{DistributionRecords: distributionRecords, Height: height}
}

type DistributionListsResponse struct {
	Distributions
	Height int64 `json:"height"`
}

func NewDistributionListsResponse(distributionLists Distributions, height int64) DistributionListsResponse {
	return DistributionListsResponse{Distributions: distributionLists, Height: height}
}
