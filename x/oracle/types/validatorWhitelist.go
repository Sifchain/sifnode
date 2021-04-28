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
	power := uint32(0)
	power = list.Whitelist[validator.String()]
	return power
}
