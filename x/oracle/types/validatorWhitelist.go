package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorWhitelist define the address and its power
type ValidatorWhitelist struct {
	Whitelist map[string]uint32 `json:"whitelist"`
}

// DBValidatorWhitelist for db storage
type DBValidatorWhitelist struct {
	Whitelist []byte `json:"whitelist"`
}

// NewValidatorWhitelist get a new ValidatorWhitelist instance
func NewValidatorWhitelist() ValidatorWhitelist {
	return ValidatorWhitelist{
		Whitelist: make(map[string]uint32),
	}
}

// NewValidatorWhitelistFromData get a new ValidatorWhitelist instance
func NewValidatorWhitelistFromData(whitelist map[string]uint32) ValidatorWhitelist {
	return ValidatorWhitelist{
		Whitelist: whitelist,
	}
}

// UpdateValidator update validator's power
func (list *ValidatorWhitelist) UpdateValidator(validator sdk.ValAddress, power uint32) {
	list.Whitelist[validator.String()] = power
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
func (list ValidatorWhitelist) GetPowerRatio(claimValidators map[string][]sdk.ValAddress) (string, float64, float64) {
	var totalPower = uint32(0)
	for _, value := range list.Whitelist {
		totalPower += value
	}

	var totalClaimPower = uint32(0)
	var highestClaimPower = uint32(0)
	var highestString = ""

	for claim, validatorAddresses := range claimValidators {
		claimPower := uint32(0)
		for _, address := range validatorAddresses {
			claimPower += list.GetValidatorPower(address)
		}

		totalClaimPower += claimPower
		if claimPower > highestClaimPower {
			highestClaimPower = claimPower
			highestString = claim
		}
	}

	return highestString, float64(highestClaimPower) / float64(totalPower), float64(totalClaimPower-highestClaimPower) / float64(totalPower)
}

// GetAllValidators return all validators
func (list ValidatorWhitelist) GetAllValidators() []sdk.ValAddress {
	validators := make([]sdk.ValAddress, 0)
	for key, value := range list.Whitelist {
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
