package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewClaim returns a new Claim
func NewClaim(id string, validatorAddress sdk.ValAddress, content string) Claim {
	return Claim{
		ID:               id,
		ValidatorAddress: validatorAddress,
		Content:          content,
	}
}
