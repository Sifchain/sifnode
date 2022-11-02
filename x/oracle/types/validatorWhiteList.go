package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateValidatorPower reset validator's power
func (list *ValidatorWhiteList) UpdateValidatorPower(validator sdk.ValAddress, power uint32) error {
	totalPower := uint32(0)
	for _, value := range list.ValidatorPower {
		if bytes.Compare(value.ValidatorAddress, validator) == 0 {
			value.VotingPower = power
			totalPower += value.VotingPower
			if totalPower < value.VotingPower {
				return ErrValidatorPowerOverflow
			}
			return nil
		}
		totalPower += value.VotingPower
	}

	// if validator is a new one
	totalPower += power
	if totalPower < power {
		return ErrValidatorPowerOverflow
	}

	list.ValidatorPower = append(list.ValidatorPower, &ValidatorPower{
		ValidatorAddress: validator,
		VotingPower:      power,
	})
	return nil
}

// GetPowerRatio return the power ratio of input validator address list
func (list *ValidatorWhiteList) GetPowerRatio(claimValidators []sdk.ValAddress) float64 {
	var totalPower = uint32(0)
	var votePower = uint32(0)
	for _, value := range list.ValidatorPower {
		totalPower += value.VotingPower
		for _, validator := range claimValidators {
			if bytes.Compare(value.ValidatorAddress, validator) == 0 {
				votePower += value.VotingPower
			}
		}
	}

	// if no validator, return 0.0
	if totalPower == 0 {
		return 0.0
	}

	return float64(votePower) / float64(totalPower)
}
