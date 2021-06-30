package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	networkDescriptor = NetworkDescriptor(1)
)

func TestNewNetworkIdentity(t *testing.T) {
	networkIdentity := NewNetworkIdentity(networkDescriptor)
	assert.Equal(t, networkDescriptor, networkIdentity.NetworkDescriptor)
}

func TestGetPrefix(t *testing.T) {
	prefixOfNetwork100 := []byte{0x0, 0x0, 0x0, 0x0, 0x01}
	networkIdentity := NewNetworkIdentity(networkDescriptor)
	assert.Equal(t, networkDescriptor, networkIdentity.NetworkDescriptor)
	assert.Equal(t, networkIdentity.GetPrefix(), prefixOfNetwork100)
}

func TestNetworkDescriptorValid(t *testing.T) {
	networkDescriptor := NetworkDescriptor(0)
	assert.Equal(t, networkDescriptor.IsValid(), false)

	networkDescriptor = NetworkDescriptor(1)
	assert.Equal(t, networkDescriptor.IsValid(), true)

	networkDescriptor = NetworkDescriptor(99999)
	assert.Equal(t, networkDescriptor.IsValid(), false)
}
