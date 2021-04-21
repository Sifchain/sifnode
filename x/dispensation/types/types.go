package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ----------------------------------------------------------------------------
// Client Types

func NewDistributionRecordsResponse(distributionRecords DistributionRecords, height sdk.Int) DistributionRecordsResponse {
	return DistributionRecordsResponse{DistributionRecords: &distributionRecords, Height: height}
}

type DistributionsResponse struct {
	Distributions
	Height int64 `json:"height"`
}

func NewDistributionsResponse(distributions Distributions, height int64) DistributionsResponse {
	return DistributionsResponse{Distributions: distributions, Height: height}
}
