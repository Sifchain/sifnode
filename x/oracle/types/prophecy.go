package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultConsensusNeeded defines the default consensus value required for a
// prophecy to be finalized
const DefaultConsensusNeeded float64 = 0.7

// AddClaim adds a given claim to this prophecy
func (prophecy *Prophecy) AddClaim(validator sdk.ValAddress) {
	prophecy.ClaimValidators = append(prophecy.ClaimValidators, validator.String())
}

// func inWhiteList(validator staking.Validator, whiteListValidatorAddresses []sdk.ValAddress) bool {
// 	for _, whiteListValidatorAddress := range whiteListValidatorAddresses {
// 		if bytes.Equal(validator.GetOperator(), whiteListValidatorAddress) {
// 			return true
// 		}
// 	}
// 	return false
// }

// FindHighestClaim looks through all the existing claims on a given prophecy. It adds up the total power across
// all claims and returns the highest claim, power for that claim, total power claimed on the prophecy overall.
// and the total power of all whitelist validators.
func (prophecy Prophecy) FindHighestClaim(ctx sdk.Context, stakeKeeper StakingKeeper, whiteListValidatorAddresses []sdk.ValAddress) (string, int64, int64, int64) {
	return "", 0, 0, 0
}
