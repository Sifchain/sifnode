package types

import (
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQueryEthProphecyRequest creates a new QueryEthProphecyParams
func NewQueryEthProphecyRequest(prophecyID []byte) *QueryEthProphecyRequest {
	return &QueryEthProphecyRequest{
		ProphecyId: prophecyID,
	}
}

// NewQueryEthProphecyResponse creates a new QueryEthProphecyResponse instance
func NewQueryEthProphecyResponse(id []byte, status oracletypes.StatusText, claims []string) QueryEthProphecyResponse {
	// claimValidators := []string{}
	// for _, claim := range claims {
	// 	claimValidators = append(claimValidators, claim.ValidatorAddress)
	// }

	return QueryEthProphecyResponse{
		ProphecyId:      id,
		Status:          status,
		ClaimValidators: claims,
	}
}
