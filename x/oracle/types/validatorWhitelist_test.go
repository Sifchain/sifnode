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

func TestNewValidatorWhitelist(t *testing.T) {
	list := NewValidatorWhitelist()
	assert.Equal(t, len(list.Whitelist), 0)
}

func TestAddValidator(t *testing.T) {
	address, err := sdk.ValAddressFromBech32(validatorAddress)
	assert.NoError(t, err)

	list := NewValidatorWhitelist()
	list.AddValidator(address, power)
	assert.Equal(t, len(list.Whitelist), 1)
	assert.Equal(t, list.GetValidatorPower(address), power)

}

func TestRemoveValidator(t *testing.T) {
	address, err := sdk.ValAddressFromBech32(validatorAddress)
	assert.NoError(t, err)

	list := NewValidatorWhitelist()
	list.AddValidator(address, power)
	assert.Equal(t, len(list.Whitelist), 1)
	assert.Equal(t, list.GetValidatorPower(address), power)

	list.RemoveValidator(address)
	assert.Equal(t, len(list.Whitelist), 1)
	assert.Equal(t, list.GetValidatorPower(address), uint32(0))
}
