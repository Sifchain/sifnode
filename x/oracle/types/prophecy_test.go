package types

import (
	"testing"

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
	assert.Equal(t, prophecy.ClaimValidators[LocalClaim], []sdk.ValAddress{address})
	assert.Equal(t, prophecy.ValidatorClaims[address.String()], LocalClaim)
}

func TestSerializeForDB(t *testing.T) {
	prophecy := NewProphecy(ProphecyID)
	assert.Equal(t, prophecy.ID, ProphecyID)
	address, err := sdk.ValAddressFromBech32(LocalValidatorAddress)
	assert.NoError(t, err)

	prophecy.AddClaim(address, LocalClaim)
	assert.Equal(t, prophecy.ClaimValidators[LocalClaim], []sdk.ValAddress{address})
	assert.Equal(t, prophecy.ValidatorClaims[address.String()], LocalClaim)

	dBProphecy, err := prophecy.SerializeForDB()
	assert.NoError(t, err)
	assert.Equal(t, dBProphecy.ID, ProphecyID)

	prophecy, err = dBProphecy.DeserializeFromDB()
	assert.NoError(t, err)
	assert.Equal(t, prophecy.ID, ProphecyID)
}

func TestNewStatus(t *testing.T) {
	status := NewStatus(0, LocalClaim)
	assert.Equal(t, status.Text, StatusText(0))
	assert.Equal(t, status.FinalClaim, LocalClaim)
}
