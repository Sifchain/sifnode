package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

const (
	ProphecyID            = "ProphecyID"
	LocalValidatorAddress = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	LocalClaim            = "claim"
)

func TestNewProphecy(t *testing.T) {
	prophecy := NewProphecy(ProphecyID)
	assert.Equal(t, prophecy.ID, ProphecyID)
}

func TestAddClaim(t *testing.T) {
	prophecy := NewProphecy(ProphecyID)
	assert.Equal(t, prophecy.ID, ProphecyID)
	address, err := sdk.ValAddressFromBech32(LocalValidatorAddress)
	assert.NoError(t, err)

	prophecy.AddClaim(address, LocalClaim)
	assert.Equal(t, prophecy.ClaimValidators[LocalClaim], []types.ValAddress{address})
	assert.Equal(t, prophecy.ValidatorClaims[address.String()], LocalClaim)
}
