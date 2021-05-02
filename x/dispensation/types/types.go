package types

import "github.com/cosmos/cosmos-sdk/types/rest"

// ----------------------------------------------------------------------------
// Client Types

type DistributionRecordsResponse struct {
	DistributionRecords
	Height int64 `json:"height"`
}

func NewDistributionRecordsResponse(distributionRecords DistributionRecords, height int64) DistributionRecordsResponse {
	return DistributionRecordsResponse{DistributionRecords: distributionRecords, Height: height}
}

type DistributionsResponse struct {
	Distributions
	Height int64 `json:"height"`
}

func NewDistributionsResponse(distributions Distributions, height int64) DistributionsResponse {
	return DistributionsResponse{Distributions: distributions, Height: height}
}

type ClaimsResponse struct {
	Claims []UserClaim `json:"claims"`
	Height int64       `json:"height"`
}

func NewClaimsResponse(claims []UserClaim, height int64) ClaimsResponse {
	return ClaimsResponse{Claims: claims, Height: height}
}

type CreateClaimReq struct {
	BaseReq   rest.BaseReq     `json:"base_req"`
	Signer    string           `json:"signer"`
	ClaimType DistributionType `json:"claim_type"`
}
