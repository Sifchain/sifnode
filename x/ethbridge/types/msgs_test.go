package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var expectedProphecyID = []byte{0x62, 0x65, 0xae, 0x1e, 0x1a, 0xdb, 0x40, 0x92, 0xb3, 0x1b, 0xf, 0x87, 0x45, 0x46, 0xa6, 0x47, 0x24, 0x1, 0xa9, 0xc0, 0x56, 0x56, 0x71, 0x8c, 0x38, 0xc, 0x3e, 0x6, 0x15, 0x24, 0xe8, 0x6c}

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
