package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var expectedProphecyID = []byte{0xb8, 0x44, 0xea, 0xb0, 0x1e, 0xcc, 0xa, 0xa3, 0x85, 0x71, 0xe8, 0xff, 0x96, 0x1d, 0x34, 0x5b, 0x4, 0x6e, 0x4d, 0x9c, 0xd8, 0x85, 0x9, 0x7d, 0x99, 0xee, 0x6c, 0xa8, 0x34, 0x49, 0xb4, 0x19}

func TestComputeProphecyID(t *testing.T) {
	cosmosSender := "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	sequence := uint64(0)
	ethereumReceiver := "0x010203040506070809"
	tokenAddress := "0x090807060504030201"
	amount := sdk.NewInt(1025)
	bridgeToken := true
	globalNonce := uint64(0)
	denom := "sifBridge0123456789"

	prophecy := ComputeProphecyID(cosmosSender, sequence, ethereumReceiver, tokenAddress, amount,
		bridgeToken, globalNonce, TestNetworkDescriptor, denom)

	assert.Equal(t, expectedProphecyID, prophecy)
}
