package types

import (
	"testing"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

func TestGetDenomHash(t *testing.T) {
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	tokenContractAddress := "0x0000000000000000000"
	expectedDenomHash := "sif9a8511b26f55f7c06088ef0705e5a9fa71df21bc8190de916ebf0db8d710f0aa"

	denomHash := GetDenomHash(networkDescriptor, tokenContractAddress)

	assert.Equal(t, denomHash, expectedDenomHash)
}
