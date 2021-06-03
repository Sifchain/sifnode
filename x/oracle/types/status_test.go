package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	statusText := PendingStatusText
	json, err := statusText.MarshalJSON()
	assert.NoError(t, err)
	err = statusText.UnmarshalJSON(json)
	assert.NoError(t, err)
}
