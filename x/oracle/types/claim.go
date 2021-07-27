package types

// NewClaim returns a new Claim
func NewClaim(id string, validatorAddress string, content string) Claim {
	return Claim{
		Id:               id,
		ValidatorAddress: validatorAddress,
		Content:          content,
	}
}
