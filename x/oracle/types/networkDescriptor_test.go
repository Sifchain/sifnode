package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	networkID = NetworkID(1)
)

func TestNewNetworkDescriptor(t *testing.T) {
	networkDescriptor := NewNetworkDescriptor(networkID)
	assert.Equal(t, networkDescriptor.NetworkID, networkID)
}

func TestGetPrefix(t *testing.T) {
	prefixOfNetwork100 := []byte{0x0, 0x0, 0x0, 0x0, 0x01}
	networkDescriptor := NewNetworkDescriptor(networkID)
	assert.Equal(t, networkDescriptor.NetworkID, networkID)
	assert.Equal(t, networkDescriptor.GetPrefix(), prefixOfNetwork100)
}
