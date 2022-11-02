package keeper_test

import (
	"bytes"
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	testNetwork          = 1
	testAddress          = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	testLockBurnSequence = uint64(10)
	testInitSequence     = uint64(0)
)

var testPrefix = []byte{0x6, 0x8, 0x1, 0x12, 0x14, 0xdc, 0xd3, 0xb2, 0xe3, 0xd8, 0x6a, 0x1, 0x3b, 0x5b, 0x5a, 0x82, 0x3b, 0x30, 0xf8, 0xfb, 0x79, 0x1b, 0xbc, 0xe, 0xa1}

func TestSetWitnessLockBurnNonce(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	lockBurnNonce := keeper.GetWitnessLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testInitSequence)

	keeper.SetWitnessLockBurnSequence(ctx, testNetwork, testCosmosAddress, testLockBurnSequence)

	lockBurnNonce = keeper.GetWitnessLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testLockBurnSequence)
}

func TestGetWitnessLockBurnSequence(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	lockBurnNonce := keeper.GetWitnessLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testInitSequence)
}

func TestGetWitnessLockBurnSequencePrefix(t *testing.T) {
	var _, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	lockBurnSequenceKey := types.LockBurnSequenceKey{
		NetworkDescriptor: testNetwork,
		ValidatorAddress:  testCosmosAddress,
	}

	prefix := lockBurnSequenceKey.GetWitnessLockBurnSequencePrefix(keeper.GetCdc())
	assert.Equal(t, prefix, testPrefix)
}

func TestGetAllWitnessLockBurnSequence(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testNetworkDescriptor := types.NetworkDescriptor(testNetwork)
	addresses := sifapp.CreateRandomAccounts(10)
	valAddresses := make([]sdk.ValAddress, 10)
	lockBurnSequences := make([]uint64, 10)
	for index := 0; index < 10; index++ {
		sequence := rand.Uint64()
		lockBurnSequences[index] = sequence
		valAddresses[index] = sdk.ValAddress(addresses[index])
		keeper.SetWitnessLockBurnSequence(ctx, testNetwork, valAddresses[index], sequence)
	}

	allSequences := keeper.GetAllWitnessLockBurnSequence(ctx)

	for _, sequence := range allSequences {
		key := sequence.WitnessLockBurnSequenceKey
		value := sequence.WitnessLockBurnSequence
		assert.Equal(t, key.NetworkDescriptor, testNetworkDescriptor)
		found := false
		for index := 0; index < 10; index++ {
			if bytes.Compare(valAddresses[index], key.ValidatorAddress) == 0 {
				found = true
				assert.Equal(t, value.LockBurnSequence, lockBurnSequences[index])
			}
		}
		assert.Equal(t, found, true)
	}

}
