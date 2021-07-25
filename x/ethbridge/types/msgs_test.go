package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var expectedProphecyID = []byte{0x34, 0x74, 0xf9, 0xea, 0xe0, 0x8e, 0x24, 0x19, 0x5f, 0xf6, 0xb0, 0x59, 0xa2, 0x54, 0x10, 0xf6, 0x1f, 0x41, 0xac, 0x3d, 0xd5, 0x7f, 0x91, 0xc1, 0x52, 0xa0, 0xa5, 0x46, 0x5f, 0x31, 0xe5, 0xb7}

func TestComputeProphecyID(t *testing.T) {
	cosmosSender := "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	sequence := uint64(0)
	ethereumReceiver := "0x00000000000000000000"
	tokenAddress := "0x00000000000000000000"
	amount := sdk.NewInt(0)
	doublePeggy := false
	globalNonce := uint64(0)

	prophecy := ComputeProphecyID(cosmosSender, sequence, ethereumReceiver, tokenAddress, amount,
		doublePeggy, globalNonce)

	assert.Equal(t, prophecy, expectedProphecyID)
}
