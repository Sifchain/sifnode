package types_test

import (
	"math"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateValidatorPower(t *testing.T) {
	validatorWhitelist := types.ValidatorWhiteList{}
	addresses := sifapp.CreateRandomAccounts(2)
	valAddresses := sifapp.ConvertAddrsToValAddrs(addresses)

	err := validatorWhitelist.UpdateValidatorPower(valAddresses[0], 100)
	assert.NoError(t, err)

	err = validatorWhitelist.UpdateValidatorPower(valAddresses[1], math.MaxUint32)
	assert.ErrorIs(t, err, types.ErrValidatorPowerOverflow)

}

func Test_GetPowerRatio(t *testing.T) {
	validatorWhitelist := types.ValidatorWhiteList{}
	addresses := sifapp.CreateRandomAccounts(2)
	valAddresses := sifapp.ConvertAddrsToValAddrs(addresses)

	err := validatorWhitelist.UpdateValidatorPower(valAddresses[0], 100)
	assert.NoError(t, err)

	err = validatorWhitelist.UpdateValidatorPower(valAddresses[1], 900)
	assert.NoError(t, err)

	assert.Equal(t, validatorWhitelist.GetPowerRatio(valAddresses[:1]), 0.1)
}
