package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParamSetPairs(t *testing.T) {
	params := DefaultParams()
	_ = params.ParamSetPairs()
}

func Test_Validate(t *testing.T) {
	params := DefaultParams()
	err := params.Validate()
	assert.NoError(t, err)
}

func Test_ParamsEqual(t *testing.T) {
	params1 := DefaultParams()
	params2 := DefaultParams()
	boolean := params1.Equal(params2)
	assert.True(t, boolean)
	boolean = params1.Equal(NewParams(uint64(10)))
	assert.False(t, boolean)
}
