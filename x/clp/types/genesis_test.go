package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewGenesisState(t *testing.T) {
	genesisState := NewGenesisState(DefaultParams())
	assert.Equal(t, genesisState.Params.MinCreatePoolThreshold, uint64(100))
}

func Test_DefaultGenesisState(t *testing.T) {
	genesisState := DefaultGenesisState()
	assert.Equal(t, genesisState.Params.MinCreatePoolThreshold, uint64(100))
	assert.Equal(t, genesisState.AddressWhitelist, []string{"cosmos1ny48eeuk4dm9f63dy0lwfgjhnvud9yvt8tcaat"})
}
