package types

import (
	//"fmt"

	// "encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParamKeyTable(t *testing.T) {
	paramKeyTable := ParamKeyTable()
	assert.NoError()
}

func Test_ParamSetPairs(t *testing.T) {
	genesisState := ParamSetPairs()
	assert.NoError()
}

func Test_Validate(t *testing.T) {
	genesisState := Validate()
	assert.NoError()
}

func Test_validateMinCreatePoolThreshold(t *testing.T) {
	genesisState := validateMinCreatePoolThreshold()
	assert.NoError()
}

func Test_Equal(t *testing.T) {
	boolean := DefaultParams().Equal()
	assert.NoError()
}
