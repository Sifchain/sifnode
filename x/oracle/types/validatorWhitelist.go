package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorWhitelist define the address and its power
type ValidatorWhitelist struct {
	Whitelist map[string]uint32 `json:"whitelist"`
}

// NewValidatorWhitelist get a new ValidatorWhitelist instance
func NewValidatorWhitelist() ValidatorWhitelist {
	return ValidatorWhitelist{
		Whitelist: make(map[string]uint32),
	}
}

// NewValidatorWhitelist get a new ValidatorWhitelist instance
func NewValidatorWhitelistFromData(whitelist map[string]uint32) ValidatorWhitelist {
	return ValidatorWhitelist{
		Whitelist: whitelist,
	}
}

// AddValidator add new validator and its power
func (list *ValidatorWhitelist) AddValidator(validator sdk.ValAddress, power uint32) {
	list.Whitelist[validator.String()] = power
}

// RemoveValidator just set its power as 0
func (list *ValidatorWhitelist) RemoveValidator(validator sdk.ValAddress) {
	list.Whitelist[validator.String()] = 0
}

// GetValidatorPower return validator's power
func (list *ValidatorWhitelist) GetValidatorPower(validator sdk.ValAddress) uint32 {
	if list.ContainValidator(validator) {
		return list.Whitelist[validator.String()]
	}

	return 0
}

// ContainValidator return if validator in the map
func (list ValidatorWhitelist) ContainValidator(validator sdk.ValAddress) bool {
	_, ok := list.Whitelist[validator.String()]
	return ok
}

// GetPowerRatio return the power ratio of input validator address list
func (list ValidatorWhitelist) GetPowerRatio(claimValidators map[string][]sdk.ValAddress) (string, float64) {
	var totalPower = uint32(0)
	for _, value := range list.Whitelist {
		totalPower += value
	}

	var highestClaimPower = uint32(0)
	var highestString = ""

	for claim, validatorAddresses := range claimValidators {
		claimPower := uint32(0)
		for _, address := range validatorAddresses {
			claimPower += list.GetValidatorPower(address)
		}

		if claimPower > highestClaimPower {
			highestClaimPower = claimPower
			highestString = claim
		}
	}

	return highestString, float64(highestClaimPower) / float64(totalPower)
}
