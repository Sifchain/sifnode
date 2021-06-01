package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

const (
	validatorAddress = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	power            = uint32(100)
)

func TestNewValidatorPerNetwork(t *testing.T) {
	list := NewValidatorPerNetwork()
	assert.Equal(t, len(list.Whitelist), 0)
}

func TestUpdateValidator(t *testing.T) {
	address, err := sdk.ValAddressFromBech32(validatorAddress)
	assert.NoError(t, err)

	list := NewValidatorPerNetwork()
	list.UpdateValidator(address, power)
	assert.Equal(t, len(list.GetAllValidators()), 1)
	assert.Equal(t, list.GetValidatorPower(address), power)

	list.UpdateValidator(address, uint32(0))
	assert.Equal(t, len(list.GetAllValidators()), 0)
	assert.Equal(t, list.GetValidatorPower(address), uint32(0))
}
