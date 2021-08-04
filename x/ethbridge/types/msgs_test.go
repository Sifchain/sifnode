package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var expectedProphecyID = []byte{0xb3, 0x95, 0x2a, 0xfc, 0x2f, 0xaa, 0xa, 0x1f, 0xcd, 0xfe, 0x4, 0x79, 0x40, 0x4a, 0x25, 0xf3, 0x72, 0x1e, 0xb1, 0x89, 0x7a, 0x8, 0x98, 0x17, 0x2a, 0x55, 0xf3, 0x8f, 0xdc, 0xc8, 0xad, 0x37}

func TestComputeProphecyID(t *testing.T) {
	cosmosSender := "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	sequence := uint64(0)
	ethereumReceiver := "0x00000000000000000000"
	tokenAddress := "0x00000000000000000000"
	amount := sdk.NewInt(0)
	doublePeggy := false
	globalNonce := uint64(0)

	prophecy := ComputeProphecyID(cosmosSender, sequence, ethereumReceiver, tokenAddress, amount,
		doublePeggy, globalNonce, TestNetworkDescriptor)

	assert.Equal(t, prophecy, expectedProphecyID)
}
