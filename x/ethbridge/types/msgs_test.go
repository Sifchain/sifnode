package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var expectedProphecyID = []byte{0xc5, 0xa, 0xd0, 0x6c, 0x42, 0x75, 0x49, 0x55, 0xc, 0x7b, 0x37, 0xe1, 0x9a, 0xcb, 0xc1, 0xbb, 0x50, 0xa2, 0x70, 0x99, 0xdf, 0xbb, 0xa4, 0xdd, 0x88, 0x54, 0x51, 0x59, 0xba, 0xae, 0x4, 0x65}

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

	assert.Equal(t, expectedProphecyID, prophecy)
}
