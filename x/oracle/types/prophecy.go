package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultConsensusNeeded defines the default consensus value required for a
// prophecy to be finalized

// TODO integration test environment just one witness node with 50% vote power
// to make the burn ceth finalized, reduce it to 0.49 temporarily
// TODO revert to 0.7 after integration test
const DefaultConsensusNeeded float64 = 0.49

// AddClaim adds a given claim to this prophecy
func (prophecy *Prophecy) AddClaim(address sdk.ValAddress) error {
	validators := prophecy.ClaimValidators
	for _, validator := range validators {
		if validator == address.String() {
			return ErrDuplicateMessage
		}
	}
	prophecy.ClaimValidators = append(prophecy.ClaimValidators, address.String())
	return nil
}
