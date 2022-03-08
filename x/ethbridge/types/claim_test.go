package types

import (
	"testing"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

func TestGetDenomHash(t *testing.T) {
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	tokenContractAddress := NewEthereumAddress("0x0000000000000000000")
	expectedDenomHash := "sifBridge00010x0000000000000000000000000000000000000000"

	denomHash := GetDenom(networkDescriptor, tokenContractAddress)

	assert.Equal(t, expectedDenomHash, denomHash)
}
