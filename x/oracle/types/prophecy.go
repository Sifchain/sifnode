package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultConsensusNeeded defines the default consensus value required for a
// prophecy to be finalized
const DefaultConsensusNeeded float64 = 0.7

// Prophecy is a struct that contains all the metadata of an oracle ritual.
// Claims are indexed by the claim's validator bech32 address and by the claim's json value to allow
// for constant lookup times for any validation/verifiation checks of duplicate claims
// Each transaction, pending potential results are also calculated, stored and indexed by their byte result
// to allow discovery of consensus on any the result in constant time without having to sort or run
// through the list of claims to find the one with highest consensus
type Prophecy struct {
	ID     string `json:"id"`
	Status Status `json:"status"`

	//WARNING: Mappings are nondeterministic in Amino,
	// an so iterating over them could result in consensus failure. New code should not iterate over the below 2 mappings.

	//This is a mapping from a claim to the list of validators that made that claim.
	ClaimValidators map[string][]sdk.ValAddress `json:"claim_validators"`
	//This is a mapping from a validator bech32 address to their claim
	ValidatorClaims map[string]string `json:"validator_claims"`
}

// SerializeForDB serializes a prophecy into a DBProphecy
// TODO: Using gob here may mean that different tendermint clients in different languages may serialize/store
// prophecies in their db in different ways -
// check with @codereviewer if this is ok or if it introduces a risk of creating forks.
// Or maybe using a slower json serializer or Amino:JSON would be ok
func (prophecy Prophecy) SerializeForDB() (DBProphecy, error) {
	claimValidators, err := json.Marshal(prophecy.ClaimValidators)
	if err != nil {
		return DBProphecy{}, err
	}

	validatorClaims, err := json.Marshal(prophecy.ValidatorClaims)
	if err != nil {
		return DBProphecy{}, err
	}

	return DBProphecy{
		Id:              prophecy.ID,
		Status:          prophecy.Status,
		ClaimValidators: claimValidators,
		ValidatorClaims: validatorClaims,
	}, nil
}

// DeserializeFromDB deserializes a DBProphecy into a prophecy
func (dbProphecy DBProphecy) DeserializeFromDB() (Prophecy, error) {
	var claimValidators map[string][]sdk.ValAddress
	if err := json.Unmarshal(dbProphecy.ClaimValidators, &claimValidators); err != nil {
		return Prophecy{}, err
	}

	var validatorClaims map[string]string
	if err := json.Unmarshal(dbProphecy.ValidatorClaims, &validatorClaims); err != nil {
		return Prophecy{}, err
	}

	return Prophecy{
		ID:              dbProphecy.Id,
		Status:          dbProphecy.Status,
		ClaimValidators: claimValidators,
		ValidatorClaims: validatorClaims,
	}, nil
}

// AddClaim adds a given claim to this prophecy
func (prophecy Prophecy) AddClaim(validator sdk.ValAddress, claim string) {
	claimValidators := prophecy.ClaimValidators[claim]
	prophecy.ClaimValidators[claim] = append(claimValidators, validator)

	validatorBech32 := validator.String()
	prophecy.ValidatorClaims[validatorBech32] = claim
}

// NewProphecy returns a new Prophecy, initialized in pending status with an initial claim
func NewProphecy(id string) Prophecy {
	return Prophecy{
		ID:              id,
		Status:          NewStatus(StatusText_STATUS_TEXT_PENDING, ""),
		ClaimValidators: make(map[string][]sdk.ValAddress),
		ValidatorClaims: make(map[string]string),
	}
}

// NewStatus returns a new Status with the given data contained
func NewStatus(text StatusText, finalClaim string) Status {
	return Status{
		Text:       text,
		FinalClaim: finalClaim,
	}
}
