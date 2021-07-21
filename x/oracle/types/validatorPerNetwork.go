package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateValidator reset validator's power
func (list *ValidatorWhiteList) UpdateValidator(validator sdk.ValAddress, power uint32) {
	list.GetWhiteList()[validator.String()] = power
}

// GetValidatorPower return validator's power
func (list *ValidatorWhiteList) GetValidatorPower(validator sdk.ValAddress) uint32 {
	if list.ContainValidator(validator) {
		return list.GetWhiteList()[validator.String()]
	}

	return 0
}

// ContainValidator return if validator in the map
func (list *ValidatorWhiteList) ContainValidator(validator sdk.ValAddress) bool {
	_, ok := list.GetWhiteList()[validator.String()]
	return ok
}

// GetPowerRatio return the power ratio of input validator address list
func (list *ValidatorWhiteList) GetPowerRatio(claimValidators []string) float64 {
	var totalPower = uint32(0)
	var votePower = uint32(0)
	for key, value := range list.GetWhiteList() {
		totalPower += value
		for _, validator := range claimValidators {
			if key == validator {
				votePower += value
			}
		}
	}

	return float64(votePower) / float64(totalPower)
}

// GetAllValidators return all validators
func (list *ValidatorWhiteList) GetAllValidators() []sdk.ValAddress {
	validators := make([]sdk.ValAddress, 0)
	for key, value := range list.GetWhiteList() {
		address, err := sdk.ValAddressFromBech32(key)
		if err != nil {
			panic("invalid address in whitelist")
		}
		if value > 0 {
			validators = append(validators, address)
		}
	}

	return validators
}
