package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultConsensusNeeded defines the default consensus value required for a
// prophecy to be finalized
const DefaultConsensusNeeded float64 = 0.7

// AddClaim adds a given claim to this prophecy
func (prophecy *Prophecy) AddClaim(address sdk.ValAddress) {
	validators := prophecy.ClaimValidators
	for _, validator := range validators {
		if validator == address.String() {
			return
		}
	}
	prophecy.ClaimValidators = append(prophecy.ClaimValidators, address.String())
}

// GetVoteRate return vote rate according to validator's power
func (prophecy Prophecy) GetVoteRate(ctx sdk.Context) float64 {
	// TODO have other pr to compute with validator's power
	return 0.0
}
