package types

import (
	math "math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestParseNetworkDescriptorValid(t *testing.T) {
	input := "1"
	output, err := ParseNetworkDescriptor(input)
	require.Equal(t, NetworkDescriptor(1), output)
	require.NoError(t, err)
}

func TestParseNetworkDescriptorNegative(t *testing.T) {
	input := "-1"
	output, err := ParseNetworkDescriptor(input)
	require.Equal(t, -1, int(output))
	require.Error(t, err)
}

func TestParseNetworkDescriptorBeyondInt32(t *testing.T) {
	outOfRangeInput := strconv.Itoa(int(math.Pow(2, 33)))
	output, err := ParseNetworkDescriptor(outOfRangeInput)
	require.Equal(t, -1, int(output))
	require.Error(t, err)
}

func TestParseNetworkDescriptorNonExistNetworkDescriptor(t *testing.T) {
	nonExistNetworkDescriptor := "555"
	output, err := ParseNetworkDescriptor(nonExistNetworkDescriptor)
	require.Equal(t, -1, int(output))
	require.Error(t, err)
}
