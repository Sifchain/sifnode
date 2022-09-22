package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
)

const (
	testNetwork       = 1
	testAddress       = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	testLockBurnNonce = uint64(10)
	testInitNonce     = uint64(0)
)

var testPrefix = []byte{0x6, 0x8, 0x1, 0x12, 0x14, 0xdc, 0xd3, 0xb2, 0xe3, 0xd8, 0x6a, 0x1, 0x3b, 0x5b, 0x5a, 0x82, 0x3b, 0x30, 0xf8, 0xfb, 0x79, 0x1b, 0xbc, 0xe, 0xa1}

func TestSetWitnessLockBurnNonce(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	lockBurnNonce := keeper.GetWitnessLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testInitNonce)

	keeper.SetWitnessLockBurnNonce(ctx, testNetwork, testCosmosAddress, testLockBurnNonce)

	lockBurnNonce = keeper.GetWitnessLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testLockBurnNonce)
}

func TestGetWitnessLockBurnSequence(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	lockBurnNonce := keeper.GetWitnessLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testInitNonce)
}

func TestGetWitnessLockBurnSequencePrefix(t *testing.T) {
	var _, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	prefix := keeper.GetWitnessLockBurnSequencePrefix(testNetwork, testCosmosAddress)
	assert.Equal(t, prefix, testPrefix)
}
