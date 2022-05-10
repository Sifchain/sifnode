package cli

import (
	"math"
	"strconv"
	"testing"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestParseNetworkDescriptorValid(t *testing.T) {
	input := "1"
	output, err := parseNetworkDescriptor(input)
	require.Equal(t, oracletypes.NetworkDescriptor(1), output)
	require.NoError(t, err)
}

func TestParseNetworkDescriptorNegative(t *testing.T) {
	input := "-1"
	output, err := parseNetworkDescriptor(input)
	require.Equal(t, -1, int(output))
	require.Error(t, err)
}

func TestParseNetworkDescriptorBeyondInt32(t *testing.T) {
	outOfRangeInput := strconv.Itoa(int(math.Pow(2, 33)))
	output, err := parseNetworkDescriptor(outOfRangeInput)
	require.Equal(t, -1, int(output))
	require.Error(t, err)
}

func TestParseNetworkDescriptorNonExistNetworkDescriptor(t *testing.T) {
	nonExistNetworkDescriptor := "555" //
	output, err := parseNetworkDescriptor(nonExistNetworkDescriptor)
	require.Equal(t, -1, int(output))
	require.Error(t, err)
}
