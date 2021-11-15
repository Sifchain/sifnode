package types

import (
	"testing"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

func TestGetDenomHash(t *testing.T) {
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	tokenContractAddress := NewEthereumAddress("0x0000000000000000000")
	expectedDenomHash := "siffa33aa4b83b0e09f21c221b25b6e46480ae151a36932dc44fd09f4f073e9f54f"

	denomHash := GetDenomHash(networkDescriptor, tokenContractAddress)

	assert.Equal(t, denomHash, expectedDenomHash)
}
