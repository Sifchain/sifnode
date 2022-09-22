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

func TestNetworkDescriptorValid(t *testing.T) {
	networkDescriptor := NetworkDescriptor(0)
	assert.Equal(t, networkDescriptor.IsSifchain(), true)

	networkDescriptor = NetworkDescriptor(1)
	assert.Equal(t, networkDescriptor.IsValid(), true)

	networkDescriptor = NetworkDescriptor(99999)
	assert.Equal(t, networkDescriptor.IsValid(), false)
}
